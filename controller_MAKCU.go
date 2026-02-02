package main

import "fmt"

type makcuMouseKeyboard struct {
	makcu makcu
}

func NewMouseKeyboard_MAKCU(makcu makcu) *makcuMouseKeyboard {
	return &makcuMouseKeyboard{makcu: makcu}
}

func (mk *makcuMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	if Wheel != 0 {
		mk.makcu.writer <- []byte(fmt.Sprintf("km.wheel(%d)\r", Wheel))
	}
	if dx != 0 || dy != 0 {
		mk.makcu.writer <- []byte(fmt.Sprintf("km.move(%d,%d)\r", dx, dy))
	}
	return nil
}

func (mk *makcuMouseKeyboard) MouseBtnDown(keyCode byte) error {
	switch keyCode {
	case MouseBtnLeft:
		mk.makcu.writer <- []byte("km.left(1)\r")
	case MouseBtnRight:
		mk.makcu.writer <- []byte("km.right(1)\r")
	case MouseBtnMiddle:
		mk.makcu.writer <- []byte("km.middle(1)\r")
	case MouseBtnBack:
		mk.makcu.writer <- []byte("km.side1(1)\r")
	case MouseBtnForward:
		mk.makcu.writer <- []byte("km.side2(1)\r")
	default:
		return fmt.Errorf("MouseBtnDown: unknown keyCode: %v", keyCode)
	}
	return nil
}

func (mk *makcuMouseKeyboard) MouseBtnUp(keyCode byte) error {
	switch keyCode {
	case MouseBtnLeft:
		mk.makcu.writer <- []byte("km.left(0)\r")
	case MouseBtnRight:
		mk.makcu.writer <- []byte("km.right(0)\r")
	case MouseBtnMiddle:
		mk.makcu.writer <- []byte("km.middle(0)\r")
	case MouseBtnBack:
		mk.makcu.writer <- []byte("km.side1(0)\r")
	case MouseBtnForward:
		mk.makcu.writer <- []byte("km.side2(0)\r")
	default:
		return fmt.Errorf("MouseBtnUp: unknown keyCode: %v", keyCode)
	}
	return nil
}
func (mk *makcuMouseKeyboard) KeyDown(keyCode byte) error {
	mk.makcu.writer <- []byte(fmt.Sprintf("km.down(%d)\r", keyCode))
	return nil
}
func (mk *makcuMouseKeyboard) KeyUp(keyCode byte) error {
	mk.makcu.writer <- []byte(fmt.Sprintf("km.up(%d)\r", keyCode))
	return nil
}
