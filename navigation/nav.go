package navigation

// #cgo CFLAGS: -DXPLM410=1
// #include <stdlib.h>
// #include "XPLMNavigation.h"
import "C"
import (
	"errors"
	"unsafe"
)

// NavType represents the type of navigation aid.
type NavType int

const (
	NavUnknown      NavType = C.xplm_Nav_Unknown
	NavAirport      NavType = C.xplm_Nav_Airport
	NavNDB          NavType = C.xplm_Nav_NDB
	NavVOR          NavType = C.xplm_Nav_VOR
	NavILS          NavType = C.xplm_Nav_ILS
	NavLocalizer    NavType = C.xplm_Nav_Localizer
	NavGlideSlope   NavType = C.xplm_Nav_GlideSlope
	NavOuterMarker  NavType = C.xplm_Nav_OuterMarker
	NavMiddleMarker NavType = C.xplm_Nav_MiddleMarker
	NavInnerMarker  NavType = C.xplm_Nav_InnerMarker
	NavFix          NavType = C.xplm_Nav_Fix
	NavDME          NavType = C.xplm_Nav_DME
	NavLatLon       NavType = C.xplm_Nav_LatLon
	NavTACAN        NavType = C.xplm_Nav_TACAN
)

// NavRef is a reference to a navigation aid in the database.
type NavRef int

const (
	NavNotFound NavRef = C.XPLM_NAV_NOT_FOUND
)

// NavFlightPlan represents different flight plans available in the system.
type NavFlightPlan int

const (
	FplPilotPrimary     NavFlightPlan = C.xplm_Fpl_Pilot_Primary
	FplCoPilotPrimary   NavFlightPlan = C.xplm_Fpl_CoPilot_Primary
	FplPilotApproach    NavFlightPlan = C.xplm_Fpl_Pilot_Approach
	FplCoPilotApproach  NavFlightPlan = C.xplm_Fpl_CoPilot_Approach
	FplPilotTemporary   NavFlightPlan = C.xplm_Fpl_Pilot_Temporary
	FplCoPilotTemporary NavFlightPlan = C.xplm_Fpl_CoPilot_Temporary
)

// NavAidInfo contains information about a navigation aid.
type NavAidInfo struct {
	Type      NavType
	Latitude  float64
	Longitude float64
	Height    float64
	Frequency int
	Heading   float64
	ID        string
	Name      string
	Region    string
}

// FMSFlightPlanEntryInfo contains information about an FMS flight plan entry.
type FMSFlightPlanEntryInfo struct {
	Type      NavType
	ID        string
	Ref       NavRef
	Altitude  int
	Latitude  float64
	Longitude float64
}

// GetFirstNavAid returns the first navigation aid in the database.
func GetFirstNavAid() NavRef {
	return NavRef(C.XPLMGetFirstNavAid())
}

// GetNextNavAid returns the next navigation aid after the given reference.
func GetNextNavAid(ref NavRef) NavRef {
	return NavRef(C.XPLMGetNextNavAid(C.XPLMNavRef(ref)))
}

// FindFirstNavAidOfType finds the first navigation aid of the specified type.
func FindFirstNavAidOfType(navType NavType) NavRef {
	return NavRef(C.XPLMFindFirstNavAidOfType(C.XPLMNavType(navType)))
}

// FindLastNavAidOfType finds the last navigation aid of the specified type.
func FindLastNavAidOfType(navType NavType) NavRef {
	return NavRef(C.XPLMFindLastNavAidOfType(C.XPLMNavType(navType)))
}

// FindNavAid searches for navigation aids matching the specified criteria.
func FindNavAid(nameFragment, idFragment string, lat, lon *float64, frequency *int, navType NavType) NavRef {
	var cNameFragment *C.char
	var cIDFragment *C.char

	if nameFragment != "" {
		cNameFragment = C.CString(nameFragment)
		defer C.free(unsafe.Pointer(cNameFragment))
	}

	if idFragment != "" {
		cIDFragment = C.CString(idFragment)
		defer C.free(unsafe.Pointer(cIDFragment))
	}

	var cLat *C.float
	var cLon *C.float
	var cFrequency *C.int

	if lat != nil {
		cLat = (*C.float)(unsafe.Pointer(lat))
	}

	if lon != nil {
		cLon = (*C.float)(unsafe.Pointer(lon))
	}

	if frequency != nil {
		cFrequency = (*C.int)(unsafe.Pointer(frequency))
	}

	return NavRef(C.XPLMFindNavAid(
		cNameFragment,
		cIDFragment,
		cLat,
		cLon,
		cFrequency,
		C.XPLMNavType(navType),
	))
}

// GetNavAidInfo retrieves information about a navigation aid.
func GetNavAidInfo(ref NavRef) (NavAidInfo, error) {
	var outType C.XPLMNavType
	var outLatitude C.float
	var outLongitude C.float
	var outHeight C.float
	var outFrequency C.int
	var outHeading C.float
	var outID [64]C.char
	var outName [256]C.char
	var outReg [1]C.char

	C.XPLMGetNavAidInfo(
		C.XPLMNavRef(ref),
		(*C.XPLMNavType)(&outType),
		(*C.float)(&outLatitude),
		(*C.float)(&outLongitude),
		(*C.float)(&outHeight),
		(*C.int)(&outFrequency),
		(*C.float)(&outHeading),
		(*C.char)(unsafe.Pointer(&outID)),
		(*C.char)(unsafe.Pointer(&outName)),
		(*C.char)(unsafe.Pointer(&outReg)),
	)

	info := NavAidInfo{
		Type:      NavType(outType),
		Latitude:  float64(outLatitude),
		Longitude: float64(outLongitude),
		Height:    float64(outHeight),
		Frequency: int(outFrequency),
		Heading:   float64(outHeading),
		ID:        C.GoString((*C.char)(unsafe.Pointer(&outID))),
		Name:      C.GoString((*C.char)(unsafe.Pointer(&outName))),
		Region:    C.GoString((*C.char)(unsafe.Pointer(&outReg))),
	}

	return info, nil
}

// CountFMSEntries returns the number of entries in the FMS.
func CountFMSEntries() int {
	return int(C.XPLMCountFMSEntries())
}

// GetDisplayedFMSEntry returns the index of the currently displayed FMS entry.
func GetDisplayedFMSEntry() int {
	return int(C.XPLMGetDisplayedFMSEntry())
}

// GetDestinationFMSEntry returns the index of the destination FMS entry.
func GetDestinationFMSEntry() int {
	return int(C.XPLMGetDestinationFMSEntry())
}

// SetDisplayedFMSEntry sets the currently displayed FMS entry.
func SetDisplayedFMSEntry(index int) {
	C.XPLMSetDisplayedFMSEntry(C.int(index))
}

// SetDestinationFMSEntry sets the destination FMS entry.
func SetDestinationFMSEntry(index int) {
	C.XPLMSetDestinationFMSEntry(C.int(index))
}

// GetFMSEntryInfo retrieves information about a specific FMS entry.
func GetFMSEntryInfo(index int) (FMSFlightPlanEntryInfo, error) {
	var outType C.XPLMNavType
	var outID [256]C.char
	var outRef C.XPLMNavRef
	var outAltitude C.int
	var outLat C.float
	var outLon C.float

	// Initialize outRef to NavNotFound to handle the bug in X-Plane prior to 11.31
	outRef = C.XPLMNavRef(C.XPLM_NAV_NOT_FOUND)

	C.XPLMGetFMSEntryInfo(
		C.int(index),
		(*C.XPLMNavType)(&outType),
		(*C.char)(unsafe.Pointer(&outID)),
		(*C.XPLMNavRef)(&outRef),
		(*C.int)(&outAltitude),
		(*C.float)(&outLat),
		(*C.float)(&outLon),
	)

	info := FMSFlightPlanEntryInfo{
		Type:      NavType(outType),
		ID:        C.GoString((*C.char)(unsafe.Pointer(&outID))),
		Ref:       NavRef(outRef),
		Altitude:  int(outAltitude),
		Latitude:  float64(outLat),
		Longitude: float64(outLon),
	}

	return info, nil
}

// SetFMSEntryInfo sets the destination navaid and altitude for an FMS entry.
func SetFMSEntryInfo(index int, ref NavRef, altitude int) {
	C.XPLMSetFMSEntryInfo(C.int(index), C.XPLMNavRef(ref), C.int(altitude))
}

// SetFMSEntryLatLon sets a lat/lon entry in the FMS.
func SetFMSEntryLatLon(index int, lat, lon float64, altitude int) {
	C.XPLMSetFMSEntryLatLon(C.int(index), C.float(lat), C.float(lon), C.int(altitude))
}

// ClearFMSEntry clears a specific FMS entry.
func ClearFMSEntry(index int) {
	C.XPLMClearFMSEntry(C.int(index))
}

// CountFMSFlightPlanEntries returns the number of entries in the specified flight plan.
func CountFMSFlightPlanEntries(flightPlan NavFlightPlan) int {
	return int(C.XPLMCountFMSFlightPlanEntries(C.XPLMNavFlightPlan(flightPlan)))
}

// GetDisplayedFMSFlightPlanEntry returns the index of the displayed entry in the specified flight plan.
func GetDisplayedFMSFlightPlanEntry(flightPlan NavFlightPlan) int {
	return int(C.XPLMGetDisplayedFMSFlightPlanEntry(C.XPLMNavFlightPlan(flightPlan)))
}

// GetDestinationFMSFlightPlanEntry returns the index of the destination entry in the specified flight plan.
func GetDestinationFMSFlightPlanEntry(flightPlan NavFlightPlan) int {
	return int(C.XPLMGetDestinationFMSFlightPlanEntry(C.XPLMNavFlightPlan(flightPlan)))
}

// SetDisplayedFMSFlightPlanEntry sets the displayed entry in the specified flight plan.
func SetDisplayedFMSFlightPlanEntry(flightPlan NavFlightPlan, index int) {
	C.XPLMSetDisplayedFMSFlightPlanEntry(C.XPLMNavFlightPlan(flightPlan), C.int(index))
}

// SetDestinationFMSFlightPlanEntry sets the destination entry in the specified flight plan.
func SetDestinationFMSFlightPlanEntry(flightPlan NavFlightPlan, index int) {
	C.XPLMSetDestinationFMSFlightPlanEntry(C.XPLMNavFlightPlan(flightPlan), C.int(index))
}

// SetDirectToFMSFlightPlanEntry sets the direct-to entry in the specified flight plan.
func SetDirectToFMSFlightPlanEntry(flightPlan NavFlightPlan, index int) {
	C.XPLMSetDirectToFMSFlightPlanEntry(C.XPLMNavFlightPlan(flightPlan), C.int(index))
}

// GetFMSFlightPlanEntryInfo retrieves information about a specific entry in the specified flight plan.
func GetFMSFlightPlanEntryInfo(flightPlan NavFlightPlan, index int) (FMSFlightPlanEntryInfo, error) {
	var outType C.XPLMNavType
	var outID [256]C.char
	var outRef C.XPLMNavRef
	var outAltitude C.int
	var outLat C.float
	var outLon C.float

	// Initialize outRef to NavNotFound to handle the bug in X-Plane prior to 11.31
	outRef = C.XPLMNavRef(C.XPLM_NAV_NOT_FOUND)

	C.XPLMGetFMSFlightPlanEntryInfo(
		C.XPLMNavFlightPlan(flightPlan),
		C.int(index),
		(*C.XPLMNavType)(&outType),
		(*C.char)(unsafe.Pointer(&outID)),
		(*C.XPLMNavRef)(&outRef),
		(*C.int)(&outAltitude),
		(*C.float)(&outLat),
		(*C.float)(&outLon),
	)

	info := FMSFlightPlanEntryInfo{
		Type:      NavType(outType),
		ID:        C.GoString((*C.char)(unsafe.Pointer(&outID))),
		Ref:       NavRef(outRef),
		Altitude:  int(outAltitude),
		Latitude:  float64(outLat),
		Longitude: float64(outLon),
	}

	return info, nil
}

// SetFMSFlightPlanEntryInfo sets the destination navaid and altitude for an entry in the specified flight plan.
func SetFMSFlightPlanEntryInfo(flightPlan NavFlightPlan, index int, ref NavRef, altitude int) {
	C.XPLMSetFMSFlightPlanEntryInfo(C.XPLMNavFlightPlan(flightPlan), C.int(index), C.XPLMNavRef(ref), C.int(altitude))
}

// SetFMSFlightPlanEntryLatLon sets a lat/lon entry in the specified flight plan.
func SetFMSFlightPlanEntryLatLon(flightPlan NavFlightPlan, index int, lat, lon float64, altitude int) {
	C.XPLMSetFMSFlightPlanEntryLatLon(C.XPLMNavFlightPlan(flightPlan), C.int(index), C.float(lat), C.float(lon), C.int(altitude))
}

// SetFMSFlightPlanEntryLatLonWithId sets a lat/lon entry with an ID in the specified flight plan.
func SetFMSFlightPlanEntryLatLonWithId(flightPlan NavFlightPlan, index int, lat, lon float64, altitude int, id string) {
	cID := C.CString(id)
	defer C.free(unsafe.Pointer(cID))
	C.XPLMSetFMSFlightPlanEntryLatLonWithId(
		C.XPLMNavFlightPlan(flightPlan),
		C.int(index),
		C.float(lat),
		C.float(lon),
		C.int(altitude),
		cID,
		C.uint(len(id)),
	)
}

// ClearFMSFlightPlanEntry clears a specific entry in the specified flight plan.
func ClearFMSFlightPlanEntry(flightPlan NavFlightPlan, index int) {
	C.XPLMClearFMSFlightPlanEntry(C.XPLMNavFlightPlan(flightPlan), C.int(index))
}

// LoadFMSFlightPlan loads a flight plan from a buffer into the specified device.
func LoadFMSFlightPlan(device int, buffer string) {
	cBuffer := C.CString(buffer)
	defer C.free(unsafe.Pointer(cBuffer))
	C.XPLMLoadFMSFlightPlan(C.int(device), cBuffer, C.uint(len(buffer)))
}

// SaveFMSFlightPlan saves a flight plan from the specified device to a buffer.
func SaveFMSFlightPlan(device int, buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, errors.New("buffer cannot be empty")
	}

	// SaveFMSFlightPlan returns the required buffer size, so we need to call it first
	// with a small buffer to get the required size
	requiredSize := int(C.XPLMSaveFMSFlightPlan(C.int(device), (*C.char)(unsafe.Pointer(&buffer[0])), C.uint(len(buffer))))
	if requiredSize <= len(buffer) {
		// The buffer was large enough, return the actual size written
		return requiredSize, nil
	}

	// Buffer was too small, create a new buffer of the required size
	newBuffer := make([]byte, requiredSize)
	actualSize := int(C.XPLMSaveFMSFlightPlan(C.int(device), (*C.char)(unsafe.Pointer(&newBuffer[0])), C.uint(len(newBuffer))))
	return actualSize, nil
}

// GetGPSDestinationType returns the type of the currently selected GPS destination.
func GetGPSDestinationType() NavType {
	return NavType(C.XPLMGetGPSDestinationType())
}

// GetGPSDestination returns the currently selected GPS destination.
func GetGPSDestination() NavRef {
	return NavRef(C.XPLMGetGPSDestination())
}
