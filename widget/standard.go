package widget

// #cgo CFLAGS: -DXPLM410=1
// #cgo LDFLAGS: -lXPWidgets_64
// #include "XPStandardWidgets.h"
import "C"

// Standard Widget Classes
const (
	MainWindowClass WidgetClass = C.xpWidgetClass_MainWindow
	SubWindowClass  WidgetClass = C.xpWidgetClass_SubWindow
	ButtonClass     WidgetClass = C.xpWidgetClass_Button
	TextFieldClass  WidgetClass = C.xpWidgetClass_TextField
	ScrollBarClass  WidgetClass = C.xpWidgetClass_ScrollBar
	CaptionClass    WidgetClass = C.xpWidgetClass_Caption
)

// Button Properties and Messages
const (
	PropertyButtonType        PropertyID    = C.xpProperty_ButtonType
	PropertyButtonBehavior    PropertyID    = C.xpProperty_ButtonBehavior
	PropertyButtonState       PropertyID    = C.xpProperty_ButtonState
	MessagePushButtonPressed  WidgetMessage = C.xpMsg_PushButtonPressed
	MessageButtonStateChanged WidgetMessage = C.xpMsg_ButtonStateChanged
)

// Button Types
const (
	PushButton     = C.xpPushButton
	RadioButton    = C.xpRadioButton
	WindowCloseBox = C.xpWindowCloseBox
)

// Button Behaviors
const (
	ButtonBehaviorPushButton  = C.xpButtonBehaviorPushButton
	ButtonBehaviorCheckBox    = C.xpButtonBehaviorCheckBox
	ButtonBehaviorRadioButton = C.xpButtonBehaviorRadioButton
)

// ScrollBar Properties and Messages
const (
	PropertyScrollBarSliderPosition       PropertyID    = C.xpProperty_ScrollBarSliderPosition
	PropertyScrollBarMin                  PropertyID    = C.xpProperty_ScrollBarMin
	PropertyScrollBarMax                  PropertyID    = C.xpProperty_ScrollBarMax
	MessageScrollBarSliderPositionChanged WidgetMessage = C.xpMsg_ScrollBarSliderPositionChanged
)
