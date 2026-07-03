//go:build darwin
// +build darwin

// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package darwin

import (
	"math"
	"time"
)

// MouseSleep is the global mouse delay in milliseconds.
var MouseSleep = 0

// mouseButton resolves a robotgo button name to the CoreGraphics down/up
// event types and the mouse-button number.
func mouseButton(btn string) (down, up, dragged uint32, button uint32) {
	switch btn {
	case "right":
		return kCGEventRightMouseDown, kCGEventRightMouseUp, kCGEventRightMouseDragged, kCGMouseButtonRight
	case "center", "middle":
		return kCGEventOtherMouseDown, kCGEventOtherMouseUp, kCGEventOtherMouseDragged, kCGMouseButtonCenter
	default: // "left"
		return kCGEventLeftMouseDown, kCGEventLeftMouseUp, kCGEventLeftMouseDragged, kCGMouseButtonLeft
	}
}

// postMouse creates and posts a single mouse event at point p.
func postMouse(eventType uint32, p CGPoint, button uint32) {
	if !loaded {
		return
	}
	ev := cgEventCreateMouseEvent(0, eventType, p, button)
	postEvent(ev)
}

// Move moves the mouse to absolute position (x, y).
// The optional displayId is accepted for API parity but ignored.
func Move(x, y int, displayId ...int) {
	postMouse(kCGEventMouseMoved, CGPoint{X: float64(x), Y: float64(y)}, kCGMouseButtonLeft)
	mouseDelay()
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) {
	cx, cy := Location()
	Move(cx+x, cy+y)
}

// MoveSmooth moves the mouse smoothly to (x, y) with an ease-in-out curve.
// Optional args: steps (int), sleepMs (int). Returns true on success.
func MoveSmooth(x, y int, args ...interface{}) bool {
	steps := 20
	sleepMs := 5
	if len(args) >= 1 {
		if v, ok := args[0].(int); ok {
			steps = v
		}
	}
	if len(args) >= 2 {
		if v, ok := args[1].(int); ok {
			sleepMs = v
		}
	}
	if steps < 1 {
		steps = 1
	}

	sx, sy := Location()
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		// Ease-in-out cubic.
		if t < 0.5 {
			t = 4 * t * t * t
		} else {
			t = 1 - math.Pow(-2*t+2, 3)/2
		}
		cx := float64(sx) + float64(x-sx)*t
		cy := float64(sy) + float64(y-sy)*t
		Move(int(math.Round(cx)), int(math.Round(cy)))
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	return true
}

// Click clicks a mouse button. Default is the left button.
// Use "left", "right", or "center". Pass true for a double-click.
func Click(args ...interface{}) error {
	button := "left"
	double := false
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			button = v
		case bool:
			double = v
		}
	}

	down, up, _, num := mouseButton(button)
	p := locationPoint()
	count := 1
	if double {
		count = 2
	}
	for i := 0; i < count; i++ {
		postMouse(down, p, num)
		postMouse(up, p, num)
		if i < count-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}
	mouseDelay()
	return nil
}

// Toggle toggles a mouse button down or up.
//
//	Toggle("left")        // press down
//	Toggle("left", "up")  // release
func Toggle(key ...interface{}) error {
	button := "left"
	up := false
	for _, arg := range key {
		if s, ok := arg.(string); ok {
			switch s {
			case "up":
				up = true
			case "down":
				up = false
			default:
				button = s
			}
		}
	}

	downType, upType, _, num := mouseButton(button)
	p := locationPoint()
	if up {
		postMouse(upType, p, num)
	} else {
		postMouse(downType, p, num)
	}
	return nil
}

// MouseDown sends a mouse button down event.
func MouseDown(key ...interface{}) error {
	return Toggle(append(append([]interface{}{}, key...), "down")...)
}

// MouseUp sends a mouse button up event.
func MouseUp(key ...interface{}) error {
	return Toggle(append(append([]interface{}{}, key...), "up")...)
}

// Scroll scrolls the mouse. Positive y scrolls down, negative scrolls up;
// positive x scrolls right, negative scrolls left. Optional arg: delay ms.
func Scroll(x, y int, args ...int) {
	msDelay := 10
	if len(args) > 0 {
		msDelay = args[0]
	}
	if loaded {
		// Create a line-unit scroll event then set the deltas explicitly.
		// The deltas are overridden via CGEventSetIntegerValueField, so the
		// variadic wheel argument's exact value does not matter (and stays
		// correct across amd64/arm64 calling conventions).
		ev := cgEventCreateScrollWheelEvent(0, kCGScrollEventUnitLine, 2, 0)
		if ev != 0 {
			// macOS scroll: positive delta scrolls up, so negate for
			// robotgo's down-positive convention.
			cgEventSetIntegerValueField(ev, kCGScrollWheelEventDeltaAxis1, int64(-y))
			cgEventSetIntegerValueField(ev, kCGScrollWheelEventDeltaAxis2, int64(x))
			postEvent(ev)
		}
	}
	if msDelay > 0 {
		time.Sleep(time.Duration(msDelay) * time.Millisecond)
	}
}

// ScrollDir scrolls in a named direction: "up", "down", "left", "right".
func ScrollDir(x int, direction ...interface{}) {
	dir := "down"
	if len(direction) > 0 {
		if s, ok := direction[0].(string); ok {
			dir = s
		}
	}
	switch dir {
	case "up":
		Scroll(0, -x)
	case "down":
		Scroll(0, x)
	case "left":
		Scroll(-x, 0)
	case "right":
		Scroll(x, 0)
	}
}

// ScrollSmooth scrolls the mouse smoothly by `to` steps, repeating `num`
// times (default 5) with `tm` ms between steps (default 100). An optional
// third arg sets the horizontal offset per step.
func ScrollSmooth(to int, args ...int) {
	num := 5
	if len(args) > 0 {
		num = args[0]
	}
	tm := 100
	if len(args) > 1 {
		tm = args[1]
	}
	tox := 0
	if len(args) > 2 {
		tox = args[2]
	}

	for i := 0; i < num; i++ {
		Scroll(tox, to)
		MilliSleep(tm)
	}
	MilliSleep(MouseSleep)
}

// DragSmooth moves the mouse smoothly while holding the left button down.
func DragSmooth(x, y int, args ...interface{}) {
	btn := "left"
	if len(args) > 0 {
		if s, ok := args[0].(string); ok {
			btn = s
		}
	}
	_ = Toggle(btn, "down")
	time.Sleep(50 * time.Millisecond)

	// Move with the button held so apps see drag events.
	downType, upType, dragType, num := mouseButton(btn)
	_, _ = downType, upType
	sx, sy := Location()
	steps := 20
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		cx := float64(sx) + float64(x-sx)*t
		cy := float64(sy) + float64(y-sy)*t
		postMouse(dragType, CGPoint{X: math.Round(cx), Y: math.Round(cy)}, num)
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	_ = Toggle(btn, "up")
}

// MoveClick moves to (x, y) then clicks.
func MoveClick(x, y int, args ...interface{}) {
	Move(x, y)
	_ = Click(args...)
}

// locationPoint returns the current mouse position as a CGPoint.
func locationPoint() CGPoint {
	if !loaded {
		return CGPoint{}
	}
	ev := cgEventCreate(0)
	if ev == 0 {
		return CGPoint{}
	}
	p := cgEventGetLocation(ev)
	cfRelease(ev)
	return p
}

// Location returns the current mouse position.
func Location() (int, int) {
	p := locationPoint()
	return int(p.X), int(p.Y)
}

// GetMousePos returns the current mouse position. Alias of Location.
func GetMousePos() (int, int) {
	return Location()
}

func mouseDelay() {
	if MouseSleep > 0 {
		time.Sleep(time.Duration(MouseSleep) * time.Millisecond)
	}
}
