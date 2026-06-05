//go:build windows
// +build windows

// Package win provides a pure-Go Windows implementation of the robotgo API
// for desktop automation: mouse, keyboard, screen capture, window management
// and process enumeration.
//
// Unlike the upstream go-vgo/robotgo, this package uses NO Cgo. It is built
// entirely on top of the Win32 API bindings in github.com/tailscale/win and
// golang.org/x/sys/windows, so it cross-compiles with CGO_ENABLED=0.
//
// The exported API mirrors the Wayland sibling package
// (github.com/go-vgo/robotgo/wayland) so callers can swap
// implementations per platform with a build tag.
package win
