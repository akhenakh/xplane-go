package main

import (
	"fmt"
	"log"

	"github.com/akhenakh/xplane-go/camera"
	"github.com/akhenakh/xplane-go/dref"
	"github.com/akhenakh/xplane-go/menu"
	"github.com/akhenakh/xplane-go/plugin"
	"github.com/akhenakh/xplane-go/processing"
	"github.com/akhenakh/xplane-go/util"
)

// The main entry point for the plugin is to have an init function
// that registers an implementation of the plugin.Plugin interface.
func init() {
	plugin.Register(&HelloPlugin{})
}

// HelloPlugin holds the state for our plugin.
type HelloPlugin struct {
	ourMenu menu.MenuID

	// A generic cache for all datarefs used by the plugin.
	datarefCache *dref.DataRefCache

	// Flight Loop ID for cleanup
	ourFlightLoop processing.FlightLoopID

	// Camera Demo State
	isCameraShaking bool
	shakeCounter    float32
}

// Start is called by X-Plane when the plugin is first loaded.
func (p *HelloPlugin) Start() (name, sig, desc string, err error) {
	util.DebugString("HelloPlugin: Started!\n")
	return "Hello World Go Plugin", "xplane-go.example.helloworld", "A plugin in Go demonstrating menus and datarefs.", nil
}

// Stop is called by X-Plane when the plugin is unloaded.
func (p *HelloPlugin) Stop() {
	p.cleanup()
	util.DebugString("HelloPlugin: Stopped!\n")
}

// Enable is called by X-Plane when the plugin is enabled.
func (p *HelloPlugin) Enable() error {
	util.DebugString("HelloPlugin: Enabled!\n")

	pluginsMenu := menu.FindPluginsMenu()
	if pluginsMenu == nil {
		return fmt.Errorf("could not find plugins menu")
	}

	containerIndex := menu.AppendMenuItem(pluginsMenu, "HelloGo Plugin", nil)
	if containerIndex < 0 {
		return fmt.Errorf("could not create menu item in plugins menu")
	}

	p.ourMenu = menu.CreateMenu("HelloGo", pluginsMenu, containerIndex, p.helloMenuHandler, p)
	if p.ourMenu == nil {
		return fmt.Errorf("could not create the menu container")
	}

	menu.AppendMenuItem(p.ourMenu, "Say Hello", "hello_item")
	menu.AppendMenuSeparator(p.ourMenu)
	menu.AppendMenuItem(p.ourMenu, "Toggle Camera Shake", "camera_item")
	menu.AppendMenuSeparator(p.ourMenu)
	menu.AppendMenuItem(p.ourMenu, "Log Current Position", "log_pos_item") // Changed from ADS-B

	// DataRef Cache Initialization
	p.datarefCache = dref.NewDataRefCache()
	datarefsToRegister := []string{
		// For simple demo
		"sim/time/zulu_time_sec",
		"sim/graphics/view/field_of_view_deg",
		// For camera demo
		"sim/flightmodel/position/local_x",
		"sim/flightmodel/position/local_y",
		"sim/flightmodel/position/local_z",
		// For position demo
		"sim/flightmodel/position/latitude",
		"sim/flightmodel/position/longitude",
		"sim/cockpit2/gauges/indicators/altitude_ft_pilot",
	}

	for _, name := range datarefsToRegister {
		if err := p.datarefCache.Register(name); err != nil {
			log.Printf("FATAL: Failed to register dataref '%s': %v", name, err)
			return err // Fail to enable if a dataref is missing
		}
	}
	log.Println("Successfully initialized and registered all datarefs in cache.")
	util.DebugString("Successfully initialized and registered all datarefs in cache.\n")

	//  Flight Loop Setup
	p.ourFlightLoop = processing.CreateFlightLoop(processing.AfterFlightModel, p.flightLoopCallback)
	processing.ScheduleFlightLoop(p.ourFlightLoop, 2.0, true)

	return nil
}

// Disable is called by X-Plane when the plugin is disabled.
func (p *HelloPlugin) Disable() {
	p.cleanup()
	util.DebugString("HelloPlugin: Disabled!\n")
}

// flightLoopCallback is called by X-Plane at a regular interval.
func (p *HelloPlugin) flightLoopCallback(elapsedSinceLastCall, elapsedTimeSinceLastFlightLoop float32, counter int) float32 {
	zuluTime, err := p.datarefCache.GetFloat("sim/time/zulu_time_sec")
	if err != nil {
		util.DebugString(fmt.Sprintf("Go Plugin Error: %v\n", err))
		return 2.0 // Continue trying
	}
	util.DebugString(fmt.Sprintf("Go Plugin: Current Zulu Time is %f\n", zuluTime))

	currentFOV, err := p.datarefCache.GetFloat("sim/graphics/view/field_of_view_deg")
	if err != nil {
		util.DebugString(fmt.Sprintf("Go Plugin Error: %v\n", err))
		return 2.0 // Continue trying
	}
	var newFOV float32 = 60.0
	if currentFOV < 70.0 {
		newFOV = 80.0
	}

	// For setting a dataref, we still need to find it directly.
	// The cache is primarily for optimized reading.
	fovRef, err := dref.FindDataRef("sim/graphics/view/field_of_view_deg")
	if err == nil {
		dref.SetFloat(fovRef, newFOV)
		util.DebugString(fmt.Sprintf("Go Plugin: Set FOV from %f to %f\n", currentFOV, newFOV))
	}

	return 2.0
}

// cameraCallback controls the camera on a per-frame basis.
func (p *HelloPlugin) cameraCallback(isLosingControl bool) (bool, *camera.Position) {
	if isLosingControl || !p.isCameraShaking {
		p.isCameraShaking = false
		return false, nil // Surrender control
	}

	posX, errX := p.datarefCache.GetFloat("sim/flightmodel/position/local_x")
	posY, errY := p.datarefCache.GetFloat("sim/flightmodel/position/local_y")
	posZ, errZ := p.datarefCache.GetFloat("sim/flightmodel/position/local_z")
	if errX != nil || errY != nil || errZ != nil {
		log.Printf("Error reading camera position datarefs: %v, %v, %v", errX, errY, errZ)
		return true, nil // Keep control but don't move camera
	}

	p.shakeCounter += 0.2
	offset := float32(0.5 * (p.shakeCounter - float32(int(p.shakeCounter/2)*2)))

	newPos := camera.Position{
		X:       posX,
		Y:       posY + 5.0 + offset,
		Z:       posZ - 15.0,
		Pitch:   -5.0,
		Heading: 0,
		Roll:    0,
		Zoom:    1.0,
	}

	return true, &newPos
}

// helloMenuHandler handles all clicks in our menu.
func (p *HelloPlugin) helloMenuHandler(menuRef, itemRef interface{}) {
	pluginInstance, ok := menuRef.(*HelloPlugin)
	if !ok {
		util.DebugString(fmt.Sprintf("Error: menu handler called with invalid menuRef type: %T\n", menuRef))
		return
	}
	itemID, ok := itemRef.(string)
	if !ok {
		util.DebugString(fmt.Sprintf("Error: menu handler called with invalid itemRef type: %T\n", itemRef))
		return
	}

	switch itemID {
	case "hello_item":
		log.Println("Hello Go Menu Item Clicked!")
		util.DebugString("The 'Say Hello' menu item was clicked!\n")
	case "camera_item":
		pluginInstance.isCameraShaking = !pluginInstance.isCameraShaking
		if pluginInstance.isCameraShaking {
			log.Println("Starting camera shake.")
			camera.ControlCamera(camera.Forever, pluginInstance.cameraCallback)
		} else {
			log.Println("Stopping camera shake.")
			camera.DontControlCamera()
		}
	case "log_pos_item":
		// This case now calls our new, simpler function.
		pluginInstance.logCurrentPosition()
	}
}

// logCurrentPosition reads lat, lon, and alt from the cache and prints them.
func (p *HelloPlugin) logCurrentPosition() {
	if p.datarefCache == nil {
		log.Println("Error: Dataref cache is not initialized.")
		util.DebugString("Error: Dataref cache is not initialized.\n")
		return
	}

	// Read the required datarefs from our cache
	lat, errLat := p.datarefCache.GetDouble("sim/flightmodel/position/latitude")
	lon, errLon := p.datarefCache.GetDouble("sim/flightmodel/position/longitude")
	alt, errAlt := p.datarefCache.GetFloat("sim/cockpit2/gauges/indicators/altitude_ft_pilot")

	// Check if any of the reads failed
	if errLat != nil || errLon != nil || errAlt != nil {
		errMsg := fmt.Sprintf("Error reading position datarefs: latErr=%v, lonErr=%v, altErr=%v", errLat, errLon, errAlt)
		log.Println(errMsg)
		util.DebugString(errMsg + "\n")
		return
	}

	// Format and log the successfully retrieved data
	logMsg := fmt.Sprintf("Current Position -> Lat: %.4f, Lon: %.4f, Alt: %.0f ft", lat, lon, alt)
	log.Println(logMsg)
	util.DebugString(logMsg + "\n")
}

// cleanup is a helper function to avoid duplicating code in Stop() and Disable().
func (p *HelloPlugin) cleanup() {
	if p.ourMenu != nil {
		menu.DestroyMenu(p.ourMenu)
		p.ourMenu = nil
	}
	if p.ourFlightLoop != nil {
		processing.DestroyFlightLoop(p.ourFlightLoop)
		p.ourFlightLoop = nil
	}
	if p.isCameraShaking {
		camera.DontControlCamera()
		p.isCameraShaking = false
	}
	p.datarefCache = nil // Clear the cache reference
}

// This empty main function is required by c-shared builds, but is not executed.
func main() {}
