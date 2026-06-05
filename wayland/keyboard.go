//go:build linux
// +build linux

package wayland

import (
	"errors"
	"time"
)

// KeySleep is the global keyboard delay in milliseconds.
var KeySleep = 10

// Key constants matching robotgo's API.
const (
	KeyA = "a"
	KeyB = "b"
	KeyC = "c"
	KeyD = "d"
	KeyE = "e"
	KeyF = "f"
	KeyG = "g"
	KeyH = "h"
	KeyI = "i"
	KeyJ = "j"
	KeyK = "k"
	KeyL = "l"
	KeyM = "m"
	KeyN = "n"
	KeyO = "o"
	KeyP = "p"
	KeyQ = "q"
	KeyR = "r"
	KeyS = "s"
	KeyT = "t"
	KeyU = "u"
	KeyV = "v"
	KeyW = "w"
	KeyX = "x"
	KeyY = "y"
	KeyZ = "z"

	Backspace = "backspace"
	Delete    = "delete"
	Enter     = "enter"
	Tab       = "tab"
	Esc       = "esc"
	Escape    = "escape"
	Up        = "up"
	Down      = "down"
	Right     = "right"
	Left      = "left"
	Home      = "home"
	End       = "end"
	Pageup    = "pageup"
	Pagedown  = "pagedown"

	F1  = "f1"
	F2  = "f2"
	F3  = "f3"
	F4  = "f4"
	F5  = "f5"
	F6  = "f6"
	F7  = "f7"
	F8  = "f8"
	F9  = "f9"
	F10 = "f10"
	F11 = "f11"
	F12 = "f12"

	Shift    = "shift"
	Ctrl     = "ctrl"
	Alt      = "alt"
	Cmd      = "cmd"
	Lshift   = "lshift"
	Rshift   = "rshift"
	Lctrl    = "lctrl"
	Rctrl    = "rctrl"
	Lalt     = "lalt"
	Ralt     = "ralt"
	Space    = "space"
	Capslock = "capslock"
	Print    = "print"
	Insert   = "insert"
	Menu     = "menu"
)

// Evdev key state values
const (
	keyStateReleased = 0
	keyStatePressed  = 1
)

// KeyTap taps a key (press + release). Optional modifiers or pid.
//
//	robotgo.KeyTap("a")
//	robotgo.KeyTap("a", "ctrl")
//	robotgo.KeyTap("a", "ctrl", "shift")
func KeyTap(key string, args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	if c.keyboard == nil || !c.keymapSet {
		return ErrNotSupported
	}

	modifiers := extractModifiers(args)

	// Press modifiers
	for _, mod := range modifiers {
		code, ok := keyToEvdev(mod)
		if !ok {
			continue
		}
		if err := c.keyboard.Key(timestamp(), code, keyStatePressed); err != nil {
			return err
		}
	}

	// Press and release the key
	code, ok := keyToEvdev(key)
	if !ok {
		return errors.New("robotgo: unknown key: " + key)
	}

	ts := timestamp()
	if err := c.keyboard.Key(ts, code, keyStatePressed); err != nil {
		return err
	}
	time.Sleep(time.Duration(KeySleep) * time.Millisecond)
	if err := c.keyboard.Key(ts+uint32(KeySleep), code, keyStateReleased); err != nil {
		return err
	}

	// Release modifiers (reverse order)
	for i := len(modifiers) - 1; i >= 0; i-- {
		code, ok := keyToEvdev(modifiers[i])
		if !ok {
			continue
		}
		if err := c.keyboard.Key(timestamp(), code, keyStateReleased); err != nil {
			return err
		}
	}

	return nil
}

// KeyToggle toggles a key. Default is "down".
//
//	robotgo.KeyToggle("a")        // press
//	robotgo.KeyToggle("a", "up")  // release
func KeyToggle(key string, args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	if c.keyboard == nil || !c.keymapSet {
		return ErrNotSupported
	}

	state := uint32(keyStatePressed)
	for _, arg := range args {
		if s, ok := arg.(string); ok && s == "up" {
			state = keyStateReleased
		}
	}

	code, ok := keyToEvdev(key)
	if !ok {
		return errors.New("robotgo: unknown key: " + key)
	}

	return c.keyboard.Key(timestamp(), code, state)
}

// KeyDown presses a key down.
func KeyDown(key string, args ...interface{}) error {
	return KeyToggle(key, "down")
}

// KeyUp releases a key.
func KeyUp(key string, args ...interface{}) error {
	return KeyToggle(key, "up")
}

// KeyPress presses a key (down + delay + up).
func KeyPress(key string, args ...interface{}) error {
	return KeyTap(key, args...)
}

// Type types a string character by character.
// Uses evdev key codes for ASCII characters.
func Type(str string, args ...int) {
	for _, ch := range str {
		key := string(ch)
		needShift := false

		// Handle uppercase and shifted characters
		if ch >= 'A' && ch <= 'Z' {
			key = string(ch + 32) // lowercase
			needShift = true
		} else if shifted, ok := shiftedChars[ch]; ok {
			key = shifted
			needShift = true
		}

		if needShift {
			_ = KeyTap(key, "shift")
		} else {
			_ = KeyTap(key)
		}
	}
}

// TypeStr types a string character by character.
// It is an alias of Type, mirroring the robotgo API.
func TypeStr(str string, args ...int) {
	Type(str, args...)
}

// TypeDelay types a string with a per-character delay in milliseconds.
func TypeDelay(str string, delay int) {
	old := KeySleep
	KeySleep = delay
	Type(str)
	KeySleep = old
}

// SetDelay sets both KeySleep and MouseSleep.
func SetDelay(d ...int) {
	delay := 10
	if len(d) > 0 {
		delay = d[0]
	}
	KeySleep = delay
	MouseSleep = delay
}

// CmdCtrl returns "cmd" on macOS, "ctrl" on other platforms.
// On Wayland (Linux), always returns "ctrl".
func CmdCtrl() string {
	return "ctrl"
}

func extractModifiers(args []interface{}) []string {
	var mods []string
	for _, arg := range args {
		if s, ok := arg.(string); ok {
			switch s {
			case "ctrl", "control", "lctrl", "rctrl":
				mods = append(mods, s)
			case "shift", "lshift", "rshift":
				mods = append(mods, s)
			case "alt", "lalt", "ralt":
				mods = append(mods, s)
			case "cmd", "lcmd", "rcmd":
				mods = append(mods, s)
			}
		}
	}
	return mods
}

// shiftedChars maps shifted characters to their base key names.
var shiftedChars = map[rune]string{
	'!': "1", '@': "2", '#': "3", '$': "4", '%': "5",
	'^': "6", '&': "7", '*': "8", '(': "9", ')': "0",
	'_': "-", '+': "=", '{': "[", '}': "]", '|': "\\",
	':': ";", '"': "'", '<': ",", '>': ".", '?': "/",
	'~': "`",
}
