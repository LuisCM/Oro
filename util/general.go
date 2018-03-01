// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package util implements functions to environment.
package util

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	OroName                          = "Oro Programming Language"
	OroVersion                       = "1.0.0"
	OroReplSignal                    = "oro> "
	OroAuthorName                    = "LuisCM"
	OroAuthorEmail                   = "tcljava@gmail.com"
	OroCopyrightDescription          = "Copyright "
	OroCopyright                     = "\u00a9 2011-%d LuisCM All Rights Reserved."
	OroFileExtension                 = ".oro"
	OroCmdNotFound                   = "Command %q doesn't exist.\n"
	OroCliCommandExit                = "Use CTRL+C or quit() return to exit."
	OroCliCommandNameRun             = "run"
	OroCliCommandUsageRun            = "Run an source file."
	OroCliCommandActionRunExistFile  = "The file extension '%s' should be '" + OroFileExtension + "'."
	OroCliCommandActionRunSourceFile = "Run expects a source file as argument."
	OroCliCommandActionRunReadFile   = "Couldn't read '%s'."
	OroCliCommandNameRepl            = "repl"
	OroCliCommandUsageRepl           = "Start the interactive Read-Eval-Print Loop."
)

func Environment() string {
	return fmt.Sprintf("%s/%s %d-CPU(s) %s", strings.ToUpper(runtime.GOOS), strings.ToUpper(runtime.GOARCH), runtime.NumCPU(), time.Now().Format("2006-01-02 15:04:05 Monday"))
}

func NameVersionEnvironment() string {
	return Name() + " " + Version() + " " + Environment()
}

func ReplSignal() string {
	return OroReplSignal
}

func Name() string {
	return OroName
}

func Version() string {
	return OroVersion
}

func NameVersion() string {
	return Name() + " " + Version()
}

func AuthorName() string {
	return OroAuthorName
}

func AuthorEmail() string {
	return OroAuthorEmail
}

func Copyright() string {
	return OroCopyright
}

func CopyrightDescription() string {
	return OroCopyrightDescription + fmt.Sprintf(OroCopyright, time.Now().Year())
}

func CommandNotFound() string {
	return OroCmdNotFound
}

func FileExtension() string {
	return OroFileExtension
}

func CommandExit() string {
	return OroCliCommandExit
}

func CliCommandNameRun() string {
	return OroCliCommandNameRun
}

func CliCommandUsageRun() string {
	return OroCliCommandUsageRun
}

func CliCommandActionRunExistFile() string {
	return OroCliCommandActionRunExistFile
}

func CliCommandActionRunSourceFile() string {
	return OroCliCommandActionRunSourceFile
}

func CliCommandActionRunReadFile() string {
	return OroCliCommandActionRunReadFile
}

func CliCommandNameRepl() string {
	return OroCliCommandNameRepl
}

func CliCommandUsageRepl() string {
	return OroCliCommandUsageRepl
}
