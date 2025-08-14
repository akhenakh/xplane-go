package widget

// #cgo CFLAGS: -DXPLM410=1
// #cgo LDFLAGS: -lXPWidgets_64
// #include <stdlib.h>
// #include "XPWidgets.h"
import "C"

import "unsafe"

type WidgetID C.XPWidgetID
type WidgetClass int
type WidgetMessage int
type PropertyID int
type DispatchMode int

// WidgetFunc is the callback for custom widget behavior.
// Return true if the message was handled.
type WidgetFunc func(message WidgetMessage, widget WidgetID, param1, param2 int) bool

// CreateWidget creates a new widget.
func CreateWidget(left, top, right, bottom int, visible bool, desc string, isRoot bool, container WidgetID, class WidgetClass) WidgetID {
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	vis := 0
	if visible {
		vis = 1
	}
	root := 0
	if isRoot {
		root = 1
	}

	return WidgetID(C.XPCreateWidget(
		C.int(left), C.int(top), C.int(right), C.int(bottom),
		C.int(vis),
		cDesc,
		C.int(root),
		C.XPWidgetID(container),
		C.XPWidgetClass(class),
	))
}

// DestroyWidget destroys a widget and optionally its children.
func DestroyWidget(id WidgetID, destroyChildren bool) {
	dc := 0
	if destroyChildren {
		dc = 1
	}
	C.XPDestroyWidget(C.XPWidgetID(id), C.int(dc))
}

// SetWidgetDescriptor sets the text associated with a widget.
func SetWidgetDescriptor(id WidgetID, desc string) {
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	C.XPSetWidgetDescriptor(C.XPWidgetID(id), cDesc)
}

// GetWidgetDescriptor gets the text associated with a widget.
func GetWidgetDescriptor(id WidgetID) string {
	buf := make([]byte, 1024)
	len := C.XPGetWidgetDescriptor(C.XPWidgetID(id), (*C.char)(unsafe.Pointer(&buf[0])), 1024)
	return string(buf[:len])
}

// GetWidgetProperty retrieves a property value from a widget.
func GetWidgetProperty(id WidgetID, propID PropertyID) (value int, exists bool) {
	var ex C.int
	val := C.XPGetWidgetProperty(C.XPWidgetID(id), C.XPWidgetPropertyID(propID), &ex)
	return int(val), ex != 0
}

// SetWidgetProperty sets a property value on a widget.
func SetWidgetProperty(id WidgetID, propID PropertyID, value int) {
	C.XPSetWidgetProperty(C.XPWidgetID(id), C.XPWidgetPropertyID(propID), C.long(value))
}

// ShowWidget makes a widget visible.
func ShowWidget(id WidgetID) {
	C.XPShowWidget(C.XPWidgetID(id))
}

// HideWidget makes a widget invisible.
func HideWidget(id WidgetID) {
	C.XPHideWidget(C.XPWidgetID(id))
}

// IsWidgetVisible checks if a widget and its ancestors are visible.
func IsWidgetVisible(id WidgetID) bool {
	return C.XPIsWidgetVisible(C.XPWidgetID(id)) != 0
}
