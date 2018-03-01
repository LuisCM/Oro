// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

package ast

import (
	"github.com/luiscm/oro/token"
	"testing"
)

func TestVal(t *testing.T) {
	val := &Val{
		Token: token.Token{Type: token.Val, Literal: "val"},
		Name: &Identifier{
			Token: token.Token{Type: token.Identifier, Literal: "myVal"},
			Value: "myVal",
		},
		Value: &Identifier{
			Token: token.Token{Type: token.Identifier, Literal: "anotherVal"},
			Value: "anotherVal",
		},
	}
	if val.Check() != "val myVal = anotherVal" {
		t.Errorf("val.Check() wrong. got=%q", val.Check())
	}
}
