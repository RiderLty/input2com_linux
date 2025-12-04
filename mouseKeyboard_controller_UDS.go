package main

import (
	"sync"

	"github.com/kenshaw/evdev"
)

type UDSMouseKeyboard struct {
	writer UDSWriter
	mu     sync.Mutex
}

func NewMouseKeyboard_UDS(address string) *UDSMouseKeyboard {
	writer, err := CreateUDSWriter(address)
	if err != nil {
		logger.Errorf("CreateUDSWriter error : %v", err)
		return nil
	}
	return &UDSMouseKeyboard{
		writer: *writer,
	}
}

func (mk *UDSMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	events := make([]evdev.Event, 0, 3)
	if dx != 0 {
		events = append(events, evdev.Event{
			Type:  evdev.EventRelative,
			Code:  uint16(evdev.RelativeX),
			Value: dx,
		})
	}
	if dy != 0 {
		events = append(events, evdev.Event{
			Type:  evdev.EventRelative,
			Code:  uint16(evdev.RelativeY),
			Value: dy,
		})
	}
	if Wheel != 0 {
		events = append(events, evdev.Event{
			Type:  evdev.EventRelative,
			Code:  uint16(evdev.RelativeWheel),
			Value: Wheel,
		})
	}
	buffer := append(eventPacker(events), []byte{0, 114, 107, 109}...) // 0(type_mouse) ,  r  k m
	mk.mu.Lock()
	defer mk.mu.Unlock()
	_, err := mk.writer.Write(buffer)
	if err != nil {
		logger.Errorf("UDSMouseKeyboard MouseMove error : %v", err)
		return err
	}
	return nil
}

func (mk *UDSMouseKeyboard) MouseBtnDown(keyCode byte) error {
	var err error = nil
	mk.mu.Lock()
	defer mk.mu.Unlock()
	switch keyCode {
	case MouseBtnLeft:
		_, err = mk.writer.Write([]byte{1, 1, 0, 16, 1, 1, 0, 0, 0, 114, 107, 109})
	case MouseBtnRight:
		_, err = mk.writer.Write([]byte{1, 1, 0, 17, 1, 1, 0, 0, 0, 114, 107, 109})
	case MouseBtnMiddle:
		_, err = mk.writer.Write([]byte{1, 1, 0, 18, 1, 1, 0, 0, 0, 114, 107, 109})
	case MouseBtnForward:
		_, err = mk.writer.Write([]byte{1, 1, 0, 21, 1, 1, 0, 0, 0, 114, 107, 109})
	case MouseBtnBack:
		_, err = mk.writer.Write([]byte{1, 1, 0, 22, 1, 1, 0, 0, 0, 114, 107, 109})
	}
	return err
}

func (mk *UDSMouseKeyboard) MouseBtnUp(keyCode byte) error {
	var err error = nil
	mk.mu.Lock()
	defer mk.mu.Unlock()
	switch keyCode {
	case MouseBtnLeft:
		_, err = mk.writer.Write([]byte{1, 1, 0, 16, 1, 0, 0, 0, 0, 114, 107, 109})
	case MouseBtnRight:
		_, err = mk.writer.Write([]byte{1, 1, 0, 17, 1, 0, 0, 0, 0, 114, 107, 109})
	case MouseBtnMiddle:
		_, err = mk.writer.Write([]byte{1, 1, 0, 18, 1, 0, 0, 0, 0, 114, 107, 109})
	case MouseBtnForward:
		_, err = mk.writer.Write([]byte{1, 1, 0, 21, 1, 0, 0, 0, 0, 114, 107, 109})
	case MouseBtnBack:
		_, err = mk.writer.Write([]byte{1, 1, 0, 22, 1, 0, 0, 0, 0, 114, 107, 109})
	}
	return err
}

func (mk *UDSMouseKeyboard) KeyDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	_, err := mk.writer.Write([]byte{1, 1, 0, hid2linux[keyCode], 0, 1, 0, 0, 0, 1, 114, 107, 109})
	return err
}

func (mk *UDSMouseKeyboard) KeyUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	_, err := mk.writer.Write([]byte{1, 1, 0, hid2linux[keyCode], 0, 0, 0, 0, 0, 1, 114, 107, 109})
	return err
}
