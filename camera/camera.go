package camera

// #cgo CFLAGS: -DXPLM410=1
// #include <stdlib.h>
// #include "XPLMCamera.h"
//
// extern int cameraControl_cgo(XPLMCameraPosition_t* outCameraPosition, int inIsLosingControl, void* inRefcon);
import "C"

import (
	"sync"
	"unsafe"
)

// Position holds all data for a camera's location and orientation.
type Position struct {
	X, Y, Z float32
	Pitch   float32
	Heading float32
	Roll    float32
	Zoom    float32
}

// CameraControlFunc is the signature for a callback that controls the camera.
// It is called every frame.
// - isLosingControl: True if another part of X-Plane is taking control.
// - Returns:
//   - keepControl: True if you want to continue controlling the camera.
//   - newPos: The desired new position, or nil to not change the position.
type CameraControlFunc func(isLosingControl bool) (keepControl bool, newPos *Position)

// ControlDuration determines how long the plugin retains camera control.
type ControlDuration int

const (
	UntilViewChanges ControlDuration = C.xplm_ControlCameraUntilViewChanges
	Forever          ControlDuration = C.xplm_ControlCameraForever
)

var (
	registry      = make(map[uintptr]CameraControlFunc)
	registryMutex sync.RWMutex
	currentID     uintptr
)

func registerCallback(callback CameraControlFunc) uintptr {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	// For camera control, only one plugin can have it at a time.
	// We'll manage just one callback.
	currentID++
	registry[currentID] = callback
	return currentID
}

func unregisterCallback(id uintptr) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	delete(registry, id)
}

func getCallback(id uintptr) CameraControlFunc {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return registry[id]
}

//export cameraControl_cgo
func cameraControl_cgo(outCameraPosition *C.XPLMCameraPosition_t, inIsLosingControl C.int, inRefcon unsafe.Pointer) C.int {
	id := uintptr(inRefcon)
	callback := getCallback(id)
	if callback == nil {
		return 0 // Callback not found, surrender control.
	}

	keepControl, newPos := callback(inIsLosingControl != 0)

	if !keepControl {
		unregisterCallback(id) // Clean up if we are surrendering control.
		return 0
	}

	// If we are keeping control and a new position is provided, copy it to the C struct.
	if newPos != nil && outCameraPosition != nil {
		outCameraPosition.x = C.float(newPos.X)
		outCameraPosition.y = C.float(newPos.Y)
		outCameraPosition.z = C.float(newPos.Z)
		outCameraPosition.pitch = C.float(newPos.Pitch)
		outCameraPosition.heading = C.float(newPos.Heading)
		outCameraPosition.roll = C.float(newPos.Roll)
		outCameraPosition.zoom = C.float(newPos.Zoom)
	}

	return 1 // Keep control
}

// ControlCamera takes control of the X-Plane camera.
// You provide a duration and a callback function that will be executed every frame.
func ControlCamera(duration ControlDuration, callback CameraControlFunc) {
	id := registerCallback(callback)
	refcon := unsafe.Pointer(id)
	C.XPLMControlCamera(C.XPLMCameraControlDuration(duration), (C.XPLMCameraControl_f)(C.cameraControl_cgo), refcon)
}

// DontControlCamera surrenders control of the camera back to X-Plane.
func DontControlCamera() {
	C.XPLMDontControlCamera()
}

// IsCameraBeingControlled returns true if a plugin is controlling the camera,
// and if so, returns the duration of that control.
func IsCameraBeingControlled() (isControlled bool, duration ControlDuration) {
	var cDuration C.XPLMCameraControlDuration
	controlled := C.XPLMIsCameraBeingControlled(&cDuration)
	return controlled != 0, ControlDuration(cDuration)
}

// ReadCameraPosition reads the current position of the camera.
func ReadCameraPosition() Position {
	var cpos C.XPLMCameraPosition_t
	C.XPLMReadCameraPosition(&cpos)
	return Position{
		X:       float32(cpos.x),
		Y:       float32(cpos.y),
		Z:       float32(cpos.z),
		Pitch:   float32(cpos.pitch),
		Heading: float32(cpos.heading),
		Roll:    float32(cpos.roll),
		Zoom:    float32(cpos.zoom),
	}
}
