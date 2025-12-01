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
