// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

package rerror

import (
	"fmt"
	"github.com/luiscm/oro/token"
	"testing"
)

func TestError(t *testing.T) {
	errors = []string{}
	Error(Parse, token.Position{1, 1}, "Test rerror 1")
	Error(Parse, token.Position{1, 1}, "Test rerror 2")
	if len(errors) != 2 {
		t.Errorf("Expected %d but got %d", 2, len(errors))
	}
}

func TestGetErrors(t *testing.T) {
	errors = []string{}
	Error(Parse, token.Position{1, 1}, "Test rerror 1")
	Error(Runtime, token.Position{2, 1}, "Test rerror 2")
	expected := []string{
		fmt.Sprintf("%s [Line %d:%d]: %s", Parse, 1, 1, "Test rerror 1"),
		fmt.Sprintf("%s [Line %d:%d]: %s", Runtime, 2, 1, "Test rerror 2"),
	}
	for i, k := range errors {
		if k != expected[i] {
			t.Errorf("Expected %s but got %s", expected[i], k)
		}
	}
}

func TestClearErrors(t *testing.T) {
	errors = []string{}
	Error(Parse, token.Position{1, 1}, "Test rerror 1")
	Error(Parse, token.Position{1, 1}, "Test rerror 2")
	ClearErrors()
	if len(errors) != 0 {
		t.Errorf("Expected %d but got %d", 0, len(errors))
	}
}
