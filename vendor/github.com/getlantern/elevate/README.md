[Godoc](http://godoc.org/github.com/getlantern/elevate)

elevate currently only works for OS X and Windows. The Windows support
currently uses a Visual Basic script that ends up displaying a confusing prompt
and is generally hoaky - it will be replaced by a C++ program that does the same
thing but with a better prompt.

On OS X, it uses cocoasudo from here - https://github.com/getlantern/cocoasudo,
forked from https://github.com/kalikaneko/cocoasudo to explicitly support OSX
10.6.

On Windows, it uses elevate from here - http://code.kliu.org/misc/elevate/. The
source code lives in elevate-1.3.0/src and can be built from a Visual Studio
command line by running `nmake elevat.mak`. The elevate makefile has been
modified from the original to 1. build as a windows GUI app instead of a console
app, 2. always build as 32 bit and 3. statically link the runtime. `elevate.c`
was modified to include a WinMain function.
