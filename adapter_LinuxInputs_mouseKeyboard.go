package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	"unsafe"

	"github.com/kenshaw/evdev"
)

type eventPack struct {
	//表示一个动作 由一系列event组成
	devName string
	events  []*evdev.Event
}

func devReader(eventReader chan *eventPack, index int) {
	fd, err := os.OpenFile(fmt.Sprintf("/dev/input/event%d", index), os.O_RDONLY, 0)
	if err != nil {
		logger.Errorf("读取设备失败 : %v", err)
		return
	}
	d := evdev.Open(fd)
	defer d.Close()
	eventCh := d.Poll(context.Background())
	events := make([]*evdev.Event, 0)
	devName := d.Name()
	logger.Infof("开始读取设备 : %s", devName)
	d.Lock()
	defer d.Unlock()
	for {
		select {
		case <-globalCloseSignal:
			logger.Infof("释放设备 : %s", devName)
			return
		case event := <-eventCh:
			if event == nil {
				logger.Warnf("移除设备 : %s", devName)
				return
			} else if event.Type == evdev.SyncReport {
				pack := &eventPack{
					devName: devName,
					events:  events,
				}
				eventReader <- pack
				events = make([]*evdev.Event, 0)
			} else {
				events = append(events, &event.Event)
			}
		}
	}
}

func checkDevType(dev *evdev.Evdev, fd *os.File) dev_type {
	abs := dev.AbsoluteTypes()
	key := dev.KeyTypes()
	rel := dev.RelativeTypes()
	_, MTPositionX := abs[evdev.AbsoluteMTPositionX]
	_, MTPositionY := abs[evdev.AbsoluteMTPositionY]
	_, MTSlot := abs[evdev.AbsoluteMTSlot]
	_, MTTrackingID := abs[evdev.AbsoluteMTTrackingID]
	if MTPositionX && MTPositionY && MTSlot && MTTrackingID {
		var bits int32
		ioctl(fd.Fd(), EVIOCGPROP(), uintptr(unsafe.Pointer(&bits)))
		if bits&(1<<_INPUT_PROP_DIRECT) != 0 {
			return type_touch_screen
		}
		if bits&(1<<_INPUT_PROP_POINTER) != 0 {
			return type_touch_pad
		}
		return type_unknown
		// return type_touch_screen //触屏检测这几个abs类型即可
	}
	_, RelX := rel[evdev.RelativeX]
	_, RelY := rel[evdev.RelativeY]
	_, Wheel := rel[evdev.RelativeWheel]
	_, MouseLeft := key[evdev.BtnLeft]
	_, MouseRight := key[evdev.BtnRight]
	if RelX && RelY && Wheel && MouseLeft && MouseRight {
		return type_mouse //鼠标 检测XY 滚轮 左右键
	}
	keyboard_keys := true
	for i := evdev.KeyEscape; i <= evdev.KeyScrollLock; i++ {
		_, ok := key[i]
		keyboard_keys = keyboard_keys && ok
	}
	if keyboard_keys {
		return type_keyboard //键盘 检测keycode(1-70)
	}

	axis_count := len(abs)
	btn_count := len(key)
	if axis_count >= 4 { //检测轴的数量是否有两个摇杆 可能存在误报
		if btn_count == 0 { //如果大于4轴，且没有按键，认为是运动传感器
			return type_motion_sensors
		} else if btn_count > 8 {
			return type_joystick //按键大于8个
		}
	}
	return type_unknown
}

func getPossibleDeviceIndexes(skipList map[int]bool) map[int]dev_type {
	files, _ := os.ReadDir("/dev/input")
	result := make(map[int]dev_type)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if len(file.Name()) <= 5 {
			continue
		}
		if file.Name()[:5] != "event" {
			continue
		}
		index, _ := strconv.Atoi(file.Name()[5:])
		reading, exist := skipList[index]
		if exist && reading {
			continue
		} else {
			fd, err := os.OpenFile(fmt.Sprintf("/dev/input/%s", file.Name()), os.O_RDONLY, 0)
			if err != nil {
				logger.Errorf("读取设备/dev/input/%s失败 : %v ", file.Name(), err)
			}
			d := evdev.Open(fd)
			defer d.Close()
			devType := checkDevType(d, fd)
			if devType != type_unknown {
				result[index] = devType
			}
		}
	}
	return result
}

func getDevNameByIndex(index int) string {
	fd, err := os.OpenFile(fmt.Sprintf("/dev/input/event%d", index), os.O_RDONLY, 0)
	if err != nil {
		return "读取设备名称失败"
	}
	d := evdev.Open(fd)
	defer d.Close()
	return d.Name()
}

func autoDetectAndRead(eventChan chan *eventPack, patern string, loop bool, types map[dev_type]bool) {
	//自动检测设备并读取 循环检测 自动管理设备插入移除
	devices := make(map[int]bool)
	for {
		select {
		case <-globalCloseSignal:
			return
		default:
			autoDetectResult := getPossibleDeviceIndexes(devices)
			devTypeFriendlyName := map[dev_type]string{
				type_mouse:          "鼠标",
				type_keyboard:       "键盘",
				type_joystick:       "手柄",
				type_touch_pad:      "触摸板",
				type_motion_sensors: "运动传感器",
				type_unknown:        "未知",
			}
			for index, devType := range autoDetectResult {
				devName := getDevNameByIndex(index)
				if devName == Cfg.Dst.Uinput.DeviceName {
					continue //跳过生成的虚拟设备
				}
				re := regexp.MustCompile(patern)
				if !re.MatchString(devName) {
					logger.Debugf("设备名称 %s 不匹配模式 %s", devName, patern)
					continue
				}
				if types[devType] {
					logger.Infof("检测到设备 %s(/dev/input/event%d) : %s", devName, index, devTypeFriendlyName[devType])
					localIndex := index
					go func() {
						devices[localIndex] = true
						devReader(eventChan, localIndex)
						devices[localIndex] = false
					}()
				}
			}
			if !loop {
				return
			}
			time.Sleep(time.Duration(400) * time.Millisecond)
		}
	}
}

func initInputAdapter_LinuxInputs_MouseKeyboard(mk mouseKeyboard, hotPlug bool, patern string) {
	//初始化linux输入设备适配器
	// 自动检测设备并读取事件 循环检测 自动管理设备插入移除
	// 事件会使用mk的方法处理
	// 使用携程启用
	eventsCh := make(chan *eventPack) //主要设备事件管道
	go autoDetectAndRead(eventsCh, patern, hotPlug, map[dev_type]bool{type_mouse: true, type_keyboard: true})
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
		case eventPack := <-eventsCh:
			if eventPack == nil {
				continue
			}
			for _, event := range eventPack.events {
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

			var perfPoint time.Time
			perfPoint = time.Now()
			handelRelEvent(x, y, HWhell, Wheel)
			relSin := time.Since(perfPoint)
			perfPoint = time.Now()
			handelKeyEvents(keyEvents, eventPack.devName)
			keySin := time.Since(perfPoint)
			perfPoint = time.Now()
			handelAbsEvents(absEvents, eventPack.devName)
			absSin := time.Since(perfPoint)
			logger.Debugf("")
			logger.Debugf("handel rel_event\t%v \n", relSin)
			logger.Debugf("handel key_events\t%v \n", keySin)
			logger.Debugf("handel abs_events\t%v \n", absSin)
		}
	}
}
