keyman [![Travis CI Status](https://travis-ci.org/getlantern/keyman.svg?branch=master)](https://travis-ci.org/getlantern/keyman)&nbsp;[![Coverage Status](https://coveralls.io/repos/getlantern/keyman/badge.png)](https://coveralls.io/r/getlantern/keyman)&nbsp;[![GoDoc](https://godoc.org/github.com/getlantern/keyman?status.png)](http://godoc.org/github.com/getlantern/keyman)
======

Easy golang RSA key and certificate management.

API documentation available on [godoc](https://godoc.org/github.com/getlantern/keyman).

### Build Notes

On Windows, keyman uses a custom executable for importing certificates into the
system trust store.  This executable is built using Visual Studio from this
[solution](certimporter).

The resulting executable is packaged into go using `embedbinaries.bash`.
