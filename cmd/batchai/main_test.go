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
	// Mocking the Version and CommitId variables
	Version = "1.0.0"
	CommitId = "abcdefg"

	os.Exit(m.Run())
}

func TestAppRun(t *testing.T) {
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

	app := &cli.App{
		Version:                fmt.Sprintf("%s (%s)", Version, CommitId),
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
	c.Defaultln(`
                 Thanks for using batchai 🙏
                 Please consider starring to my work: 
               🍷  https://github.com/qiangyt/batchai`)
}

func TestAppRunWithArgs(t *testing.T) {
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

	app := &cli.App{
		Version:                fmt.Sprintf("%s (%s)", Version, CommitId),
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
	if err := app.Run([]string{"batchai", "explain"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAppRunWithInvalidArgs(t *testing.T) {
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

	app := &cli.App{
		Version:                fmt.Sprintf("%s (%s)", Version, CommitId),
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
	if err := app.Run([]string{"batchai", "invalid"}); err == nil {
		t.Errorf("Expected error, got none")
	}
}

func TestAppRunWithFlags(t *testing.T) {
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

	app := &cli.App{
		Version:                fmt.Sprintf("%s (%s)", Version, CommitId),
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
	if err := app.Run([]string{"batchai", "explain", "--inline"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAppRunWithInvalidFlags(t *testing.T) {
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
