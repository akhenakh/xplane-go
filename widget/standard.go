package widget

// #cgo CFLAGS: -DXPLM410=1
// #cgo LDFLAGS: -lXPWidgets_64
// #include "XPStandardWidgets.h"
// #include "XPWidgetDefs.h"
// #define xpProperty_ObjectClass 100
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

// Common Properties
const (
	PropertyRefcon                      PropertyID    = C.xpProperty_Refcon
	PropertyObjectClass                 PropertyID    = C.xpProperty_ObjectClass // The widget's class
	PropertyMainWindowHasCloseBoxes     PropertyID    = C.xpProperty_MainWindowHasCloseBoxes
	MessageMainWindowCloseButtonPressed WidgetMessage = C.xpMessage_CloseButtonPushed
)

// Button Properties and Messages
const (
	PropertyButtonType        PropertyID    = C.xpProperty_ButtonType
	PropertyButtonBehavior    PropertyID    = C.xpProperty_ButtonBehavior
	PropertyButtonState       PropertyID    = C.xpProperty_ButtonState
	MessagePushButtonPressed  WidgetMessage = C.xpMsg_PushButtonPressed
	MessageButtonStateChanged WidgetMessage = C.xpMsg_ButtonStateChanged
)

// TextField Properties
const (
	PropertyTextFieldPasswordMode PropertyID = C.xpProperty_PasswordMode
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
