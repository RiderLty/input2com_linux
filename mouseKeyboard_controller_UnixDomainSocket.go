package main

import (
	"sync"

	"go.bug.st/serial"
)

type UnixDomainSocketMouseKeyboard struct {
	port            serial.Port
	mouseButtonByte byte
	keyBytes        []byte
	mu              sync.Mutex
}

func NewMouseKeyboard_UnixDomainSocket(mousDevPath string, keyboardDevPath string) *UnixDomainSocketMouseKeyboard {
	return &UnixDomainSocketMouseKeyboard{
		mouseButtonByte: 0x00,
		keyBytes:        make([]byte, 6),
	}
}
