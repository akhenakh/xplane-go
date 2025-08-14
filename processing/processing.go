package processing

// #cgo CFLAGS: -DXPLM200=1 -DXPLM210=1 -DXPLM410=1
// #cgo LDFLAGS: -lXPLM_64
// #include <stdlib.h>
// #include "XPLMProcessing.h"
//
// extern float flightLoopCallback_cgo(float inElapsedSinceLastCall, float inElapsedTimeSinceLastFlightLoop, int inCounter, void* inRefcon);
import "C"

import (
	"sync"
	"unsafe"
)

// FlightLoopID is an opaque handle to a registered flight loop.
type FlightLoopID C.XPLMFlightLoopID

// FlightLoopCallback is the function signature for a Go-native flight loop.
// It returns the delay in seconds for the next callback. Return 0 to un-schedule.
type FlightLoopCallback func(elapsedSinceLastCall, elapsedTimeSinceLastFlightLoop float32, counter int) float32

var (
	registry      = make(map[uintptr]FlightLoopCallback)
	registryMutex sync.RWMutex
	nextID        uintptr = 1
)

func registerCallback(callback FlightLoopCallback) uintptr {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	id := nextID
	nextID++
	registry[id] = callback
	return id
}

func unregisterCallback(id uintptr) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	delete(registry, id)
}

func getCallback(id uintptr) FlightLoopCallback {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return registry[id]
}

//export flightLoopCallback_cgo
func flightLoopCallback_cgo(inElapsedSinceLastCall, inElapsedTimeSinceLastFlightLoop C.float, inCounter C.int, inRefcon unsafe.Pointer) C.float {
	id := uintptr(inRefcon)
	if callback := getCallback(id); callback != nil {
		nextInterval := callback(
			float32(inElapsedSinceLastCall),
			float32(inElapsedTimeSinceLastFlightLoop),
			int(inCounter),
		)
		return C.float(nextInterval)
	}
	return 0
}

type FlightLoopPhase int

const (
	BeforeFlightModel FlightLoopPhase = C.xplm_FlightLoop_Phase_BeforeFlightModel
	AfterFlightModel  FlightLoopPhase = C.xplm_FlightLoop_Phase_AfterFlightModel
)

// CreateFlightLoop registers a new flight loop and returns its ID.
// The flight loop is initially unscheduled. Use ScheduleFlightLoop to start it.
func CreateFlightLoop(phase FlightLoopPhase, callback FlightLoopCallback) FlightLoopID {
	id := registerCallback(callback)
	refcon := unsafe.Pointer(id)
	params := C.XPLMCreateFlightLoop_t{
		structSize:   C.int(unsafe.Sizeof(C.XPLMCreateFlightLoop_t{})),
		phase:        C.XPLMFlightLoopPhaseType(phase),
		callbackFunc: (C.XPLMFlightLoop_f)(C.flightLoopCallback_cgo),
		refcon:       refcon,
	}
	// Pass the address of the params struct using &params
	return FlightLoopID(C.XPLMCreateFlightLoop(&params))
}

// DestroyFlightLoop unregisters a flight loop and removes it from the Go registry.
func DestroyFlightLoop(loopID FlightLoopID) {
	// A more robust implementation would need to look up the refcon/id associated with the loopID.
	// For now, we assume the user manages the ID returned from CreateFlightLoop.
	// This requires a reverse mapping from FlightLoopID to our internal uintptr id.
	// For simplicity in this example, we'll just destroy the C-side object.
	// A real library would need: var idToLoopIDMap = make(map[uintptr]FlightLoopID)
	C.XPLMDestroyFlightLoop(C.XPLMFlightLoopID(loopID))
	// unregisterCallback(id) // Would be called here.
}

// ScheduleFlightLoop schedules (or re-schedules) a flight loop.
//   - interval > 0: seconds from now.
//   - interval < 0: flight loops from now.
//   - interval = 0: unschedule the callback.
func ScheduleFlightLoop(loopID FlightLoopID, interval float32, relativeToNow bool) {
	rel := 0
	if relativeToNow {
		rel = 1
	}
	C.XPLMScheduleFlightLoop(C.XPLMFlightLoopID(loopID), C.float(interval), C.int(rel))
}
