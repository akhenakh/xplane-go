# xplane-go

`xplane-go` is a set of Go packages that provide bindings to the X-Plane SDK, allowing you to write X-Plane plugins in the Go programming language. The goal is to offer a clean, safe, and idiomatic Go interface to the underlying C APIs.

## Project Status

The library is under active development but provides wrappers for many of the most common X-Plane SDK features. It is suitable for building functional plugins.

## Getting Started

### Building the Example Plugin

1. **SDK** Download the X-Plane SDK from the [X-Plane website](https://www.x-plane.com/download/sdk/).
   Copy the libraries into the `lib` directory at the root of the repository.
   Copy the header files.

1.  **Build:**  Using `go build`. The `-buildmode=c-shared` flag is essential to create a shared library that X-Plane can load.

    Example for Linux:
    ```bash
    CGO_CFLAGS="-DLIN=1 -I/usr/include/xplane_sdk/XPLM -I/usr/include/xplane_sdk/Widgets"  CGO_LDFLAGS="-L$(pwd)/lib"  go build -buildmode=c-shared -a -o hello.xpl  ./cmd/hello .
    ```

1.  **Install the Plugin:**
    *   Copy the output file (`hello.xpl`) to your X-Plane installation's plugin directory: `X-Plane 12/Resources/plugins/`.
    *   Create a folder for your plugin inside the `plugins` directory (e.g., `X-Plane 12/Resources/plugins/HelloGo/`).


1.  **Run X-Plane:** Start X-Plane. You should see messages from the plugin in the `Log.txt` file, and a new "HelloGo Plugin" menu will appear under the "Plugins" menu.

## Library Packages

The repository is organized into several packages, each wrapping a specific part of the X-Plane SDK.

### `plugin`

This is the most important package. It defines the core `plugin.Plugin` interface that every plugin must implement. It handles the C-to-Go bridge for the five required X-Plane entry points (`XPluginStart`, `XPluginStop`, `XPluginEnable`, `XPluginDisable`, `XPluginReceiveMessage`).

### `dref`

Provides access to X-Plane's data system (datarefs).

The key feature is the `dref.DataRefCache`, which allows you to find and cache dataref handles for high-performance access. This avoids the overhead of calling `XPLMFindDataRef` repeatedly in performance-critical code like flight loops.

**Recommended Usage:**
1.  Create a `NewDataRefCache` when your plugin is enabled.
2.  `Register()` all required datarefs by name.
3.  During runtime, use the `GetFloat()`, `GetInt()`, etc., methods to read values quickly from the cache.

```go
// In Enable():
p.datarefCache = dref.NewDataRefCache()
p.datarefCache.Register("sim/flightmodel/position/latitude")

// In a flight loop or menu handler:
lat, err := p.datarefCache.GetDouble("sim/flightmodel/position/latitude")
if err == nil {
    util.DebugString(fmt.Sprintf("Latitude is %f\n", lat))
}
```

### `processing`

Wraps the `XPLMProcessing` API. It allows you to register flight loop callbacks that are executed by X-Plane at a specified interval or phase (e.g., before or after the flight model). This is the primary mechanism for doing work on every frame or on a timer.

### `menu`

Wraps the `XPLMMenus` API for creating and managing plugin menus. You can create top-level menus, add items and separators, and handle user clicks.

### `camera`

Wraps the `XPLMCamera` API, allowing you to take programmatic control of the X-Plane camera. You can use `ControlCamera()` with a callback function that is executed every frame to set the camera's position, orientation, and zoom.

### `util`

Contains wrappers for simple utility functions, most notably `util.DebugString`, which is the standard way to write messages to X-Plane's `Log.txt` file for debugging.

### `display` & `widget`

Provide wrappers and constants for the `XPLMDisplay` and `XPWidgets` APIs, which are used for creating 2D user interfaces, windows, and standard UI controls like buttons and text fields.
