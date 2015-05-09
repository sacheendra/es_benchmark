// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"github.com/joshlf13/erreq"
	"github.com/joshlf13/errlist"
)

type equalError struct {
	i int
}

func (e1 *equalError) Equals(e erreq.Error) bool {
	e2, ok := e.(*equalError)

	if !ok {
		return false
	}

	return e1.i == e2.i
}

func (e *equalError) Error() string {
	return fmt.Sprintf("%d", e.i)
}

func main() {
	var erl *errlist.Errlist
	// var erq *erreq.Error
	// var err error

	fmt.Println("First we will create an empty list")
	// erl = errlist.EmptyList()
	erl = nil

	fmt.Printf("It has length %d\n", erl.Num())
	fmt.Println()

	fmt.Println("Then we will add 100 nil errors just for fun")

	for i := 0; i < 100; i++ {
		erl = erl.AddError(nil)
	}

	fmt.Printf("Now it has length %d\n", erl.Num())
	fmt.Println()

	fmt.Println("Now we will add 10 real errors")

	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%d", i)

		if i%2 == 0 {
			erl = erl.AddString(str)
		} else {
			erl = erl.AddError(errors.New(str))
		}
	}

	fmt.Println("...and print them all out")
	fmt.Println("First using errlist.Error():")
	fmt.Println(erl)
	fmt.Println("Then by getting the errors as a slice:")

	sl := erl.Slice()
	for _, e := range sl {
		fmt.Println(e)
	}
	fmt.Println()

	fmt.Println("Now we will make another list and check that they are equal")

	var erl2 *errlist.Errlist
	erl2 = nil

	erl2 = errlist.FromSlice(sl)

	if erl.Equals(erl2) {
		fmt.Println("They are! Yay!")
	} else {
		fmt.Println("Oh noes, they are not equal!")
	}

	fmt.Println()
	fmt.Println("Now let's try it with errors which can be checked for equality")
	fmt.Println("Just to make it tricky, we'll make sure that the pointers won't be equal. Muahaha!")

	e1 := errlist.EmptyList()
	e2 := errlist.EmptyList()

	for i := 0; i < 10; i++ {
		e1 = e1.AddError(&equalError{i})
		e2 = e2.AddError(&equalError{i})
	}

	if e1.Equals(e2) {
		fmt.Println("They are! Yay!")
	} else {
		fmt.Println("Oh noes, they are not equal!")
	}

	fmt.Println("And check for equality...")
}
