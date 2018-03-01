// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package main implements functions to main application.
package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/luiscm/oro/interpreter"
	"github.com/luiscm/oro/lexer"
	"github.com/luiscm/oro/parser"
	"github.com/luiscm/oro/rerror"
	"github.com/luiscm/oro/runtime"
	"github.com/luiscm/oro/util"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = util.Name()
	app.Usage = ""
	app.Authors = []cli.Author{{
		Name:  util.AuthorName(),
		Email: util.AuthorEmail(),
	}}
	app.Version = util.Version()
	app.Compiled = time.Now()
	app.Copyright = fmt.Sprintf(util.Copyright(), time.Now().Year())
	app.Commands = []cli.Command{
		{
			Name:  util.CliCommandNameRun(),
			Usage: util.CliCommandUsageRun(),
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 1 {
					color.Red(util.CliCommandActionRunSourceFile())
					return nil
				}
				file := c.Args()[0]
				ext := filepath.Ext(file)
				if ext == "" || ext != util.FileExtension() {
					color.Red(util.CliCommandActionRunExistFile(), file)
					return nil
				}
				source, err := ioutil.ReadFile(file)
				if err != nil {
					color.Red(util.CliCommandActionRunReadFile(), file)
					return nil
				}
				lex := lexer.New(source)
				if rerror.HasErrors() {
					rerror.PrintErrors()
					return nil
				}
				parse := parser.New(lex)
				program := parse.Parse()
				if rerror.HasErrors() {
					rerror.PrintErrors()
					return nil
				}
				runner := interpreter.New()
				runner.Interpreter(program, runtime.NewScope())
				if rerror.HasErrors() {
					rerror.PrintErrors()
					return nil
				}
				return nil
			},
		},
		{
			Name:  util.CliCommandNameRepl(),
			Usage: util.CliCommandUsageRepl(),
			Action: func(c *cli.Context) error {
				input := bufio.NewReader(os.Stdin)
				color.HiGreen(util.NameVersionEnvironment())
				color.HiBlue(util.CommandExit())
				sc := runtime.NewScope()
				for {
					color.Set(color.FgHiWhite)
					fmt.Print(util.ReplSignal())
					color.Unset()
					source, _ := input.ReadBytes('\n')
					lex := lexer.New(source)
					if rerror.HasErrors() {
						rerror.PrintErrors()
						continue
					}
					parse := parser.New(lex)
					program := parse.Parse()
					if rerror.HasErrors() {
						rerror.PrintErrors()
						continue
					}
					runner := interpreter.New()
					object := runner.Interpreter(program, sc)
					if rerror.HasErrors() {
						rerror.PrintErrors()
						continue
					}
					if object != nil {
						fmt.Println(object.Check())
					}
				}
			},
		},
	}
	app.CommandNotFound = func(ctx *cli.Context, command string) {
		color.Set(color.FgHiRed)
		fmt.Fprintf(ctx.App.Writer, util.CommandNotFound(), command)
		color.Unset()
	}
	app.Run(os.Args)
}
