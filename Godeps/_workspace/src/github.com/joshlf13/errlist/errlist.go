// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errlist contains a type compatible with the error interface which handles lists of
// errors. All of the methods in this package are nil-safe; that is, calling them on nil
// pointers is expected behavior. errlist implements the erreq interface 
// (github.com/joshlf13/erreq) which allows two lists to be checked for equality.
//
// Note that functions in this package which return Errlist pointers do not do so merely
// as a convenience. You should not expect that calling such functions on a pointer
// will necessarily modify the object pointed at - it may return a new pointer.
// For example, many functions treat a nil pointer as a valid list (ie, it is possible to
// append to a nil list).
package errlist

import (
	"errors"
	"github.com/joshlf13/erreq"
)

// An error type which supports lists of errors.
type Errlist struct {
	hd  *errnode
	tl  *errnode
	num int
}

type errnode struct {
	error
	next *errnode
}

// Create a new empty error list
// (returns nil)
func EmptyList() *Errlist {
	return nil
}

// Create a new error list starting with
// an error created from e. If an empty
// error string is provided, a nil
// pointer is returned.
func NewString(e string) *Errlist {
	if e == "" {
		return nil
	}
	var erl Errlist
	erl.hd = &errnode{errors.New(e), nil}
	erl.tl = erl.hd
	erl.num = 1
	return &erl
}

// Create a new error list starting
// with e. If a nil error is provided,
// a nil pointer is returned.
func NewError(e error) *Errlist {
	if e == nil {
		return nil
	}
	var erl Errlist
	erl.hd = &errnode{e, nil}
	erl.tl = erl.hd
	erl.num = 1
	return &erl
}

// Create an error from e and append 
// it to the error list, or if the
// list is nil, create a new list
// with the error as its first element.
// In either case, return the resultant list.
// If e is an empty error string, do not
// append an error. If AddString was called
// on a nil list and with an emtpy string
// as the argument, it returns a nil pointer.
func (erl *Errlist) AddString(e string) *Errlist {
	if erl == nil {
		return NewString(e)
	}
	if e == "" {
		return erl
	}
	ern := new(errnode)
	ern.error = errors.New(e)
	erl.tl.next = ern
	erl.tl = ern
	erl.num++
	return erl
}

// Append e to the error list,
// or if the list is nil, create
// a new list with e as its first
// element. In either case, return
// the resultant list. If e is nil,
// do not append an error. If
// AddError was called on a nil
// list and with a nil argument,
// it returns a nil pointer.
func (erl *Errlist) AddError(e error) *Errlist {
	if erl == nil {
		return NewError(e)
	}
	if e == nil {
		return erl
	}
	ern := new(errnode)
	ern.error = e
	erl.tl.next = ern
	erl.tl = ern
	erl.num++
	return erl
}

// Return a string consisting of
// each error in the list printed
// and separated by newlines, or
// an empty string if called on
// a nil pointer.
func (erl *Errlist) Error() string {
	out := ""
	if erl == nil {
		return out
	}
	for n := erl.hd; n != nil; n = n.next {
		out += n.error.Error() + "\n"
	}
	return out[:len(out)-1]
}

// Return the errors as a slice.
// If called on a nil pointer,
// returns an empty slice.
func (erl *Errlist) Slice() []error {
	if erl == nil {
		return make([]error, 0)
	}
	esl := make([]error, erl.num)
	for i, n := 0, erl.hd; i < erl.num; i, n = (i + 1), n.next {
		esl[i] = n.error
	}
	return esl
}

// Creates an error list whose
// elements are the elements
// of the argument slice.
func FromSlice(e []error) *Errlist {
	erl := EmptyList()

	for _, v := range e {
		erl = erl.AddError(v)
	}
	return erl
}

// Return the number of errors
// in the list, or 0 if called
// on a nil pointer.
func (erl *Errlist) Num() int {
	if erl == nil {
		return 0
	}
	return erl.num
}

// Err returns an error equivalent 
// to this error list. If the list 
// is empty, Err returns nil.
// This is meant primarily for
// checking against nil values,
// since interface types with
// nil values are not equal to nil.
func (erl *Errlist) Err() error {
	if erl == nil {
		return nil
	}
	if erl.num == 1 {
		// This is ugly and seemingly pointless,
		// but doing it directly was giving me
		// cryptic errors from other packages
		// which imported this one.
		n := erl.hd
		return n.error
	}
	return erl
}

// errlist implements the erreq interface
// (github.com/joshlf13/erreq). Equals
// checks for pairwise equality between
// two lists. Equality is determined by
// first checking to see if both elements
// of the pair implement the erreq interface.
// If they do, erreq.Equals is used.
// Otherwise, pointer equality is used.
// Lists of different length are never equal. 
// Two nil lists are always equal. A nil list 
// is never equal to a non-nil list.
func (erl1 *Errlist) Equals(e erreq.Error) bool {
	erl2, ok := e.(*Errlist)
	if !ok {
		return false
	}

	if erl1 == erl2 {
		return true
	}

	// Since empty lists are always nil pointers,
	// and since we have already checked for
	// pointer equality, this check is now valid.
	// Plus, we avoid nil pointer dereferences.
	if erl1 == nil || erl2 == nil {
		return false
	}

	if erl1.num != erl2.num {
		return false
	}

	ern1 := erl1.hd
	ern2 := erl2.hd
	for ern1 != nil && ern2 != nil {
		erreq1, ok1 := ern1.error.(erreq.Error)
		erreq2, ok2 := ern2.error.(erreq.Error)

		if (ok1 && ok2) && (!erreq1.Equals(erreq2)) {
			return false
		} else if ern1.error != ern2.error {
			return false
		}

		ern1 = ern1.next
		ern2 = ern2.next
	}

	// In case the Errlist.num field
	// is inaccurate
	return ern1 == ern2
}
