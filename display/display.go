package widget

// #cgo CFLAGS: -DXPLM410=1
// #cgo LDFLAGS: -lXPWidgets_64
// #include <stdlib.h>
// #include "XPWidgets.h"
//
// extern int widgetCallback_cgo(XPWidgetMessage inMessage, XPWidgetID inWidget, intptr_t inParam1, intptr_t inParam2);
import "C"

import (
	"sync"
	"unsafe"
)

type WidgetID C.XPWidgetID
type WidgetClass int
type WidgetMessage int
type PropertyID int
type DispatchMode int

// WidgetFunc is the callback for custom widget behavior.
// It receives the message, the widget that received it, and two message-specific parameters.
// Return true (1) if you handled the message, or false (0) to pass it to a parent widget.
type WidgetFunc func(message WidgetMessage, widget WidgetID, param1, param2 int) bool

var (
	widgetCallbackRegistry         = make(map[uintptr]WidgetFunc)
	nextCallbackID         uintptr = 1
	registryMutex          sync.RWMutex
)

func registerWidgetCallback(callback WidgetFunc) uintptr {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	id := nextCallbackID
	nextCallbackID++
	widgetCallbackRegistry[id] = callback
	return id
}

func getWidgetCallback(id uintptr) WidgetFunc {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return widgetCallbackRegistry[id]
}

//export widgetCallback_cgo
func widgetCallback_cgo(inMessage C.XPWidgetMessage, inWidget C.XPWidgetID, inParam1 C.intptr_t, inParam2 C.intptr_t) C.int {
	id := uintptr(C.XPGetWidgetProperty(inWidget, C.xPProperty_Refcon, nil))
	if callback := getWidgetCallback(id); callback != nil {
		handled := callback(
			WidgetMessage(inMessage),
			WidgetID(inWidget),
			int(inParam1),
			int(inParam2),
		)
		if handled {
			return 1
		}
	}
	return 0
}

// AddWidgetCallback attaches a callback function to a widget.
// The callback will be executed when the widget receives messages.
func AddWidgetCallback(id WidgetID, callback WidgetFunc) {
	callbackID := registerWidgetCallback(callback)
	// We store our Go callback's ID in the widget's 'refcon' property.
	// The C trampoline function will use this to find the correct Go func.
	SetWidgetProperty(id, PropertyRefcon, int(callbackID))
	C.XPAddWidgetCallback(C.XPWidgetID(id), C.widgetCallback_cgo)
}

// GetWidgetClass returns the class of a given widget.
func GetWidgetClass(id WidgetID) WidgetClass {
	return WidgetClass(C.XPGetWidgetClass(C.XPWidgetID(id)))
}

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
	length := C.XPGetWidgetDescriptor(C.XPWidgetID(id), (*C.char)(unsafe.Pointer(&buf[0])), 1024)
	if length <= 0 {
		return ""
	}
	return string(buf[:length])
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
