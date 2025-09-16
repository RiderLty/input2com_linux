package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
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

var globalCloseSignal = make(chan bool) //仅会在程序退出时关闭  不用于其他用途

type devType uint8

const (
	typeMouse    = devType(0)
	typeKeyboard = devType(1)
	typeJoystick = devType(2)
	typeTouch    = devType(3)
	typeUnknown  = devType(4)
)

func checkDevType(dev *evdev.Evdev) devType {
	abs := dev.AbsoluteTypes()
	key := dev.KeyTypes()
	rel := dev.RelativeTypes()
	_, MTPositionX := abs[evdev.AbsoluteMTPositionX]
	_, MTPositionY := abs[evdev.AbsoluteMTPositionY]
	_, MTSlot := abs[evdev.AbsoluteMTSlot]
	_, MTTrackingID := abs[evdev.AbsoluteMTTrackingID]
	if MTPositionX && MTPositionY && MTSlot && MTTrackingID {
		return typeTouch //触屏检测这几个abs类型即可
	}
	_, RelX := rel[evdev.RelativeX]
	_, RelY := rel[evdev.RelativeY]
	_, HWheel := rel[evdev.RelativeHWheel]
	_, MouseLeft := key[evdev.BtnLeft]
	_, MouseRight := key[evdev.BtnRight]
	_, MouseMiddle := key[evdev.BtnMiddle]
	if RelX && RelY && HWheel && MouseLeft && MouseRight && MouseMiddle {
		return typeMouse //鼠标 检测XY 滚轮 左右中键
	}
	keyboardKeys := true
	for i := evdev.KeyEscape; i <= evdev.KeyScrollLock; i++ {
		_, ok := key[i]
		keyboardKeys = keyboardKeys && ok
	}
	if keyboardKeys {
		return typeKeyboard //键盘 检测keycode(1-70)
	}

	axisCount := 0
	for i := evdev.AbsoluteX; i <= evdev.AbsoluteRZ; i++ {
		_, ok := abs[i]
		if ok {
			axisCount++
		}
	}
	LsRs := axisCount >= 4

	keyCount := 0
	for i := evdev.BtnA; i <= evdev.BtnZ; i++ {
		_, ok := key[i]
		if ok {
			keyCount++
		}
	}
	ABXY := keyCount >= 4

	if LsRs && ABXY {
		return typeJoystick //手柄 检测LS,RS A,B,X,Y
	}
	return typeUnknown
}

func getPossibleDeviceIndexes(skipList map[int]bool) map[int]devType {
	files, _ := os.ReadDir("/dev/input")
	result := make(map[int]devType)
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
			devType := checkDevType(d)
			if devType != typeUnknown {
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

func autoDetectAndRead(eventChan chan *eventPack, patern string) {
	//自动检测设备并读取 循环检测 自动管理设备插入移除
	devices := make(map[int]bool)
	for {
		select {
		case <-globalCloseSignal:
			return
		default:
			autoDetectResult := getPossibleDeviceIndexes(devices)
			devTypeFriendlyName := map[devType]string{
				typeMouse:    "鼠标",
				typeKeyboard: "键盘",
				typeJoystick: "手柄",
				typeTouch:    "触屏",
				typeUnknown:  "未知",
			}
			for index, devType := range autoDetectResult {
				devName := getDevNameByIndex(index)
				if devName == "input2com-virtual-device" {
					continue //跳过生成的虚拟设备
				}
				re := regexp.MustCompile(patern)
				if !re.MatchString(devName) {
					logger.Debugf("设备名称 %s 不匹配模式 %s", devName, patern)
					continue
				}

				if devType == typeMouse || devType == typeKeyboard || devType == typeJoystick {
					logger.Infof("检测到设备 %s(/dev/input/event%d) : %s", devName, index, devTypeFriendlyName[devType])

					localIndex := index
					go func() {
						devices[localIndex] = true
						devReader(eventChan, localIndex)
						devices[localIndex] = false
					}()

				}
			}
			time.Sleep(time.Duration(400) * time.Millisecond)
		}
	}
}

func main() {
	//如果有参数-n 则测试模式
	parser := argparse.NewParser("input2com", " ")

	var debug = parser.Flag("d", "debug", &argparse.Options{
		Required: false,
		Default:  false,
		Help:     "调试模式",
	})

	var badurate = parser.Int("b", "badurate", &argparse.Options{
		Required: false,
		Help:     "波特率",
		Default:  2000000,
	})

	var ttyPath = parser.String("t", "tty", &argparse.Options{
		Required: false,
		Default:  "/dev/ttyUSB*",
		Help:     "串口设备路径，可以使用通配符来匹配第一个设备",
	})

	var sbDesc = parser.String("", "sbdesc", &argparse.Options{
		Required: false,
		Default:  "",
		Help:     "自定义设备描述符",
	})
	var csDesc = parser.String("", "csdesc", &argparse.Options{
		Required: false,
		Default:  "",
		Help:     "自定义厂商描述符",
	})
	var cpDesc = parser.String("", "cpdesc", &argparse.Options{
		Required: false,
		Default:  "",
		Help:     "自定义产品描述符",
	})
	var xlDesc = parser.String("", "xldesc", &argparse.Options{
		Required: false,
		Default:  "",
		Help:     "自定义序列号描述符",
	})

	var patern = parser.String("p", "pattern", &argparse.Options{
		Required: false,
		Default:  ".*",
		Help:     "捕获设备名称的通配符模式",
	})

	go serve() //启动配置服务器

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *debug {
		logger.WithDebug()
	}

	matches, err := filepath.Glob(*ttyPath)
	if err != nil {
		logger.Fatalf("无法匹配设备路径: %v", err)
	}
	if len(matches) == 0 {
		logger.Fatalf("没有找到匹配的设备路径: %s", *ttyPath)
	}

	devpath := matches[0] // 取第一个匹配的设备路径
	logger.Infof("使用设备路径: %s", devpath)
	logger.Infof("波特率: %d", *badurate)

	eventsCh := make(chan *eventPack) //主要设备事件管道
	go autoDetectAndRead(eventsCh, *patern)
	comKB := NewComMouseKeyboard(devpath, *badurate, *sbDesc, *csDesc, *cpDesc, *xlDesc)
	macroKB := NewMacroMouseKeyboard(comKB)

	handelRelEvent := func(x, y, HWhell, Wheel int32) {
		if x != 0 || y != 0 || HWhell != 0 || Wheel != 0 {
			macroKB.MouseMove(x, y, Wheel)
		}
	}
	handelKeyEvents := func(events []*evdev.Event, devName string) {
		for _, event := range events {
			if event.Value == 0 {
				logger.Debugf("%v 按键释放: %v", devName, event.Code)
				if event.Code == uint16(evdev.BtnLeft) { // 鼠标左键释放
					macroKB.MouseBtnUp(MouseBtnLeft)
				} else if event.Code == uint16(evdev.BtnRight) { // 鼠标右键释放
					macroKB.MouseBtnUp(MouseBtnRight)
				} else if event.Code == uint16(evdev.BtnMiddle) { // 鼠标中键释放
					macroKB.MouseBtnUp(MouseBtnMiddle)
				} else if event.Code == uint16(evdev.BtnSide) { // 鼠标后退键释放
					macroKB.MouseBtnUp(MouseBtnBack)
				} else if event.Code == uint16(evdev.BtnExtra) { // 鼠标前进键释放
					macroKB.MouseBtnUp(MouseBtnForward)
				} else {
					macroKB.KeyUp(event.Code) // 其他按键释放
				}
			} else if event.Value == 1 {
				logger.Debugf("%v 按键按下: %v", devName, event.Code)
				if event.Code == uint16(evdev.BtnLeft) { // 鼠标左键释放
					macroKB.MouseBtnDown(MouseBtnLeft)
				} else if event.Code == uint16(evdev.BtnRight) { // 鼠标右键释放
					macroKB.MouseBtnDown(MouseBtnRight)
				} else if event.Code == uint16(evdev.BtnMiddle) { // 鼠标中键释放
					macroKB.MouseBtnDown(MouseBtnMiddle)
				} else if event.Code == uint16(evdev.BtnSide) { // 鼠标后退键释放
					macroKB.MouseBtnDown(MouseBtnBack)
				} else if event.Code == uint16(evdev.BtnExtra) { // 鼠标前进键释放
					macroKB.MouseBtnDown(MouseBtnForward)
				} else {
					macroKB.KeyDown(event.Code) // 其他按键释放
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

	go func() {
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
	}()

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-exitChan
	close(globalCloseSignal)
	logger.Info("已停止")
	time.Sleep(time.Millisecond * 40)
}
