package main

import (
	"github.com/kenshaw/evdev"
)

func initInputAdapter_LinuxInputs_Touchpad(mk mouseKeyboard, hotPlug bool, patern string) {
	//接收 手柄事件，转为键鼠输出
	eventsCh := make(chan *eventPack, 10) // 增加缓冲
	go autoDetectAndRead(eventsCh, patern, hotPlug, map[dev_type]bool{type_touch_pad: true})

	handelKeyEvents := func(events []*evdev.Event) {
		for _, event := range events {
			logger.Infof("key event : %v", event)
		}
	}

	handelAbsEvents := func(events []*evdev.Event) {
		for _, event := range events {
			logger.Infof("abs event : %v", event)
		}
	}
	// 主事件循环
	for {
		select {
		case <-globalCloseSignal:
			return
		case eventPack := <-eventsCh:
			if eventPack == nil {
				continue
			}
			keyEvents := make([]*evdev.Event, 0, len(eventPack.events))
			absEvents := make([]*evdev.Event, 0, len(eventPack.events))
			for _, event := range eventPack.events {
				switch event.Type {
				case evdev.EventKey:
					keyEvents = append(keyEvents, event)
				case evdev.EventAbsolute:
					absEvents = append(absEvents, event)
				}
			}
			if len(keyEvents) > 0 {
				handelKeyEvents(keyEvents)
			}
			if len(absEvents) > 0 {
				handelAbsEvents(absEvents)
			}
		}
	}
}
