// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

package cmd

import (
	"github.com/luiscm/oro/token"
	"testing"
)

func TestCommandInsert(t *testing.T) {
	table = make(map[string]token.TType)
	command := &Command{}
	command.Insert("val", token.Val)
	command.Insert("repeat", token.Repeat)
	expected := 2
	if len(table) != 2 {
		t.Errorf("Expected %d but got %d", expected, len(table))
	}
}

func TestCommandLookup(t *testing.T) {
	table = make(map[string]token.TType)
	command := &Command{}
	command.Insert("val", token.Val)
	command.Insert("repeat", token.Repeat)
	tok, found := command.Lookup("repeat")
	if !found {
		t.Errorf("Expected to find a cmd but didn't.")
	}
	if tok != token.Repeat {
		t.Errorf("Expected %s but got %s.", token.Repeat, tok)
	}
}
