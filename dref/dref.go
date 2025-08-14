package dref

// #cgo CFLAGS: -DXPLM410=1
// #cgo LDFLAGS: -lXPLM_64
// #include <stdlib.h>
// #include "XPLMDataAccess.h"
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

type DataRef unsafe.Pointer

var (
	ErrDataRefNotFound  = errors.New("dataref not found")
	ErrRefNotRegistered = errors.New("dataref not registered in cache")
)

// DataRefCache holds pre-found DataRef handles for fast, repeated access.
type DataRefCache struct {
	refs  map[string]DataRef
	mutex sync.RWMutex
}

// NewDataRefCache creates a new, empty cache for datarefs.
func NewDataRefCache() *DataRefCache {
	return &DataRefCache{
		refs: make(map[string]DataRef),
	}
}

// Register finds a dataref by name and stores its handle in the cache.
// This should be done during plugin initialization.
func (c *DataRefCache) Register(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Avoid re-finding if already present
	if _, exists := c.refs[name]; exists {
		return nil
	}

	ref, err := FindDataRef(name)
	if err != nil {
		return fmt.Errorf("could not register dataref '%s': %w", name, err)
	}

	c.refs[name] = ref
	return nil
}

// getRef is an internal helper to safely get a registered dataref handle.
func (c *DataRefCache) getRef(name string) (DataRef, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ref, ok := c.refs[name]
	if !ok {
		return nil, fmt.Errorf("dataref '%s': %w", name, ErrRefNotRegistered)
	}
	return ref, nil
}

// GetInt retrieves the value of a pre-registered integer dataref.
func (c *DataRefCache) GetInt(name string) (int, error) {
	ref, err := c.getRef(name)
	if err != nil {
		return 0, err
	}
	return GetInt(ref), nil
}

// GetFloat retrieves the value of a pre-registered float dataref.
func (c *DataRefCache) GetFloat(name string) (float32, error) {
	ref, err := c.getRef(name)
	if err != nil {
		return 0.0, err
	}
	return GetFloat(ref), nil
}

// GetDouble retrieves the value of a pre-registered double dataref.
func (c *DataRefCache) GetDouble(name string) (float64, error) {
	ref, err := c.getRef(name)
	if err != nil {
		return 0.0, err
	}
	return GetDouble(ref), nil
}

// GetString retrieves the value of a pre-registered string/byte dataref.
func (c *DataRefCache) GetString(name string) (string, error) {
	ref, err := c.getRef(name)
	if err != nil {
		return "", err
	}
	return GetString(ref), nil
}

// GetBytes retrieves the value of a pre-registered byte array dataref.
func (c *DataRefCache) GetBytes(name string, buffer []byte) (int, error) {
	ref, err := c.getRef(name)
	if err != nil {
		return 0, err
	}
	return GetBytes(ref, buffer), nil
}

// FindDataRef looks up a dataref by its string identifier.
// Returns an error if the dataref cannot be found.
func FindDataRef(name string) (DataRef, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ref := C.XPLMFindDataRef(cName)
	if ref == nil {
		return nil, ErrDataRefNotFound
	}
	return DataRef(ref), nil
}

// GetFloat reads the value of a float dataref.
func GetFloat(ref DataRef) float32 {
	return float32(C.XPLMGetDataf(C.XPLMDataRef(ref)))
}

// SetFloat sets the value of a float dataref.
func SetFloat(ref DataRef, value float32) {
	C.XPLMSetDataf(C.XPLMDataRef(ref), C.float(value))
}

// GetInt reads the value of an integer dataref.
func GetInt(ref DataRef) int {
	return int(C.XPLMGetDatai(C.XPLMDataRef(ref)))
}

// SetInt sets the value of an integer dataref.
func SetInt(ref DataRef, value int) {
	C.XPLMSetDatai(C.XPLMDataRef(ref), C.int(value))
}

// GetDouble reads the value of a double-precision float dataref.
func GetDouble(ref DataRef) float64 {
	return float64(C.XPLMGetDatad(C.XPLMDataRef(ref)))
}

// GetString reads a string dataref. It reads up to 256 bytes and trims null characters.
func GetString(ref DataRef) string {
	buffer := make([]byte, 256)
	count := GetBytes(ref, buffer)
	return strings.TrimRight(string(buffer[:count]), "\x00")
}

// SetDouble sets the value of a double-precision float dataref.
func SetDouble(ref DataRef, value float64) {
	C.XPLMSetDatad(C.XPLMDataRef(ref), C.double(value))
}

// GetBytes reads a byte array dataref into the provided slice.
// It returns the number of bytes actually read.
func GetBytes(ref DataRef, buffer []byte) int {
	if len(buffer) == 0 {
		return 0
	}
	return int(C.XPLMGetDatab(
		C.XPLMDataRef(ref),
		unsafe.Pointer(&buffer[0]),
		0,
		C.int(len(buffer)),
	))
}

// SetBytes writes a byte slice to a dataref.
func SetBytes(ref DataRef, data []byte) {
	if len(data) == 0 {
		return
	}
	C.XPLMSetDatab(
		C.XPLMDataRef(ref),
		unsafe.Pointer(&data[0]),
		0,
		C.int(len(data)),
	)
}
