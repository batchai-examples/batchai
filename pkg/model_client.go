package batchai

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/pkg/errors"
	"github.com/tiktoken-go/tokenizer"

	"github.com/qiangyt/batchai/comm"
)

type ModelClientT struct {
	config                ModelConfig
	openAiClient          *openai.Client
	openAiStreamingClient *openai.Client
	semaphore             chan struct{}
	codec                 tokenizer.Codec
}

type ModelClient = *ModelClientT

func NewModelClient(config ModelConfig) ModelClient {
	codec, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize Cl100kBase tokenizer codec"))
	}

	return &ModelClientT{
		config:                config,
		openAiClient:          buildOpenAiClient(config, false),
		openAiStreamingClient: buildOpenAiClient(config, true),
		semaphore:             make(chan struct{}, 1),
		codec:                 codec,
	}
}

func buildModelClient(config AppConfig, modelId string) ModelClient {
	model := config.LoadModel(modelId)
	return NewModelClient(model)
}

func (me ModelClient) acquire() {
	me.semaphore <- struct{}{}
}

func (me ModelClient) release() {
	<-me.semaphore
}

func (me ModelClient) Encode(msg string) ([]uint, []string) {
	tokens, texts, err := me.codec.Encode(msg)
	if err != nil {
		panic(errors.Wrap(err, "failed to encode text"))
	}
	return tokens, texts
}

func (me ModelClient) Decode(tokens []uint) string {
	text, err := me.codec.Decode(tokens)
	if err != nil {
		panic(errors.Wrap(err, "failed to decode tokens"))
	}
	return text
}

func (me ModelClient) EvaluatedTokens(prompt string) int {
	tokens, _ := me.Encode(prompt)
	return len(tokens)
}

func (me ModelClient) Chat(x Kontext, c comm.Console, saveIntoMemory bool, memory ChatMemory, writer io.Writer) (*openai.ChatCompletion, time.Duration) {
	me.acquire()
	defer me.release()

	startTime := time.Now()
	x = x.Timeouted(me.config.Timeout)

	cfg := me.config

	requestedMaxCompletionTokens := cfg.MaxCompletionTokens
	if requestedMaxCompletionTokens > 0 {
		promptTokens, _ := me.Encode(memory.Format())
		lenOfPromptTokens := int64(len(promptTokens) + 16)

		contextWindow := cfg.ContextWindow
		allowedMaxCompletionTokens := contextWindow - lenOfPromptTokens

		if requestedMaxCompletionTokens >= allowedMaxCompletionTokens {
			c.Yellowf("pre-check warning: %s maximum context length is %d tokens. However, you requested %d tokens (%d in the messages, %d in the completion). Force to reduce the length of the messages or completion to be %d.",
				cfg.Id, contextWindow, requestedMaxCompletionTokens+lenOfPromptTokens, lenOfPromptTokens, requestedMaxCompletionTokens, allowedMaxCompletionTokens)
		}

		requestedMaxCompletionTokens = allowedMaxCompletionTokens
	}

	var r *openai.ChatCompletion
	if writer == nil {
		r = me.chat(x, memory, requestedMaxCompletionTokens)
	} else {
		r = me.chatStream(x, memory, writer, requestedMaxCompletionTokens)
	}

	if saveIntoMemory {
		content := r.Choices[0].Message.Content
		memory.AddAssistantMessage(content)
	}

	return r, time.Since(startTime)
}

func (me ModelClient) chat(x Kontext, memory ChatMemory, requestedMaxCompletionTokens int64) *openai.ChatCompletion {
	params := openai.ChatCompletionNewParams{
		Messages:    openai.F(memory.ToChatCompletionMessageParamUnion()),
		Temperature: openai.F(me.config.Temperature),
		Seed:        openai.Int(1),
		Model:       openai.F(me.config.Name),
	}
	if requestedMaxCompletionTokens > 0 {
		params.MaxCompletionTokens = openai.Int(requestedMaxCompletionTokens)
	}

	r, err := me.openAiClient.Chat.Completions.New(x.Context, params)
	if err != nil {
		panic(errors.Wrap(err, "failed to call chat completions API"))
	}
	return r
}

func (me ModelClient) chatStream(x Kontext, memory ChatMemory, output io.Writer, requestedMaxCompletionTokens int64) *openai.ChatCompletion {
	params := openai.ChatCompletionNewParams{
		Messages:    openai.F(memory.ToChatCompletionMessageParamUnion()),
		Temperature: openai.F(me.config.Temperature),
		Seed:        openai.Int(0),
		Model:       openai.F(me.config.Name),
	}
	if requestedMaxCompletionTokens > 0 {
		params.MaxCompletionTokens = openai.Int(requestedMaxCompletionTokens)
	}

	stream := me.openAiStreamingClient.Chat.Completions.NewStreaming(x.Context, params)

	// optionally, an accumulator helper can be used
	acc := &openai.ChatCompletionAccumulator{}

	var buffer string // Buffer to hold partial lines

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		// if content, ok := acc.JustFinishedContent(); ok {
		// 	println("Content stream finished:", content)
		// }

		// if using tool calls
		// if tool, ok := acc.JustFinishedToolCall(); ok {
		// 	println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
		// }

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			panic(errors.New("Refusalstream finished:" + refusal))
		}

		// it's best to use chunks after handling JustFinished events
		// Process chunk content
		if len(chunk.Choices) > 0 {
			buffer += chunk.Choices[0].Delta.Content

			// Split the buffer by lines
			lines := comm.SplitBufferByLines(&buffer)
			for _, line := range lines {
				output.Write([]byte(line + "\n"))
			}
		}
	}

	// Output the last incomplete part in the buffer
	if len(buffer) > 0 {
		output.Write([]byte(buffer + "\n"))
	}

	if err := stream.Err(); err != nil {
		panic(err)
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	return &acc.ChatCompletion
}

func buildOpenAiClient(model ModelConfig, streaming bool) *openai.Client {
	options := []option.RequestOption{}
	if len(model.ApiKey) > 0 {
		options = append(options, option.WithAPIKey(model.ApiKey))
	}
	if len(model.BaseUrl) > 0 {
		options = append(options, option.WithBaseURL(model.BaseUrl))
	}

	if !streaming {
		if model.Timeout.Seconds() > 0 {
			options = append(options, option.WithRequestTimeout(model.Timeout))
		}
	}

	if len(model.ProxyUrl) > 0 {
		proxyURL, err := url.Parse(model.ProxyUrl)
		if err != nil {
			panic(fmt.Errorf("error parsing proxy URL: %+v", err))
		}

		var transport *http.Transport = nil
		if len(model.ProxyUser) == 0 {
			transport = &http.Transport{
				Proxy:           http.ProxyURL(proxyURL),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: model.ProxyInsecureSkipVerify},
			}
		} else {
			dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := &net.Dialer{}
				conn, err := dialer.DialContext(ctx, network, proxyURL.Host)
				if err != nil {
					return nil, err
				}

				auth := model.ProxyUser + ":" + model.ProxyPass
				authBase64 := base64.StdEncoding.EncodeToString([]byte(auth))

				_, err = conn.Write([]byte("CONNECT " + addr + " HTTP/1.0\r\n"))
				if err != nil {
					return nil, err
				}
				_, err = conn.Write([]byte("Proxy-Authorization: Basic " + authBase64 + "\r\n"))
				if err != nil {
					return nil, err
				}
				_, err = conn.Write([]byte("\r\n"))
				if err != nil {
					return nil, err
				}

				buffer := make([]byte, 1024)
				n, err := conn.Read(buffer)
				if err != nil {
					return nil, err
				}
				response := string(buffer[:n])
				if !strings.Contains(response, "200 Connection established") {
					return nil, fmt.Errorf("proxy connection failed: %s", response)
				}

				return conn, nil
			}

			transport = &http.Transport{
				Proxy:           http.ProxyURL(proxyURL),
				DialContext:     dialContext,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: model.ProxyInsecureSkipVerify},
			}
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		options = append(options, option.WithHTTPClient(httpClient))
	}

	return openai.NewClient(options...)
}
