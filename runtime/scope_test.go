// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

package runtime

import (
	"testing"
)

func TestScopeReadWrite(t *testing.T) {
	s := NewScope()
	s.Write("num", &TInteger{Value: 20})
	val, ok := s.Read("num")
	if !ok {
		t.Errorf("Expected a value but got nothing")
	}
	value, ok := val.(*TInteger)
	if !ok {
		t.Errorf("Expected an Integer type but got %T", val)
	}
	if value.Value != 20 {
		t.Errorf("Expected %d but got %d", 20, value.Value)
	}
}

func TestScopeParent(t *testing.T) {
	sp := NewScope()
	sp.Write("num", &TInteger{Value: 20})
	s := NewScopeFrom(sp)
	val, ok := s.Read("num")
	if !ok {
		t.Errorf("Expected a value but got nothing")
	}
	value, ok := val.(*TInteger)
	if !ok {
		t.Errorf("Expected an Integer type but got %T", val)
	}
	if value.Value != 20 {
		t.Errorf("Expected %d but got %d", 20, value.Value)
	}
}

func TestScopeUpdate(t *testing.T) {
	sp := NewScope()
	sp.Write("num", &TInteger{Value: 20})
	s := NewScopeFrom(sp)
	s.Update("num", &TInteger{Value: 30})
	val, ok := s.Read("num")
	valP, okP := sp.Read("num")
	if !ok || !okP {
		t.Errorf("Expected a value but got nothing")
	}
	value, ok := val.(*TInteger)
	valueP, okP := valP.(*TInteger)
	if !ok || !okP {
		t.Errorf("Expected an Integer type but got %T", val)
	}
	if value.Value != 30 || valueP.Value != 30 {
		t.Errorf("Expected %d but got %d", 30, value.Value)
	}
}

func TestScopeMerge(t *testing.T) {
	sp := NewScope()
	sp.Write("num", &TInteger{Value: 20})
	s := NewScope()
	s.Write("str", &TString{Value: "test"})
	s.Merge(sp)
	val, ok := s.Read("num")
	if !ok {
		t.Errorf("Expected a value but got nothing")
	}
	value, ok := val.(*TInteger)
	if !ok {
		t.Errorf("Expected an Integer type but got %T", val)
	}
	if value.Value != 20 {
		t.Errorf("Expected %d but got %d", 20, value.Value)
	}
}
