package sqlutil

import (
	"reflect"
	"unsafe"
)

// unsafeStringToBytes converts string to slice without copy.
// Use at your own risk.
func unsafeStringToBytes(s string) (b []byte) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return
}
