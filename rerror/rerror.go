// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package rerror implements functions to errors.
package rerror

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luiscm/oro/token"
)

type TError string

const (
	FoundErrors        = "Found Errors:"
	ErrorLine          = "%s [Line %d:%d]: %s"
	ParseErrors        = "Parse Errors: "
	Parse       TError = "Parse Error"
	Runtime     TError = "Runtime Error"
)

var errors []string

func Error(terror TError, location token.Position, msg string) {
	errors = append(errors, fmt.Sprintf(ErrorLine, terror, location.Row, location.Col, msg))
}

func ErrorFmt(msg string, a ...interface{}) error {
	return fmt.Errorf(msg, a...)
}

func PrintErrors() {
	color.White(FoundErrors)
	for _, e := range GetErrors() {
		color.Red(e)
	}
	ClearErrors()
}

func HasErrors() bool {
	return len(errors) > 0
}

func GetErrors() []string {
	return errors
}

func ClearErrors() {
	errors = []string{}
}
