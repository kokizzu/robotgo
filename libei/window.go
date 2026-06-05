//go:build linux
// +build linux

package libei

// Window management is not exposed by the RemoteDesktop portal. The portal's
// security model deliberately prevents clients from enumerating, focusing, or
// resizing other applications' windows. These functions exist for robotgo API
// parity and are no-ops / report ErrNotSupported.

// GetTitle returns an empty string; window titles are unavailable on this backend.
func GetTitle(args ...int) string { return "" }

// ActiveName is not supported by this backend.
func ActiveName(name string) error { return ErrNotSupported }

// MinWindow is a no-op on this backend.
func MinWindow(pid int, args ...interface{}) {}

// MaxWindow is a no-op on this backend.
func MaxWindow(pid int, args ...interface{}) {}

// CloseWindow is a no-op on this backend.
func CloseWindow(args ...int) {}
