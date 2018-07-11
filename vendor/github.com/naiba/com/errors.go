package com

//PanicIfNotNil panic if not nil
func PanicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
