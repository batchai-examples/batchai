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
  "overall_severity": "trivial" or "minor" or "major" or "critical"
}`

type CheckIssueT struct {
	ShortDescription     string   `json:"short_description"`
	DetailedExplaination string   `json:"detailed_explaination"`
	Suggestion           string   `json:"suggestion"`
	IssueLineBegin       int      `json:"issue_line_begin"`
	IssueLineEnd         int      `json:"issue_line_end"`
	IssueReferenceUrls   []string `json:"issue_reference_urls"`
	Severity             string   `json:"severity"`
	SeverityReason       string   `json:"severity_reason"`
}

type CheckIssue = *CheckIssueT

func (me CheckIssue) Print(console comm.Console) {
	console.NewLine().Printf("Short Description: %s", me.ShortDescription)
	console.NewLine().Printf("Detailed Description: %s", me.DetailedExplaination)
	console.NewLine().Printf("Severity: %s", me.Severity)
	console.NewLine().Printf("Severity Reason: %s", me.SeverityReason)
	console.NewLine().Printf("Suggestion: %s", me.Suggestion)
	console.NewLine().Printf("IssueLineBegin: %d", me.IssueLineBegin)
	console.NewLine().Printf("IssueLineEnd: %d", me.IssueLineEnd)
	console.NewLine().Printf("IssueReferenceUrls: %s", me.IssueReferenceUrls)
}

type CheckReportT struct {
	HasIssue          bool              `json:"has_issue"`
	OverallSeverity   string            `json:"overall_severity"`
	Issues            []CheckIssue      `json:"issues"`
	FixedCode         string            `json:"fixed_code"`
	OriginalCode      string            `json:"original_code"`
	Path              string            `json:"path"`
	ModelUsageMetrics ModelUsageMetrics `json:"model_usage_metrics"`
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
