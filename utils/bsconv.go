package utils

import (
	"reflect"
	"unsafe"
)

// Bcs converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func Bcs(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

// BcsChar  converts a byte to string without memory allocation.
func BcsChar(b byte) (s string) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Len = 1
	sh.Data = uintptr(unsafe.Pointer(&b))
	return
}

// Scb converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func Scb(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
