package store
import "C"
import (
	"reflect"
	"unsafe"
)

type SliceHeader struct {
	Data uintptr
	Len  int
}

func (sh *SliceHeader) ToBytes() (b []byte) {
	sb := (*reflect.SliceHeader)((unsafe.Pointer(&b)))
	sb.Data = sh.Data
	sb.Cap = sh.Len
	sb.Len = sh.Len
	return
}

func (sh *SliceHeader) enlarge(size int) {
	if sh.Len != 0 {
		sh.Data = uintptr(C.realloc(unsafe.Pointer(sh.Data), C.size_t(size)))
	} else {
		sh.Data = uintptr(C.malloc(C.size_t(size)))
	}
	sh.Len = size
}