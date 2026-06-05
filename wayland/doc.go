//go:build linux
// +build linux

// Package wayland provides a pure-Go Wayland implementation of the robotgo API
// for desktop automation on wlroots-based compositors (Sway, Hyprland, Wayfire, etc.).
//
// This package requires a running Wayland compositor that supports:
//   - zwlr_virtual_pointer_v1 (mouse control)
//   - zwp_virtual_keyboard_v1 (keyboard control)
//   - zwlr_screencopy_v1 (screen capture)
//   - zwlr_foreign_toplevel_management_v1 (window management)
//
// These protocols are available on wlroots-based compositors.
// GNOME and KDE do NOT support them natively.
package wayland
