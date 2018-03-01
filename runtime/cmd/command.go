// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package cmd implements functions to commands.
package cmd

import (
	"github.com/luiscm/oro/token"
)

type Command struct{}

var table = make(map[string]token.TType)

func (s *Command) Insert(name string, t token.TType) {
	table[name] = t
}

func (s *Command) InsertAll() {
	table[token.True] = token.Boolean
	table[token.False] = token.Boolean
	table[token.Nil] = token.Nil
	table[token.Val] = token.Val
	table[token.Var] = token.Var
	table[token.Function] = token.Function
	table[token.Do] = token.Do
	table[token.End] = token.End
	table[token.If] = token.If
	table[token.Else] = token.Else
	table[token.Repeat] = token.Repeat
	table[token.In] = token.In
	table[token.Is] = token.Is
	table[token.As] = token.As
	table[token.Return] = token.Return
	table[token.Then] = token.Then
	table[token.Match] = token.Match
	table[token.When] = token.When
	table[token.With] = token.With
	table[token.Break] = token.Break
	table[token.Continue] = token.Continue
	table[token.Module] = token.Module
	table[token.Use] = token.Use
}

func (s *Command) Lookup(name string) (token.TType, bool) {
	if tok, ok := table[name]; ok {
		return tok, true
	}
	return "", false
}
