errlist/example
===============

This test creates a number of files, and then attempts to open those files and some others, with some system calls returning errors and others returning no error.  Each error is added to a list of errors, and printed.

###Use

To use, remove the "_test.go" extension before building (the suffix prevents the "go get" command from installing this example directory).