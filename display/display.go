package display

// #cgo CFLAGS: -DXPLM410=1
// #cgo LDFLAGS: -lXPLM_64
// #include "XPLMDisplay.h"
import "C"

type WindowID C.XPLMWindowID

// KeyFlags represents modifier keys.
type KeyFlags int

const (
	ShiftFlag     KeyFlags = C.xplm_ShiftFlag
	OptionAltFlag KeyFlags = C.xplm_OptionAltFlag
	ControlFlag   KeyFlags = C.xplm_ControlFlag
	DownFlag      KeyFlags = C.xplm_DownFlag
	UpFlag        KeyFlags = C.xplm_UpFlag
)

// MouseStatus indicates the state of a mouse click.
type MouseStatus int

const (
	MouseDown MouseStatus = C.xplm_MouseDown
	MouseDrag MouseStatus = C.xplm_MouseDrag
	MouseUp   MouseStatus = C.xplm_MouseUp
)

// CursorStatus indicates the desired cursor appearance.
type CursorStatus int

const (
	CursorDefault CursorStatus = C.xplm_CursorDefault
	CursorHidden  CursorStatus = C.xplm_CursorHidden
	CursorArrow   CursorStatus = C.xplm_CursorArrow
	CursorCustom  CursorStatus = C.xplm_CursorCustom
)

// WindowLayer specifies the Z-order of a window.
type WindowLayer int

const (
	WindowLayerFlightOverlay      WindowLayer = C.xplm_WindowLayerFlightOverlay
	WindowLayerFloatingWindows    WindowLayer = C.xplm_WindowLayerFloatingWindows
	WindowLayerModal              WindowLayer = C.xplm_WindowLayerModal
	WindowLayerGrowlNotifications WindowLayer = C.xplm_WindowLayerGrowlNotifications
)

// WindowDecoration determines the visual style of a window.
type WindowDecoration int

const (
	DecorationNone                   WindowDecoration = C.xplm_WindowDecorationNone
	DecorationRoundRectangle         WindowDecoration = C.xplm_WindowDecorationRoundRectangle
	DecorationSelfDecorated          WindowDecoration = C.xplm_WindowDecorationSelfDecorated
	DecorationSelfDecoratedResizable WindowDecoration = C.xplm_WindowDecorationSelfDecoratedResizable
)

// WindowPositioningMode determines how a window is positioned on screen.
type WindowPositioningMode int

const (
	PositionFree                WindowPositioningMode = C.xplm_WindowPositionFree
	PositionCenterOnMonitor     WindowPositioningMode = C.xplm_WindowCenterOnMonitor
	PositionFullScreenOnMonitor WindowPositioningMode = C.xplm_WindowFullScreenOnMonitor
	PositionPopOut              WindowPositioningMode = C.xplm_WindowPopOut
	PositionVR                  WindowPositioningMode = C.xplm_WindowVR
)
