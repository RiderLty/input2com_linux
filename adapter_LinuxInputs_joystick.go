package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kenshaw/evdev"
)

// DeadzoneConfig 固定结构
type DeadzoneConfig struct {
	LS []float64 `json:"LS"`
	RS []float64 `json:"RS"`
}

// AxisInfo 轴的具体信息
type AxisInfo struct {
	Name    string `json:"name"`
	Range   []int  `json:"range"`
	Reverse bool   `json:"reverse"`
}

// ControllerConfig 顶层配置结构 (用于 JSON 解析)
type ControllerConfig struct {
	Deadzone    DeadzoneConfig      `json:"DEADZONE"`
	Abs         map[string]AxisInfo `json:"ABS"`
	Btn         map[string]string   `json:"BTN"`
	MapKeyboard map[string]string   `json:"MAP_KEYBOARD"`
}

// RuntimeControllerConfig 运行时优化的配置结构
type RuntimeControllerConfig struct {
	Deadzone    DeadzoneConfig
	Abs         map[uint16]AxisInfo
	Btn         map[uint16]string
	MapKeyboard map[string]string
}

// 静态映射表，避免重复创建
var friendly_name_2_mouse = map[string]byte{
	"BTN_LEFT":    MouseBtnLeft,
	"BTN_RIGHT":   MouseBtnRight,
	"BTN_MIDDLE":  MouseBtnMiddle,
	"BTN_SIDE":    MouseBtnForward,
	"BTN_EXTRA":   MouseBtnBack,
	"BTN_FORWARD": MouseBtnForward,
	"BTN_BACK":    MouseBtnBack,
	"BTN_TASK":    MouseBtnBack,
}

var DPAD_MAP = [9]int32{
	6, 4, 5, // index 0, 1, 2 (y=-1)
	2, 0, 1, // index 3, 4, 5 (y=0)
	10, 8, 9, // index 6, 7, 8 (y=1)
}

var dpadBitName = [4]string{
	"BTN_DPAD_RIGHT",
	"BTN_DPAD_LEFT",
	"BTN_DPAD_UP",
	"BTN_DPAD_DOWN",
}

func stickMapMouse(val int32) int32 {
	if val > 512 {
		return (val - 512) * (val - 512) * 15 / 262144
	} else {
		return (512 - val) * (val - 512) * 15 / 262144
	}
}

func initInputAdapter_LinuxInputs_Joystick(mk mouseKeyboard, hotPlug bool, patern string) {
	//接收 手柄事件，转为键鼠输出
	eventsCh := make(chan *eventPack, 10) // 增加缓冲
	go autoDetectAndRead(eventsCh, patern, hotPlug, map[dev_type]bool{type_joystick: true})

	joystickInfo := make(map[string]*RuntimeControllerConfig)
	path, _ := exec.LookPath(os.Args[0])
	absPath, _ := filepath.Abs(path)
	workingDir, _ := filepath.Split(absPath)
	joystickInfosDir := filepath.Join(workingDir, "joystickInfos")

	if _, err := os.Stat(joystickInfosDir); os.IsNotExist(err) {
		logger.Warnf("%s 文件夹不存在,没有载入任何手柄配置文件", joystickInfosDir)
	} else {
		files, _ := os.ReadDir(joystickInfosDir)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if len(file.Name()) < 5 || file.Name()[len(file.Name())-5:] != ".json" {
				continue
			}
			content, _ := os.ReadFile(filepath.Join(joystickInfosDir, file.Name()))
			var config ControllerConfig
			err := json.Unmarshal(content, &config)
			if err != nil {
				logger.Errorf("解析配置文件 %s 失败: %v", file.Name(), err)
				continue
			}

			// 转换为运行时配置
			runtimeConfig := &RuntimeControllerConfig{
				Deadzone:    config.Deadzone,
				Abs:         make(map[uint16]AxisInfo, len(config.Abs)),
				Btn:         make(map[uint16]string, len(config.Btn)),
				MapKeyboard: config.MapKeyboard,
			}

			for k, v := range config.Abs {
				code, err := strconv.Atoi(k)
				if err == nil {
					runtimeConfig.Abs[uint16(code)] = v
				}
			}
			for k, v := range config.Btn {
				code, err := strconv.Atoi(k)
				if err == nil {
					runtimeConfig.Btn[uint16(code)] = v
				}
			}

			joystickInfo[file.Name()[:len(file.Name())-5]] = runtimeConfig
			logger.Infof("手柄配置文件已载入 : %s", file.Name()[:len(file.Name())-5])
		}
	}

	// 状态变量
	var LS_X_val int32 = 512
	_ = LS_X_val
	var LS_Y_val int32 = 512
	var RS_X_val int32 = 512
	var RS_Y_val int32 = 512
	var HAT0X_val int32 = 0
	var HAT0Y_val int32 = 0
	var LT_val int32 = 0
	var RT_val int32 = 0
	// 事件处理函数
	handelKeyEvents := func(events []*evdev.Event, config *RuntimeControllerConfig) {
		for _, event := range events {
			keyName, ok := config.Btn[event.Code]
			if !ok {
				continue
			}
			if mappedKeyName, ok := config.MapKeyboard[keyName]; ok {
				if len(mappedKeyName) >= 3 && mappedKeyName[:3] == "BTN" { //鼠标
					mouseCode, ok := friendly_name_2_mouse[mappedKeyName]
					if ok {
						if event.Value == 0 {
							mk.MouseBtnUp(mouseCode)
						} else {
							mk.MouseBtnDown(mouseCode)
						}
					}
				} else {
					if keyCode, ok := friendly_name_2_keycode[mappedKeyName]; ok {
						if event.Value == 0 {
							mk.KeyUp(byte(keyCode))
						} else {
							mk.KeyDown(byte(keyCode))
						}
					}
				}
			}
		}
	}

	handelAbsEvents := func(events []*evdev.Event, config *RuntimeControllerConfig) {
		lastDpadState := HAT0X_val + HAT0Y_val*3 + 4

		for _, event := range events {
			axisInfo, ok := config.Abs[event.Code]
			if !ok {
				continue
			}
			valMini := int32(axisInfo.Range[0])
			valMax := int32(axisInfo.Range[1])

			// 避免除以零
			rangeDiff := valMax - valMini
			if rangeDiff == 0 {
				rangeDiff = 1
			}

			switch axisInfo.Name {
			case "LS_X":
				LS_X_val = ((event.Value - valMini) << 10) / rangeDiff
			case "LS_Y":
				LS_Y_val = ((event.Value - valMini) << 10) / rangeDiff
			case "RS_X":
				RS_X_val = ((event.Value - valMini) << 10) / rangeDiff
			case "RS_Y":
				RS_Y_val = ((event.Value - valMini) << 10) / rangeDiff
			case "LT":
				lastVal := LT_val
				LT_val = ((event.Value - valMini) << 10) / rangeDiff
				if lastVal < 256 && LT_val >= 256 {
					mk.MouseBtnDown(MouseBtnRight)
				} else if lastVal >= 256 && LT_val < 256 {
					mk.MouseBtnUp(MouseBtnRight)
				}
			case "RT":
				lastVal := RT_val
				RT_val = ((event.Value - valMini) << 10) / rangeDiff
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

		nowDpadState := HAT0X_val + HAT0Y_val*3 + 4
		if lastDpadState != nowDpadState {
			justPressed := DPAD_MAP[nowDpadState] &^ DPAD_MAP[lastDpadState]
			justReleased := DPAD_MAP[lastDpadState] &^ DPAD_MAP[nowDpadState]
			for index, bitName := range dpadBitName {
				mask := int32(1 << index)
				if justPressed&mask != 0 {
					if mappedKeyName, ok := config.MapKeyboard[bitName]; ok {
						if code, ok := friendly_name_2_keycode[mappedKeyName]; ok {
							mk.KeyDown(byte(code))
						}
					}
				}
				if justReleased&mask != 0 {
					if mappedKeyName, ok := config.MapKeyboard[bitName]; ok {
						if code, ok := friendly_name_2_keycode[mappedKeyName]; ok {
							mk.KeyUp(byte(code))
						}
					}
				}
			}
		}
	}

	// 鼠标移动处理 Goroutine
	go func() {
		ticker := time.NewTicker(4 * time.Millisecond)
		defer ticker.Stop()
		counter := 0

		for {
			select {
			case <-globalCloseSignal:
				return
			case <-ticker.C:
				counter++
				if counter >= 10 {
					counter = 0
				}
				var dx, dy, dw int32
				// dx = (RS_X_val - 512) * 15 / 512
				// dy = (RS_Y_val - 512) * 15 / 512
				dx = stickMapMouse(RS_X_val)
				dy = stickMapMouse(RS_Y_val)
				if counter == 0 {
					dw = -1 * (LS_Y_val - 512) * 3 / 512
				}
				// 只有当有实际移动或滚动时才调用
				if dx != 0 || dy != 0 || dw != 0 {
					mk.MouseMove(dx, dy, dw)
				}
			}
		}
	}()

	// 主事件循环
	for {
		select {
		case <-globalCloseSignal:
			return
		case eventPack := <-eventsCh:
			if eventPack == nil {
				continue
			}

			config, ok := joystickInfo[eventPack.devName]
			if !ok {
				// logger.Errorf("未找到 %s 的手柄配置文件", eventPack.devName)
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
				handelKeyEvents(keyEvents, config)
			}
			if len(absEvents) > 0 {
				handelAbsEvents(absEvents, config)
			}
		}
	}
}
