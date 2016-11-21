package xinput

import (
	"syscall"
	"unsafe"
)

type State struct {
	PacketNumber uint32 // increased for every controller change
	Gamepad      Gamepad
}

type Gamepad struct {
	Buttons      uint16
	LeftTrigger  uint8
	RightTrigger uint8
	ThumbLX      int16
	ThumbLY      int16
	ThumbRX      int16
	ThumbRY      int16
}

type Vibration struct {
	LeftMotorSpeed  uint16
	RightMotorSpeed uint16
}

// CONTROLLER_MAX valid controller numbers are 0-3
const CONTROLLER_MAX = 4

var (
	// TRIGGER_TRESHOLD Threshold for the left and right triggers (0 to 255)
	TRIGGER_TRESHOLD uint8 = 50

	// LEFT_THUMB_DEADZONE Deadzone for the Left Thumb Stick (-32767 to 32767)
	LEFT_THUMB_DEADZONE int16 = 7849

	// LEFT_THUMB_DEADZONE Deadzone for the Right Thumb Stick (-32767 to 32767)
	RIGHT_THUMB_DEADZONE int16 = 8689
)

// Digital Input
const (
	DPAD_UP        uint = 0x0001
	DPAD_DOWN      uint = 0x0002
	DPAD_LEFT      uint = 0x0004
	DPAD_RIGHT     uint = 0x0008
	START          uint = 0x0010
	BACK           uint = 0x0020
	LEFT_THUMB     uint = 0x0040
	RIGHT_THUMB    uint = 0x0080
	LEFT_SHOULDER  uint = 0x0100
	RIGHT_SHOULDER uint = 0x0200
	BUTTON_A       uint = 0x1000
	BUTTON_B       uint = 0x2000
	BUTTON_X       uint = 0x4000
	BUTTON_Y       uint = 0x8000
)

// Analog to Digital Input
const (
	RIGHT_STICK_UP    uint = 0x10000
	RIGHT_STICK_DOWN  uint = 0x20000
	RIGHT_STICK_LEFT  uint = 0x40000
	RIGHT_STICK_RIGHT uint = 0x80000

	LEFT_STICK_UP    uint = 0x100000
	LEFT_STICK_DOWN  uint = 0x200000
	LEFT_STICK_LEFT  uint = 0x400000
	LEFT_STICK_RIGHT uint = 0x800000

	LEFT_TRIGGER  uint = 0x1000000
	RIGHT_TRIGGER uint = 0x2000000
)

var (
	loadError          error
	procXInputGetState *syscall.Proc
	procXInputSetState *syscall.Proc
)

func init() {
	loadError = load()
}

func load() error {
	dll, err := syscall.LoadDLL("xinput1_4.dll")
	defer func() {
		if err != nil {
			dll.Release()
		}
	}()
	if err != nil {
		dll, err = syscall.LoadDLL("xinput1_3.dll")
		if err != nil {
			dll, err = syscall.LoadDLL("xinput9_1_0.dll")
			return err
		}
	}
	procXInputGetState, err = dll.FindProc("XInputGetState")
	if err != nil {
		return err
	}
	procXInputSetState, err = dll.FindProc("XInputSetState")
	return err
}

// IsLoaded checks if XInput was successfully loaded.
// Other functions of the library may be used safely only if it returns no error.
func IsLoaded() (bool, error) {
	return loadError == nil, loadError
}

// GetState retrieves the current state of the controller including analog buttons as digital.
// The analog inputs (Thumbsticks and Triggers) will be shifted into digitalAnalogButtonState
func GetState(controller uint, state *State, digitalAnalogButtonState *uint) error {
	r, _, _ := procXInputGetState.Call(uintptr(controller), uintptr(unsafe.Pointer(state)))

	if r == 0 {
		return nil
	}

	*digitalAnalogButtonState = uint(state.Gamepad.Buttons)

	analogToDigital(digitalAnalogButtonState, state)

	return syscall.Errno(r)
}

// GetSimpleState retrieves the current state of the controller excluding analog buttons as digital.
func GetSimpleState(controller uint, state *State) error {
	r, _, _ := procXInputGetState.Call(uintptr(controller), uintptr(unsafe.Pointer(state)))

	if r == 0 {
		return nil
	}

	return syscall.Errno(r)
}

// IsDown Example: IsDown(allButtonState, DPAD_UP)
func IsDown(digitalAnalogButtonState uint, button uint) bool {
	return (digitalAnalogButtonState & button) == button
}

func analogToDigital(digitalAnalogButtonState *uint, state *State) {

	//Left Thumbstick
	if state.Gamepad.ThumbLX > LEFT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= LEFT_STICK_RIGHT

	} else if state.Gamepad.ThumbLX < -LEFT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= LEFT_STICK_LEFT
	}

	if state.Gamepad.ThumbLY > LEFT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= LEFT_STICK_UP

	} else if state.Gamepad.ThumbLY < -LEFT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= LEFT_STICK_DOWN
	}

	//Right ThumbStick
	if state.Gamepad.ThumbRX > RIGHT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= RIGHT_STICK_RIGHT

	} else if state.Gamepad.ThumbRX < -RIGHT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= RIGHT_STICK_LEFT
	}

	if state.Gamepad.ThumbRY > RIGHT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= RIGHT_STICK_UP

	} else if state.Gamepad.ThumbRY < -RIGHT_THUMB_DEADZONE {

		*digitalAnalogButtonState |= RIGHT_STICK_DOWN
	}

	//Triggers
	if state.Gamepad.LeftTrigger > TRIGGER_TRESHOLD {
		*digitalAnalogButtonState |= LEFT_TRIGGER
	}

	if state.Gamepad.RightTrigger > TRIGGER_TRESHOLD {
		*digitalAnalogButtonState |= RIGHT_TRIGGER
	}
}

// SetState sets the vibration for the controller.
func SetState(controller uint, vibration *Vibration) error {
	r, _, _ := procXInputSetState.Call(uintptr(controller), uintptr(unsafe.Pointer(vibration)))
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}
