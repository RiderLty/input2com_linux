package main

import (
	"sync"

	"go.bug.st/serial"
)

type USBGadgetMouseKeyboard struct {
	port            serial.Port
	mouseButtonByte byte
	keyBytes        []byte
	mu              sync.Mutex
}

func NewMouseKeyboard_USBGadget(mousDevPath string, keyboardDevPath string) *USBGadgetMouseKeyboard {
	return &USBGadgetMouseKeyboard{
		mouseButtonByte: 0x00,
		keyBytes:        make([]byte, 6),
	}
}

func (mk *USBGadgetMouseKeyboard) MouseMove(dx, dy, Wheel int32) error { return nil }
func (mk *USBGadgetMouseKeyboard) MouseBtnDown(keyCode byte) error     { return nil }
func (mk *USBGadgetMouseKeyboard) MouseBtnUp(keyCode byte) error       { return nil }
func (mk *USBGadgetMouseKeyboard) KeyDown(keyCode byte) error          { return nil }
func (mk *USBGadgetMouseKeyboard) KeyUp(keyCode byte) error            { return nil }
