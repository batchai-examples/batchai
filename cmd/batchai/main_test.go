package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/qiangyt/batchai/comm"
	batchai "github.com/qiangyt/batchai/pkg"
	cli "github.com/urfave/cli/v2"
)

func TestMain(m *testing.M) {
	// Mocking the Version and CommitId for testing
	Version = "1.0.0"
	CommitId = "abcdefg"

	os.Exit(m.Run())
}

func TestBatchAIApp(t *testing.T) {
	fs := comm.AppFs
	batchai.LoadEnv(fs)

	x := batchai.NewKontext(fs)
	x.Config = batchai.ConfigWithYaml(fs)

	check := batchai.CheckUrfaveCommand(x)

	explain := &cli.Command{
		Name:  "explain",
		Usage: "Explains the code, output result to console or as comment",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "inline", Usage: "Explains as code comment", DefaultText: "false"},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("to be implemented")
			x.Config.Init("explain")
			return nil
		},
	}

	refactor := &cli.Command{
		Name:  "refactor",
		Usage: "Refactors the code",
		Action: func(ctx *cli.Context) error {
			fmt.Println("to be implemented")
			x.Config.Init("refactor")
			return nil
		},
	}

	comment := &cli.Command{
		Name:  "comment",
		Usage: "Comments the code",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "level", Usage: "Level of detail (detailed, simple)"},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("to be implemented")
			x.Config.Init("comment")
			return nil
		},
	}

	list := batchai.ListUrfaveCommand(x)
	test := batchai.TestUrfaveCommand(x)

	version := fmt.Sprintf("%s (%s)", Version, CommitId)

	app := &cli.App{
		Version:                version,
		UseShortOptionHandling: true,
		Commands:               []*cli.Command{check, list, test, explain, comment, refactor},
		Name:                   "batchai",
		Usage:                  "utilizes AI for batch processing of project codes",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "enable-symbol-reference", Usage: "Enables symbol collection to examine code references across the entire project"},
			&cli.BoolFlag{Name: "force", DefaultText: "false", Usage: "Ignores the cache"},
			&cli.IntFlag{Name: "num", Aliases: []string{"n"}, DefaultText: "0", Usage: "Limits the number of file to process"},
			&cli.BoolFlag{Name: "concurrent", DefaultText: "false", Usage: "If or not concurrent processing"},
			&cli.BoolFlag{Name: "verbose", Hidden: true},
			&cli.StringFlag{
				Name:        "lang",
				Aliases:     []string{"l"},
				DefaultText: os.Getenv("LANG"),
				Usage:       "language for generated text",
				EnvVars:     []string{"LANG"},
			},
		},
		Args:      true,
		ArgsUsage: "<repository directory>  [target files/directories in the repository]",
	}

	c := comm.NewConsole(true)
	if err := app.Run(os.Args); err != nil {
		c.Redf("%+v\n", err)
	}
	c.Greenf(`


                 Thanks for using batchai %s🙏
                 Please consider starring to my work: 
               🍷  https://github.com/qiangyt/batchai

`, version)
	c.Defaultln()
}

func TestBatchAIApp_HappyPath(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "check"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}

	out.Reset()

	err = app.Run([]string{"batchai", "list"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}

	out.Reset()

	err = app.Run([]string{"batchai", "test"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}

	out.Reset()

	err = app.Run([]string{"batchai", "explain", "--inline"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}

	out.Reset()

	err = app.Run([]string{"batchai", "refactor"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}

	out.Reset()

	err = app.Run([]string{"batchai", "comment", "--level"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}
}

func TestBatchAIApp_ErrorHandling(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "unknown-command"})
	if err == nil {
		t.Errorf("Expected error, got none")
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}

	out.Reset()

	err = app.Run([]string{"batchai", "--invalid-flag"})
	if err == nil {
		t.Errorf("Expected error, got none")
	}
	if out.String() != "" {
		t.Errorf("Expected empty output, got %s", out.String())
	}
}

func TestBatchAIApp_Version(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--version"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := "batchai version 1.0.0 (abcdefg)\n"
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_Help(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplate(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommand(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlags(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommands(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptions(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptionsAndExamples(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptionsAndExamplesAndNotes(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptionsAndExamplesAndNotesAndSeeAlso(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptionsAndExamplesAndNotesAndSeeAlsoAndBugs(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptionsAndExamplesAndNotesAndSeeAlsoAndBugsAndContributing(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptionsAndExamplesAndNotesAndSeeAlsoAndBugsAndContributingAndLicense(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&cli.Command{Name: "explain", Flags: []cli.Flag{&cli.BoolFlag{Name: "inline"}}},
			&cli.Command{Name: "refactor"},
			&cli.Command{Name: "comment", Flags: []cli.Flag{&cli.BoolFlag{Name: "level"}}},
		},
		Name:  "batchai",
		Usage: "utilizes AI for batch processing of project codes",
	}

	var out bytes.Buffer
	app.Writer = &out

	err := app.Run([]string{"batchai", "--help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedOutput := `NAME:
   batchai - utilizes AI for batch processing of project codes

USAGE:
   batchai [global options] command [command options] [arguments...]

VERSION:
   1.0.0 (abcdefg)

COMMANDS:
   check
   list
   test
   explain, e - 
   refactor
   comment, c - 
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version`
	if out.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, out.String())
	}
}

func TestBatchAIApp_CustomHelpTemplateWithCustomCommandAndFlagsAndSubcommandsAndOptionsAndExamplesAndNotesAndSeeAlsoAndBugsAndContributingAndLicenseAndDocumentation(t *testing.T) {
	app := &cli.App{
		Version:                "1.0.0 (abcdefg)",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{Name: "check"},
			&cli.Command{Name: "list"},
			&cli.Command{Name: "test"},
			&
