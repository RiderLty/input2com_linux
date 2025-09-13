package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type macroMouseKeyboard struct {
	mouseBtnArgs map[byte]chan bool
	keyArgs      map[byte]chan bool
	ctrl         *comMouseKeyboard
	macros       map[string]macro // 存储宏函数
}

func clamp(value, min, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func configInit() {
	// mouseConfigDict[MouseBtnLeft] = "K437_downdrag"
	mouseConfigDict[MouseBtnForward] = "btn_left"
}

type macro struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	fn          func(*macroMouseKeyboard, chan bool)
}

var macros = make(map[string]macro)

//===========================================================================================================

func downDragMacroFactory(path string) func(mk *macroMouseKeyboard, ch chan bool) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var result [][4]int32 // 存储结果的二维数组
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line) // 按空格分割每行
		if len(fields) != 4 {
			logger.Errorf("跳过无效行: %s (需要4个数字)\n", line)
			continue
		}
		var arr [4]int32
		for i, field := range fields {
			num, err := strconv.Atoi(field)
			if err != nil {
				logger.Errorf("跳过无效数字: %s\n", field)
				continue
			}
			arr[i] = int32(num)
		}
		result = append(result, arr)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return func(mk *macroMouseKeyboard, ch chan bool) {
		counter := int32(0)
		mk.ctrl.MouseBtnDown(MouseBtnLeft)
		for {
			select {
			case <-ch:
				mk.ctrl.MouseBtnUp(MouseBtnLeft)
				return
			default:
				for _, row := range result {
					if row[0] > counter {
						mk.ctrl.MouseMove(row[1], row[2], 0)
						time.Sleep(time.Duration(row[3]) * time.Millisecond)
						break
					}
				}
				counter++
			}
		}
	}
}

func NewMacroMouseKeyboard(controler *comMouseKeyboard) *macroMouseKeyboard {
	configInit()
	mouseBtnArgs := make(map[byte]chan bool)
	keyArgs := make(map[byte]chan bool)
	for i := 0; i < 8; i++ {
		mouseBtnArgs[byte(1<<i)] = make(chan bool, 1)
	}
	for i := 0; i < 256; i++ {
		keyArgs[byte(i)] = make(chan bool, 1)
	}

	macros["btn_left_hold_autofire"] = macro{
		Name:        "左键按住连发",
		Description: "按住左键 = 连点左键",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			for {
				select {
				case <-ch:
					return
				default:
					mk.ctrl.MouseBtnDown(MouseBtnLeft)
					time.Sleep(8 * time.Millisecond)
					mk.ctrl.MouseBtnUp(MouseBtnLeft)
					time.Sleep(8 * time.Millisecond)
				}
			}
		},
	}

	macros["QBZ95_1_downdrag"] = macro{
		Name:        "qbz配置",
		Description: "QBZ95-1默认预设，站立模式下压枪",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			counter := 0
			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			for {
				select {
				case <-ch:
					mk.ctrl.MouseBtnUp(MouseBtnLeft)
					return
				default:
					if counter < 16 {
						mk.ctrl.MouseMove(0, 9, 0) // 向下拖动
					} else if counter < 25 {
						mk.ctrl.MouseMove(0, 9, 0) // 向下拖动
					} else if counter < 30 {
						mk.ctrl.MouseMove(0, 8, 0) // 向下拖动
					} else {
						mk.ctrl.MouseMove(-3, 8, 0) // 向下拖动
					}
					time.Sleep(30 * time.Millisecond)
					counter++
				}
			}
		}}

	macros["K437_downdrag"] = macro{
		Name:        "K437压枪",
		Description: "K437盲人镜，站立模式下压枪",
		fn:          downDragMacroFactory("./config/K437.txt"),
	}

	macros["QJB201_5x"] = macro{
		Name:        "QJB201_5倍压枪",
		Description: "QJB201 默认配置 5倍镜",
		fn:          downDragMacroFactory("./config/QJB201_5倍.txt"),
	}

	macros["老王的PKM"] = macro{
		Name:        "老王的PKM",
		Description: "老王给的红点PKM 站立压枪",
		fn:          downDragMacroFactory("./config/老王的PKM.txt"),
	}

	macros["mini14_autofire_downdrag"] = macro{
		Name:        "mini14连发+压枪",
		Description: "适用于mini14，左键按住连点+压枪",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			counter := 0
			for {
				select {
				case <-ch:
					return
				default:
					mk.ctrl.MouseBtnDown(MouseBtnLeft)
					time.Sleep(5 * time.Millisecond)
					mk.ctrl.MouseMove(0, 11, 0)
					time.Sleep(5 * time.Millisecond)
					mk.ctrl.MouseBtnUp(MouseBtnLeft)
					time.Sleep(5 * time.Millisecond)
					counter++
				}
			}
		},
	}

	macros["btn_left"] = macro{
		Name:        "左键",
		Description: "普通的左键功能，用于其他按键映射",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			<-ch // 等待信号停止
			mk.ctrl.MouseBtnUp(MouseBtnLeft)
		},
	}

	macros["downdrag_args_from_file"] = macro{
		Name:        "从文件读取压枪数据",
		Description: "读取./test.txt中压枪数据进行压枪",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			file, err := os.Open("./test.txt")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			var result [][4]int32 // 存储结果的二维数组
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				fields := strings.Fields(line) // 按空格分割每行
				if len(fields) != 4 {
					logger.Errorf("跳过无效行: %s (需要4个数字)\n", line)
					continue
				}
				var arr [4]int32
				for i, field := range fields {
					num, err := strconv.Atoi(field)
					if err != nil {
						logger.Errorf("跳过无效数字: %s\n", field)
						continue
					}
					arr[i] = int32(num)
				}
				result = append(result, arr)
			}
			if err := scanner.Err(); err != nil {
				panic(err)
			}

			counter := int32(0)
			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			for {
				select {
				case <-ch:
					mk.ctrl.MouseBtnUp(MouseBtnLeft)
					return
				default:
					for _, row := range result {
						if row[0] > counter {
							mk.ctrl.MouseMove(row[1], row[2], 0)
							time.Sleep(time.Duration(row[3]) * time.Millisecond)
							break
						}
					}
					counter++
				}
			}
		},
	}

	return &macroMouseKeyboard{
		mouseBtnArgs: mouseBtnArgs,
		keyArgs:      keyArgs,
		ctrl:         controler,
		macros:       macros,
	}

}

func (mk *macroMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	// 分别处理 dx, dy, Wheel 的拆分移动
	for dx != 0 || dy != 0 || Wheel != 0 {
		stepDx := clamp(dx, -128, 127)
		stepDy := clamp(dy, -128, 127)
		stepWheel := clamp(Wheel, -128, 127)
		if err := mk.ctrl.MouseMove(stepDx, stepDy, stepWheel); err != nil {
			return err
		}
		dx -= stepDx
		dy -= stepDy
		Wheel -= stepWheel
	}
	return nil
}
func (mk *macroMouseKeyboard) MouseBtnDown(keyCode byte) error {
	value, ok := mouseConfigDict[keyCode]
	if !ok { // 如果没有配置，直接调用控制器的MouseBtnDown
		return mk.ctrl.MouseBtnDown(keyCode)
	} else {
		if macroFunc, exists := mk.macros[value]; exists { // 如果有宏函数，执行宏
			go macroFunc.fn(mk, mk.mouseBtnArgs[keyCode])
			return nil
		}
		return mk.ctrl.MouseBtnDown(keyCode) // 如果没有宏函数，直接调用控制器的MouseBtnDown
	}
}

func (mk *macroMouseKeyboard) MouseBtnUp(keyCode byte) error {
	value, ok := mouseConfigDict[keyCode]
	if !ok { // 如果没有配置，直接调用控制器的MouseBtnDown
		return mk.ctrl.MouseBtnUp(keyCode)
	} else {
		if _, exists := mk.macros[value]; exists { // 如果有宏函数，执行宏
			mk.mouseBtnArgs[keyCode] <- true // 发送信号停止宏
			return nil
		}
		return mk.ctrl.MouseBtnDown(keyCode) // 如果没有宏函数，直接调用控制器的MouseBtnDown
	}
}

func (mk *macroMouseKeyboard) KeyDown(keyCode uint16) error {
	return mk.ctrl.KeyDown(linux2hid[keyCode])
}

func (mk *macroMouseKeyboard) KeyUp(keyCode uint16) error {
	return mk.ctrl.KeyUp(linux2hid[keyCode])
}
