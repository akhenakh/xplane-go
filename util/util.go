package util

// #cgo CFLAGS: -DXPLM410=1
// #include <stdlib.h>
// #include "XPLMUtilities.h"
import "C"
import "unsafe"

// DebugString writes a string to the X-Plane Log.txt file.
func DebugString(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs)) // Now C.free is recognized
	C.XPLMDebugString(cs)
}

// GetSystemPath returns the full path to the X-Plane installation directory.
func GetSystemPath() string {
	buffer := make([]byte, 512)
	C.XPLMGetSystemPath((*C.char)(unsafe.Pointer(&buffer[0])))
	return C.GoString((*C.char)(unsafe.Pointer(&buffer[0])))
}
