// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The erreq package defines an extension of the error interface
// which supports equality checking
package erreq

type Error interface {
	error
	Equals(Error) bool
}
