package plugin

// #cgo CFLAGS: -DXPLM410=1
// #include <string.h>
// #include <stdlib.h>
// #include "XPLMPlugin.h"
//
// extern int XPluginStart(char* outName, char* outSig, char* outDesc);
// extern void XPluginStop(void);
// extern void XPluginDisable(void);
// extern int XPluginEnable(void);
// extern void XPluginReceiveMessage(XPLMPluginID inFrom, int inMsg, void* inParam);
import "C"

import (
	"unsafe"

	"github.com/akhenakh/xplane-go/util"
)

// Plugin is the primary interface that developers must implement to create an
// X-Plane plugin. It mirrors the five required C functions.
type Plugin interface {
	Start() (name, sig, desc string, err error)
	Stop()
	Enable() error
	Disable()
}

// MessageHandler is an interface for plugins that want to receive inter-plugin
// broadcast messages.
type MessageHandler interface {
	ReceiveMessage(from PluginID, msg Message, param unsafe.Pointer)
}

type PluginID int
type Message int

var (
	// The single instance of the user-provided plugin.
	pluginImpl Plugin
)

// Register registers the user's plugin implementation. This function must be
// called from an `init()` function in the plugin's `main` package.
func Register(p Plugin) {
	if pluginImpl != nil {
		util.DebugString("xplane-go: Plugin already registered.\n")
		return
	}
	pluginImpl = p
}

//export XPluginStart
func XPluginStart(outName, outSig, outDesc *C.char) C.int {
	if pluginImpl == nil {
		// If no plugin was registered, we can't start.
		// Copy a message to X-Plane's buffer so the user knows why.
		msg := "xplane-go: No plugin registered. Call plugin.Register() from init()."
		C.strncpy(outName, C.CString(msg), 255)
		C.strncpy(outSig, C.CString("xplane-go.error.no-register"), 255)
		C.strncpy(outDesc, C.CString(msg), 255)
		return 0
	}

	name, sig, desc, err := pluginImpl.Start()
	if err != nil {
		errMsg := "xplane-go: plugin start failed: " + err.Error() + "\n"
		util.DebugString(errMsg)
		C.strncpy(outName, C.CString("Error"), 255)
		C.strncpy(outSig, C.CString("xplane-go.error.start"), 255)
		C.strncpy(outDesc, C.CString(errMsg), 255)
		return 0
	}

	// Copy the plugin info into the C buffers provided by X-Plane.
	cName := C.CString(name)
	cSig := C.CString(sig)
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cSig))
	defer C.free(unsafe.Pointer(cDesc))

	C.strncpy(outName, cName, 255)
	C.strncpy(outSig, cSig, 255)
	C.strncpy(outDesc, cDesc, 255)

	return 1
}

//export XPluginStop
func XPluginStop() {
	if pluginImpl != nil {
		pluginImpl.Stop()
	}
}

//export XPluginEnable
func XPluginEnable() C.int {
	if pluginImpl != nil {
		if err := pluginImpl.Enable(); err != nil {
			util.DebugString("xplane-go: plugin enable failed: " + err.Error() + "\n")
			return 0
		}
		return 1
	}
	return 0
}

//export XPluginDisable
func XPluginDisable() {
	if pluginImpl != nil {
		pluginImpl.Disable()
	}
}

//export XPluginReceiveMessage
func XPluginReceiveMessage(inFrom C.XPLMPluginID, inMsg C.int, inParam unsafe.Pointer) {
	if handler, ok := pluginImpl.(MessageHandler); ok {
		handler.ReceiveMessage(PluginID(inFrom), Message(inMsg), inParam)
	}
}
