package menu

// #cgo CFLAGS: -DXPLM410=1
// #include <stdlib.h>
// #include "XPLMMenus.h"
//
// extern void menuHandler_cgo(void* inMenuRef, void* inItemRef);
import "C"

import (
	"sync"
	"unsafe"
)

type MenuID C.XPLMMenuID

type Handler func(menuRef, itemRef interface{})

var (
	handlerRegistry      = make(map[uintptr]Handler)
	handlerRegistryMutex sync.RWMutex
	nextHandlerID        uintptr = 1

	itemRegistry      = make(map[uintptr]interface{})
	itemRegistryMutex sync.RWMutex
	nextItemID        uintptr = 1

	menuIDToHandlerIDMap   = make(map[MenuID]uintptr)
	menuIDToHandlerIDMutex sync.RWMutex

	menuIDToItemIDsMap   = make(map[MenuID][]uintptr)
	menuIDToItemIDsMutex sync.RWMutex

	// This map stores the user's original menuRef, associated with our internal handler ID.
	handlerIDToMenuRefMap   = make(map[uintptr]interface{})
	handlerIDToMenuRefMutex sync.RWMutex
)

//  Registry Management

func registerHandler(handler Handler, menuRef interface{}) uintptr {
	handlerRegistryMutex.Lock()
	defer handlerRegistryMutex.Unlock()

	id := nextHandlerID
	nextHandlerID++

	handlerRegistry[id] = handler

	// NEW: Store the user's refcon associated with this handler ID.
	handlerIDToMenuRefMutex.Lock()
	handlerIDToMenuRefMap[id] = menuRef
	handlerIDToMenuRefMutex.Unlock()

	return id
}

func unregisterHandler(id uintptr) {
	if id == 0 {
		return
	}
	handlerRegistryMutex.Lock()
	delete(handlerRegistry, id)
	handlerRegistryMutex.Unlock()

	//  Clean up the refcon map
	handlerIDToMenuRefMutex.Lock()
	delete(handlerIDToMenuRefMap, id)
	handlerIDToMenuRefMutex.Unlock()
}

func getHandler(id uintptr) Handler {
	handlerRegistryMutex.RLock()
	defer handlerRegistryMutex.RUnlock()
	return handlerRegistry[id]
}

// Function to get the stored user menuRef.
func getMenuRef(id uintptr) interface{} {
	handlerIDToMenuRefMutex.RLock()
	defer handlerIDToMenuRefMutex.RUnlock()
	return handlerIDToMenuRefMap[id]
}

func registerItemRef(ref interface{}) uintptr {
	itemRegistryMutex.Lock()
	defer itemRegistryMutex.Unlock()
	id := nextItemID
	nextItemID++
	itemRegistry[id] = ref
	return id
}

func unregisterItemRef(id uintptr) {
	if id == 0 {
		return
	}
	itemRegistryMutex.Lock()
	defer itemRegistryMutex.Unlock()
	delete(itemRegistry, id)
}

func getItemRef(id uintptr) interface{} {
	if id == 0 {
		return nil
	}
	itemRegistryMutex.RLock()
	defer itemRegistryMutex.RUnlock()
	return itemRegistry[id]
}

// CGO Trampoline

//export menuHandler_cgo
func menuHandler_cgo(inMenuRef, inItemRef unsafe.Pointer) {
	handlerID := uintptr(inMenuRef)
	handler := getHandler(handlerID)
	if handler == nil {
		return
	}

	// Look up the user's original menuRef.
	actualMenuRef := getMenuRef(handlerID)

	itemID := uintptr(inItemRef)
	actualItemRef := getItemRef(itemID)

	// Pass the correct, original menuRef to the user's handler.
	handler(actualMenuRef, actualItemRef)
}

func FindPluginsMenu() MenuID {
	return MenuID(C.XPLMFindPluginsMenu())
}

func CreateMenu(name string, parent MenuID, parentItem int, handler Handler, menuRef interface{}) MenuID {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cMenuRef unsafe.Pointer
	var cCallback C.XPLMMenuHandler_f
	var handlerID uintptr

	if handler != nil {
		// NEW: Pass the user's menuRef to the registration function.
		handlerID = registerHandler(handler, menuRef)
		cMenuRef = unsafe.Pointer(handlerID)
		cCallback = (C.XPLMMenuHandler_f)(C.menuHandler_cgo)
	}

	newID := MenuID(C.XPLMCreateMenu(cName, C.XPLMMenuID(parent), C.int(parentItem), cCallback, cMenuRef))

	if newID != nil {
		if handler != nil {
			menuIDToHandlerIDMutex.Lock()
			menuIDToHandlerIDMap[newID] = handlerID
			menuIDToHandlerIDMutex.Unlock()
		}
		menuIDToItemIDsMutex.Lock()
		menuIDToItemIDsMap[newID] = make([]uintptr, 0)
		menuIDToItemIDsMutex.Unlock()
	}

	return newID
}

func AppendMenuItem(menuID MenuID, name string, itemRef interface{}) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var cItemRef unsafe.Pointer
	var itemID uintptr
	if itemRef != nil {
		itemID = registerItemRef(itemRef)
		cItemRef = unsafe.Pointer(itemID)
	}
	itemIndex := int(C.XPLMAppendMenuItem(C.XPLMMenuID(menuID), cName, cItemRef, 1))
	if itemIndex >= 0 && itemRef != nil {
		menuIDToItemIDsMutex.Lock()
		menuIDToItemIDsMap[menuID] = append(menuIDToItemIDsMap[menuID], itemID)
		menuIDToItemIDsMutex.Unlock()
	}
	return itemIndex
}

func AppendMenuSeparator(menuID MenuID) {
	C.XPLMAppendMenuSeparator(C.XPLMMenuID(menuID))
}

func ClearAllMenuItems(menuID MenuID) {
	if menuID == nil {
		return
	}
	menuIDToItemIDsMutex.Lock()
	itemIDs, ok := menuIDToItemIDsMap[menuID]
	if ok {
		menuIDToItemIDsMap[menuID] = make([]uintptr, 0)
		menuIDToItemIDsMutex.Unlock()
		for _, id := range itemIDs {
			unregisterItemRef(id)
		}
	} else {
		menuIDToItemIDsMutex.Unlock()
	}
	C.XPLMClearAllMenuItems(C.XPLMMenuID(menuID))
}

func DestroyMenu(menuID MenuID) {
	if menuID == nil {
		return
	}
	ClearAllMenuItems(menuID)
	menuIDToItemIDsMutex.Lock()
	delete(menuIDToItemIDsMap, menuID)
	menuIDToItemIDsMutex.Unlock()
	menuIDToHandlerIDMutex.Lock()
	handlerID, ok := menuIDToHandlerIDMap[menuID]
	if ok {
		delete(menuIDToHandlerIDMap, menuID)
		menuIDToHandlerIDMutex.Unlock()
		unregisterHandler(handlerID)
	} else {
		menuIDToHandlerIDMutex.Unlock()
	}
	C.XPLMDestroyMenu(C.XPLMMenuID(menuID))
}

func SetMenuItemName(menuID MenuID, index int, name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.XPLMSetMenuItemName(C.XPLMMenuID(menuID), C.int(index), cName, 0)
}
