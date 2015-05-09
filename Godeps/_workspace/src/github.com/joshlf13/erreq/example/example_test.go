// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/joshlf13/erreq"
)

func main() {
	e1 := New("ERROR")
	e2 := New("ERROR")
	fmt.Println(e1 == e2)
	fmt.Println(e1.Equals(e2))
}

type errorString struct {
	s string
}

func New(s string) erreq.Error {
	return &errorString{s}
}

func (e1 *errorString) Equals(e2 erreq.Error) bool {
	v, ok := e2.(*errorString)
	if !ok {
		return false
	}

	return e1.s == v.s
}

func (e *errorString) Error() string {
	return e.s
}
