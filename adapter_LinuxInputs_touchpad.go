package main

import (
	"github.com/kenshaw/evdev"
)

type Pos [2][2]int32

func initInputAdapter_LinuxInputs_Touchpad(mk mouseKeyboard, hotPlug bool, patern string) {
	//接收触摸板事件作为鼠标输出
	//支持点击，滑动，触摸板按下，双指滚动
	eventsCh := make(chan *eventPack, 10) // 增加缓冲
	go autoDetectAndRead(eventsCh, patern, hotPlug, map[dev_type]bool{type_touch_pad: true})

	button_state := map[uint16]int32{
		BTN_TOUCH:          UP,
		BTN_TOOL_FINGER:    UP,
		BTN_TOOL_DOUBLETAP: UP,
	}

	current_finger := int32(0)
	pos := Pos{{0, 0}, {0, 0}}

	moveMouse := makeScaledMover(mk.MouseMove, 0.3, 0.3, 0.03)

	handelKeyEvents := func(events []*evdev.Event) {
		for _, event := range events {
			switch event.Code {
			case BTN_TOUCH, BTN_TOOL_FINGER, BTN_TOOL_DOUBLETAP:
				button_state[event.Code] = event.Value
			case BTN_MOUSE:
				if event.Value == DOWN {
					mk.MouseBtnDown(MouseBtnLeft)
				} else {
					mk.MouseBtnUp(MouseBtnLeft)
				}
			}
		}
	}

	handelAbsEvents := func(events []*evdev.Event) {
		last_pos := pos
		for _, event := range events {
			logger.Debugf("abs event : %v", event)
			switch event.Code {
			case ABS_MT_SLOT:
				current_finger = event.Value
			case ABS_MT_TRACKING_ID:
				if event.Value == -1 {
					pos[current_finger] = [2]int32{-1, -1}
				}
			case ABS_MT_POSITION_X:
				pos[current_finger][0] = event.Value
			case ABS_MT_POSITION_Y:
				pos[current_finger][1] = event.Value
			}
		}
		if button_state[BTN_TOUCH] == DOWN {
			if button_state[BTN_TOOL_FINGER] == DOWN {
				active_finger := current_finger
				if pos[active_finger][0] == -1 || pos[active_finger][1] == -1 {
					for i := 0; i < 2; i++ {
						if pos[i][0] != -1 && pos[i][1] != -1 {
							active_finger = int32(i)
							break
						}
					}
				}

				if last_pos[active_finger][0] == -1 {
					last_pos[active_finger][0] = pos[active_finger][0]
				}
				if last_pos[active_finger][1] == -1 {
					last_pos[active_finger][1] = pos[active_finger][1]
				}
				offset_x := pos[active_finger][0] - last_pos[active_finger][0]
				offset_y := pos[active_finger][1] - last_pos[active_finger][1]
				moveMouse(offset_x, offset_y, 0)
			} else if button_state[BTN_TOOL_DOUBLETAP] == DOWN {
				if last_pos[0][1] == -1 {
					last_pos[0][1] = pos[0][1]
				}
				if last_pos[1][1] == -1 {
					last_pos[1][1] = pos[1][1]
				}

				last_mid_y := (last_pos[0][1] + last_pos[1][1]) / 2
				mid_y := (pos[0][1] + pos[1][1]) / 2
				vertical_scroll := mid_y - last_mid_y
				moveMouse(0, 0, vertical_scroll)
			}
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
