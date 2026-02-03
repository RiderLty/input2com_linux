package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/kenshaw/evdev"
)

func initInputAdapter_LinuxInputs_Joystick(mk mouseKeyboard, hotPlug bool, patern string) {
	//接收 手柄事件，转为键鼠输出
	eventsCh := make(chan *eventPack) //主要设备事件管道
	go autoDetectAndRead(eventsCh, patern, hotPlug, map[devType]bool{typeJoystick: true})

	joystickInfo := make(map[string]*simplejson.Json)
	path, _ := exec.LookPath(os.Args[0])
	abs, _ := filepath.Abs(path)
	workingDir, _ := filepath.Split(abs)
	joystickInfosDir := filepath.Join(workingDir, "joystickInfos")
	if _, err := os.Stat(joystickInfosDir); os.IsNotExist(err) {
		logger.Warnf("%s 文件夹不存在,没有载入任何手柄配置文件", joystickInfosDir)
	} else {
		files, _ := ioutil.ReadDir(joystickInfosDir)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if file.Name()[len(file.Name())-5:] != ".json" {
				continue
			}
			content, _ := ioutil.ReadFile(filepath.Join(joystickInfosDir, file.Name()))
			info, _ := simplejson.NewJson(content)
			joystickInfo[file.Name()[:len(file.Name())-5]] = info
			logger.Infof("手柄配置文件已载入 : %s", file.Name())
		}
	}
	var LS_X_val int32 = 512 //所有轴 归一化到[0~1024]
	var LS_Y_val int32 = 512
	var RS_X_val int32 = 512
	var RS_Y_val int32 = 512
	var HAT0X_val int32 = 0
	var HAT0Y_val int32 = 0
	var LT_val int32 = 0
	var RT_val int32 = 0

	const (
		DOWN int32 = 1
		UP   int32 = 0
	)

	var DPAD_MAP = [9]int32{
		6, 4, 5, // index 0, 1, 2 (y=-1)
		2, 0, 1, // index 3, 4, 5 (y=0)
		10, 8, 9, // index 6, 7, 8 (y=1)
	}

	dpadBitName := [4]string{
		"BTN_DPAD_RIGHT",
		"BTN_DPAD_LEFT",
		"BTN_DPAD_UP",
		"BTN_DPAD_DOWN",
	}

	friendly_name_2_mouse := map[string]byte{
		"BTN_LEFT":    MouseBtnLeft,
		"BTN_RIGHT":   MouseBtnRight,
		"BTN_MIDDLE":  MouseBtnMiddle,
		"BTN_SIDE":    MouseBtnForward,
		"BTN_EXTRA":   MouseBtnBack,
		"BTN_FORWARD": MouseBtnForward,
		"BTN_BACK":    MouseBtnBack,
		"BTN_TASK":    MouseBtnBack,
	}

	handelKeyEvents := func(events []*evdev.Event, devName string) {
		if len(events) == 0 {
			return
		} else {
			for _, event := range events {
				keyCode := fmt.Sprintf("%d", event.Code)
				keyName := joystickInfo[devName].Get("BTN").Get(keyCode).MustString()
				if mappedKeyName, ok := joystickInfo[devName].Get("MAP_KEYBOARD").CheckGet(keyName); ok {
					if "BTN" == mappedKeyName.MustString()[:3] { //鼠标
						mouseCode := friendly_name_2_mouse[mappedKeyName.MustString()]
						if event.Value == 0 {
							mk.MouseBtnUp(mouseCode)
						} else {
							mk.MouseBtnDown(mouseCode)
						}
					} else {
						if event.Value == 0 {
							mk.KeyUp(byte(friendly_name_2_keycode[mappedKeyName.MustString()]))
						} else {
							mk.KeyDown(byte(friendly_name_2_keycode[mappedKeyName.MustString()]))
						}
					}
				}
			}
		}
	}

	go func() {
		counter := 0
		for {
			select {
			case <-globalCloseSignal:
				return
			default:
				// logger.Debugf("LS_X_val : %d, LS_Y_val : %d, RS_X_val : %d, RS_Y_val : %d, LT_val : %d, RT_val : %d", LS_X_val, LS_Y_val, RS_X_val, RS_Y_val, LT_val, RT_val)
				counter += 1
				counter %= 10
				if counter == 0 && LS_X_val > -1 {

					mk.MouseMove((RS_X_val-512)*15/512, (RS_Y_val-512)*15/512, -1 * (LS_Y_val-512)*3/512)
				}else{
					mk.MouseMove((RS_X_val-512)*15/512, (RS_Y_val-512)*15/512, 0)
				}
				time.Sleep(time.Duration(4) * time.Millisecond)
			}
		}
	}()

	handelAbsEvents := func(events []*evdev.Event, devName string) {
		if len(events) == 0 {
			return
		}
		dpad_state := struct{X int32; Y int32}{HAT0X_val,HAT0Y_val}
		for _, event := range events {
			axisInfo := joystickInfo[devName].Get("ABS").Get(fmt.Sprintf("%d", event.Code))
			valMini := int32(axisInfo.Get("range").GetIndex(0).MustInt())
			valMax := int32(axisInfo.Get("range").GetIndex(1).MustInt())
			switch axisInfo.Get("name").MustString() {
			case "LS_X":
				LS_X_val = ((event.Value - valMini) << 10) / (valMax - valMini)
			case "LS_Y":
				LS_Y_val = ((event.Value - valMini) << 10) / (valMax - valMini)
			case "RS_X":
				RS_X_val = ((event.Value - valMini) << 10) / (valMax - valMini)
			case "RS_Y":
				RS_Y_val = ((event.Value - valMini) << 10) / (valMax - valMini)
			case "LT":
				lastVal := LT_val
				LT_val = ((event.Value - valMini) << 10) / (valMax - valMini)
				if lastVal < 256 && LT_val >= 256 {
					mk.MouseBtnDown(MouseBtnRight)
				} else if lastVal >= 256 && LT_val < 256 {
					mk.MouseBtnUp(MouseBtnRight)
				}
			case "RT":
				lastVal := RT_val
				RT_val = ((event.Value - valMini) << 10) / (valMax - valMini)
				if lastVal < 256 && RT_val >= 256 {
					mk.MouseBtnDown(MouseBtnLeft)
				} else if lastVal >= 256 && RT_val < 256 {
					mk.MouseBtnUp(MouseBtnLeft)
				}
			case "HAT0X":
				HAT0X_val = event.Value
			case "HAT0Y":
				HAT0Y_val = event.Value
			}
		}
		dpad_now_state :=  struct{X int32; Y int32}{HAT0X_val,HAT0Y_val}
		lastDpadState := dpad_state.X + dpad_state.Y * 3 + 4
		nowDpadState := dpad_now_state.X + dpad_now_state.Y * 3 + 4
		if lastDpadState != nowDpadState {
			justPressed := DPAD_MAP[nowDpadState] &^ DPAD_MAP[lastDpadState]
			justReleased := DPAD_MAP[lastDpadState] &^ DPAD_MAP[nowDpadState]
			for index , bitName := range dpadBitName{
				if justPressed & (1 << index) != 0 {
					if mappedKeyName, ok := joystickInfo[devName].Get("MAP_KEYBOARD").CheckGet(bitName); ok {
						logger.Debugf("按下 : %v", mappedKeyName.MustString())
						mk.KeyDown(byte(friendly_name_2_keycode[mappedKeyName.MustString()]))
					}
				}
				if justReleased & (1 << index) != 0 {
					if mappedKeyName, ok := joystickInfo[devName].Get("MAP_KEYBOARD").CheckGet(bitName); ok {
						logger.Debugf("松开 : %v", mappedKeyName.MustString())
						mk.KeyUp(byte(friendly_name_2_keycode[mappedKeyName.MustString()]))
					}
				}
			}
		}
	}

	for {
		keyEvents := make([]*evdev.Event, 0)
		absEvents := make([]*evdev.Event, 0)
		select {
		case <-globalCloseSignal:
			return
		case eventPack := <-eventsCh:
			if eventPack == nil {
				continue
			}
			if _, ok := joystickInfo[eventPack.devName]; !ok {
				logger.Errorf("未找到 %s 的手柄配置文件", eventPack.devName)
				continue
			}
			for _, event := range eventPack.events {
				switch event.Type {
				case evdev.EventKey:
					keyEvents = append(keyEvents, event)
				case evdev.EventAbsolute:
					absEvents = append(absEvents, event)
				}
				handelKeyEvents(keyEvents, eventPack.devName)
				handelAbsEvents(absEvents, eventPack.devName)
			}
		}
	}
}
