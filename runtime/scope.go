// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package runtime implements functions to scope.
package runtime

type Scope struct {
	store  map[string]Data
	parent *Scope
}

func NewScope() *Scope {
	return &Scope{
		store: make(map[string]Data),
	}
}

func NewScopeFrom(parent *Scope) *Scope {
	return &Scope{
		store:  make(map[string]Data),
		parent: parent,
	}
}

func (s *Scope) Read(name string) (Data, bool) {
	value, ok := s.store[name]
	if !ok && s.parent != nil {
		value, ok = s.parent.Read(name)
	}
	return value, ok
}

func (s *Scope) Write(name string, value Data) {
	s.store[name] = value
}

func (s *Scope) Update(name string, value Data) {
	s.updateParents(s, name, value)
}

func (s *Scope) updateParents(scope *Scope, name string, value Data) {
	if _, ok := scope.Read(name); ok {
		scope.Write(name, value)
	}
	if scope.parent != nil {
		s.updateParents(scope.parent, name, value)
	}
}

func (s *Scope) Merge(scope *Scope) {
	for k, v := range scope.store {
		if _, ok := s.store[k]; !ok {
			s.store[k] = v
		}
	}
}
