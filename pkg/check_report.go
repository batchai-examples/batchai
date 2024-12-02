package batchai

import (
	"bytes"
	"errors"
	"strings"

	"github.com/qiangyt/batchai/comm"

	"github.com/pkg/diff"
)

const CHECK_REPORT_JSON_FORMAT = `
{
  "has_issue": true or false,
  "issues": [
    {
	   "short_description": "...",
	   "detailed_explaination": "...",
	   "suggestion": "...",
	   "issue_line_begin": 9,
	   "issue_line_end": 12,
	   "issue_reference_urls": ["..."],
       "severity": "trivial" or "minor" or "major" or "critical",
       "severity_reason": "..."
    },
	{
	   "short_description": "...",
	   "detailed_explaination": "...",
	   "suggestion": "...",
	   "issue_line_begin": 24,
	   "issue_line_end": 29,
	   "issue_reference_urls": ["..."],
       "severity": "trivial" or "minor" or "major" or "critical",
       "severity_reason": "..."
    }
  ],
  "model_usage_metrics": {}
}`

type CheckReportT struct {
	HasIssue              bool
	Issues                []Issue
	ModelUsageMetrics     ModelUsageMetrics
}

type Issue struct {
	ShortDescription      string
	DetailedExplanation   string
	Suggestion            string
	IssueLineBegin        int
	IssueLineEnd          int
	IssueReferenceURLs  []string
	Severity              string
	SeverityReason        string
}

type CheckReport = *CheckReportT

func ExtractFixedCode(input string) (string, string) {
	begin := strings.Index(input, FIX_BEGIN_LINE)
	if begin < 0 {
		return "", input
	}
	block := input[begin+len(FIX_BEGIN_LINE):]

	end := strings.LastIndex(block, FIX_END_LINE)
	if end <= 0 {
		panic(errors.New("unmatched separator tag"))
	}
	result := block[:end]

	if strings.HasPrefix(strings.TrimSpace(result), "```") {
		result, _ = comm.ExtractMarkdownCodeBlocksP(result)
	}

	remained := input[:begin] + block[end+len(FIX_END_LINE):]
	return result, remained
}

func ExtractCheckReport(answer string, isGolang bool) CheckReport {
	jsonStr, _ := comm.ExtractMarkdownJsonBlocksP(answer)

	indexOfLeftBrace := strings.Index(jsonStr, "{")
	if indexOfLeftBrace < 0 {
		panic(errors.New("invalid json format - missing left brace"))
	}
	jsonStr = jsonStr[indexOfLeftBrace:]

	indexOfRightBrace := strings.LastIndex(jsonStr, "}")
	if indexOfRightBrace <= 0 {
		panic(errors.New("invalid json format - missing right brace"))
	}
	jsonStr = jsonStr[:indexOfRightBrace+1]

	report := &CheckReportT{}
	if err := comm.FromJson(jsonStr, false, report); err != nil {
		jsonStr = comm.FixJson(jsonStr, isGolang)
		comm.FromJsonP(jsonStr, false, report)
	}
	return report
}

func (me CheckReport) Print(console comm.Console) {
	if !me.HasIssue {
		console.NewLine().Print("no issue")
		return
	}

	console2 := console.NewIndented()

	console.NewLine().Printf("Overall severity: %s", me.OverallSeverity)

	me.ModelUsageMetrics.Print(console, comm.DEFAULT_COLOR)

	console.NewLine().Print("Total ").Yellowf("%d", len(me.Issues)).Default(" issues")
	for i, issue := range me.Issues {
		console2.Printf("\n#%d", i+1)
		issue.Print(console2)
	}

	console.NewLine().Print("Check:")
	console.NewLine()

	var buf bytes.Buffer
	if err := diff.Text("original", "fixed", me.OriginalCode, me.FixedCode, &buf); err == nil {
		for i, line := range strings.Split(buf.String(), "\n") {
			if i < 3 {
				console2.Yellowln(line)
			} else {
				if strings.HasPrefix(line, "+") {
					console2.Greenln(line)
				} else if strings.HasPrefix(line, "-") {
					console2.Redln(line)
				} else {
					console2.Defaultln(line)
				}
			}
		}
	}
}
