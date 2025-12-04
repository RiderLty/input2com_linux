package main

import (
	"encoding/binary"
	"fmt"

	"github.com/kenshaw/evdev"
)

func initInputAdapter_UDP(mk mouseKeyboard, port int) {
	conn, err := CreateUDPReader(fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.Errorf("ERROR:%v", err)
		return
	}
	eventsCh := conn.ReadChan()
	handelRelEvent := func(x, y, HWhell, Wheel int32) {
		if x != 0 || y != 0 || HWhell != 0 || Wheel != 0 {
			if !drop_move {
				mk.MouseMove(x, y, Wheel)
			}
		}
	}
	handelKeyEvents := func(events []*evdev.Event, devName string) {
		for _, event := range events {
			if event.Value == 0 {
				logger.Debugf("%v 按键释放: %v", devName, event.Code)
				if event.Code == uint16(evdev.BtnLeft) { // 鼠标左键释放
					mk.MouseBtnUp(MouseBtnLeft)
				} else if event.Code == uint16(evdev.BtnRight) { // 鼠标右键释放
					mk.MouseBtnUp(MouseBtnRight)
				} else if event.Code == uint16(evdev.BtnMiddle) { // 鼠标中键释放
					mk.MouseBtnUp(MouseBtnMiddle)
				} else if event.Code == uint16(evdev.BtnSide) { // 鼠标后退键释放
					mk.MouseBtnUp(MouseBtnBack)
				} else if event.Code == uint16(evdev.BtnExtra) { // 鼠标前进键释放
					mk.MouseBtnUp(MouseBtnForward)
				} else {
					mk.KeyUp(byte(event.Code)) // 其他按键释放
				}
			} else if event.Value == 1 {
				logger.Debugf("%v 按键按下: %v", devName, event.Code)
				if event.Code == uint16(evdev.BtnLeft) { // 鼠标左键释放
					mk.MouseBtnDown(MouseBtnLeft)
				} else if event.Code == uint16(evdev.BtnRight) { // 鼠标右键释放
					mk.MouseBtnDown(MouseBtnRight)
				} else if event.Code == uint16(evdev.BtnMiddle) { // 鼠标中键释放
					mk.MouseBtnDown(MouseBtnMiddle)
				} else if event.Code == uint16(evdev.BtnSide) { // 鼠标后退键释放
					mk.MouseBtnDown(MouseBtnBack)
				} else if event.Code == uint16(evdev.BtnExtra) { // 鼠标前进键释放
					mk.MouseBtnDown(MouseBtnForward)
				} else {
					mk.KeyDown(byte(event.Code)) // 其他按键释放
				}
			} else if event.Value == 2 {
				logger.Debugf("%v 按键重复: %v", devName, event.Code)
			}
		}
	}

	handelAbsEvents := func(events []*evdev.Event, devName string) {
		if len(events) == 0 {
			return
		}
		for _, event := range events {
			if event.Type != evdev.EventAbsolute {
				continue
			}
		}
	}

	for {
		keyEvents := make([]*evdev.Event, 0)
		absEvents := make([]*evdev.Event, 0)
		var x int32 = 0
		var y int32 = 0
		var HWhell int32 = 0
		var Wheel int32 = 0
		select {
		case <-globalCloseSignal:
			return
		case pack := <-eventsCh:
			event_count := int(pack[0])
			dev_name := string(pack[event_count*8+2:])
			// dev_type := dev_type(pack[event_count*8+1])
			for i := 0; i < event_count; i++ {
				event := &evdev.Event{
					Type:  evdev.EventType(uint16(binary.LittleEndian.Uint16(pack[8*i+1 : 8*i+3]))),
					Code:  uint16(binary.LittleEndian.Uint16(pack[8*i+3 : 8*i+5])),
					Value: int32(binary.LittleEndian.Uint32(pack[8*i+5 : 8*i+9])),
				}
				switch event.Type {
				case evdev.EventKey:
					keyEvents = append(keyEvents, event)
				case evdev.EventAbsolute:
					absEvents = append(absEvents, event)
				case evdev.EventRelative:
					switch event.Code {
					case uint16(evdev.RelativeX):
						x = event.Value
					case uint16(evdev.RelativeY):
						y = event.Value
					case uint16(evdev.RelativeHWheel):
						HWhell = event.Value
					case uint16(evdev.RelativeWheel):
						Wheel = event.Value
					}
				}
			}
			handelRelEvent(x, y, HWhell, Wheel)
			handelKeyEvents(keyEvents, dev_name)
			handelAbsEvents(absEvents, dev_name)
		}
	}
}
