package robotgo

import (
	"runtime"
	"time"
)

const (
	// Version get the robotgo version
	Version = "v1.00.0.1189, MT. Baker!"
)

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

var (
	// MouseSleep set the mouse default millisecond sleep time
	MouseSleep = 0
	// KeySleep set the key default millisecond sleep time
	KeySleep = 10

	// DisplayID set the screen display id
	DisplayID = -1

	// NotPid used the hwnd not pid in windows
	NotPid bool
	// Scale option the os screen scale
	Scale bool
)

// Bitmap define the go Bitmap struct
//
// The common type conversion of bitmap:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#type-conversion
type Bitmap struct {
	ImgBuf        *uint8
	Width, Height int

	Bytewidth     int
	BitsPixel     uint8
	BytesPerPixel uint8
}

// Point is point struct
type Point struct {
	X int
	Y int
}

// Size is size structure
type Size struct {
	W, H int
}

// Rect is rect structure
type Rect struct {
	Point
	Size
}

// Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
}

// ToInterfaces convert []string to []interface{}
func ToInterfaces(fields []string) []interface{} {
	res := make([]interface{}, 0, len(fields))
	for _, s := range fields {
		res = append(res, s)
	}
	return res
}

// ToStrings convert []interface{} to []string
func ToStrings(fields []interface{}) []string {
	res := make([]string, 0, len(fields))
	for _, s := range fields {
		res = append(res, s.(string))
	}
	return res
}

// CmdCtrl If the operating system is macOS, return the key string "cmd",
// otherwise return the key string "ctrl
func CmdCtrl() string {
	if runtime.GOOS == "darwin" {
		return "cmd"
	}
	return "ctrl" // Ctrl
}

// CmdV tap key command + v or control + v
func CmdV() error {
	return KeyTap("v", CmdCtrl())
}
