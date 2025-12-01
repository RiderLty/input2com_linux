package main

// 宏拦截器，实现了键盘鼠标接口类，接收另一个键盘鼠标接口作为控制器
// 实现无感的替换动作与控制
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var global_moved_x int64 = 0
var global_moved_y int64 = 0

type macroMouseKeyboard struct {
	mouseBtnArgs map[byte]chan bool
	keyArgs      map[byte]chan bool
	ctrl         mouseKeyboard
	macros       map[string]macro // 存储宏函数
}

func abs(x int32) int32 {
	if x < 0 {
		return x * -1
	} else {
		return x
	}
}

func configInit() {
	// // mouseConfigDict[MouseBtnLeft] = "K437_downdrag"
	// mouseConfigDict[MouseBtnLeft] = "ai_triger"
	// mouseConfigDict[MouseBtnForward] = "btn_left"

	// mouseConfigDict[MouseBtnBack] = "ai_triger_auto"
	// mouseConfigDict[MouseBtnLeft] = "ai_triger_MOUSE_LEFT"
	// mouseConfigDict[MouseBtnForward] = "btn_left"

	// mouseConfigDict[MouseBtnBack] = "ai_triger_juji_auto"

	// mouseConfigDict[MouseBtnBack] = "test_ai_speed"
	// mouseConfigDict[MouseBtnBack] = "K437_downdrag"

	// mouseConfigDict[MouseBtnMiddle] = "test_move_from_file"

	preConfigDict["清空"] = [2]map[byte]string{} // server设置的时候，都是重置然后一条一条设置。
	preConfigDict["弓箭自动扳机"] = [2]map[byte]string{
		{
			MouseBtnBack:    "ai_triger_auto",
			MouseBtnLeft:    "ai_triger_MOUSE_LEFT",
			MouseBtnForward: "btn_left",
		},
		{},
	}
	preConfigDict["老王的PKM"] = [2]map[byte]string{
		{
			MouseBtnLeft:    "老王的PKM",
			MouseBtnForward: "btn_left",
		},
		{},
	}
	preConfigDict["K437"] = [2]map[byte]string{
		{
			MouseBtnLeft:    "K437_downdrag",
			MouseBtnForward: "btn_left",
		},
		{},
	}
}

type macro struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	fn          func(*macroMouseKeyboard, chan bool)
}

var macros = make(map[string]macro)

var drop_move = false

//===========================================================================================================
//529px  1315dot

func (mk *macroMouseKeyboard) move_once_auto() {
	//0 1 2 3
	move_x := udp_ints[0] * 1315 / 529
	move_y := udp_ints[1] * 1315 / 529
	mk.MouseMove(move_x, move_y, 0)
	logger.Errorf("auto move %v,%v", move_x, move_y)
}

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

func NewMouseKeyboard_MacroInterceptor(controler mouseKeyboard) *macroMouseKeyboard {
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

	macros["ai_triger"] = macro{
		Name:        "AI扳机",
		Description: "AI自动扳机",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			counter := 0
			for {
				select {
				case <-ch:
					mk.ctrl.MouseBtnUp(MouseBtnLeft)
					return
				default:
					counter += 1
					time.Sleep(time.Duration(1) * time.Millisecond)
					if counter > 420 && udp_ints[2] != 0 && udp_ints[3] != 0 && float64(abs(udp_ints[0]))/float64(udp_ints[2]) < 0.5 && float64(abs(udp_ints[1]))/float64(udp_ints[3]) < 0.5 {
						counter = 0
						mk.ctrl.MouseBtnUp(MouseBtnLeft)
						time.Sleep(time.Duration(16) * time.Millisecond)
						mk.ctrl.MouseBtnDown(MouseBtnLeft)
						logger.Warnf("%v (%v)", udp_ints, udp_last)
					}
				}

			}

		},
	}

	macros["test_ai_speed"] = macro{
		Name:        "AI扳机测试速度",
		Description: "点击左键一次，然后当识别变化时显示时间差",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			old_res := fmt.Sprintf("%v", udp_ints)
			start := time.Now().UnixMicro()
			logger.Warnf("当前值%v", old_res)
			for {
				if old_res != fmt.Sprintf("%v", udp_ints) {
					end := time.Now().UnixMicro()
					logger.Errorf("used %v ms", (end-start)/1000)
					logger.Warnf("变化为 %v (%v)", udp_ints, udp_last)
					break
				}
			}
			time.Sleep(time.Duration(100) * time.Millisecond)
			<-ch
			mk.ctrl.MouseBtnUp(MouseBtnLeft)
		},
	}

	runningFlag := false
	switchFlasg := make(chan bool)
	releaseFlag := make(chan bool)

	macros["ai_triger_auto"] = macro{
		Name:        "AI扳机全自动",
		Description: "AI自动扳机_弓箭侧键特调_搭配释放使用",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			if !runningFlag {
				runningFlag = true
				mk.ctrl.MouseBtnDown(MouseBtnLeft)
				logger.Warn("开始长按左键")
				counter := 0
				for {
					select {
					case <-ch:
					case <-switchFlasg:
						mk.ctrl.MouseBtnUp(MouseBtnLeft)
						return
					case <-releaseFlag:
						counter = 0
						drop_move = true
						// mk.move_once_auto()
						time.Sleep(time.Duration(3) * time.Millisecond)
						mk.ctrl.MouseBtnUp(MouseBtnLeft)
						drop_move = false
						time.Sleep(time.Duration(16) * time.Millisecond)
						mk.ctrl.MouseBtnDown(MouseBtnLeft)
						logger.Warnf("用户发射")
					default:
						counter += 1
						time.Sleep(time.Duration(1) * time.Millisecond)
						if counter > 420 && udp_ints[2] != 0 && udp_ints[3] != 0 && float64(abs(udp_ints[0]))/float64(udp_ints[2]) < 0.9 && float64(abs(udp_ints[1]))/float64(udp_ints[3]) < 0.9 {
							counter = 0
							drop_move = true
							mk.move_once_auto()
							mk.ctrl.MouseBtnUp(MouseBtnLeft)
							drop_move = false
							time.Sleep(time.Duration(16) * time.Millisecond)
							mk.ctrl.MouseBtnDown(MouseBtnLeft)
							logger.Warnf("%v (%v)", udp_ints, udp_last)
						}
					}

				}

			} else {
				logger.Warnf("停止长按左键")
				runningFlag = false
				switchFlasg <- false
			}

		},
	}

	macros["ai_triger_MOUSE_LEFT"] = macro{
		Name:        "AI扳机全自动释放",
		Description: "AI自动扳机_释放",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			if runningFlag {
				releaseFlag <- false
			}
			<-ch
		},
	}

	macros["ai_triger_juji_auto"] = macro{
		Name:        "AI扳机全自动_狙击枪_开关版",
		Description: "AI自动扳机_狙击枪_开关版",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			if !runningFlag {
				runningFlag = true
				counter := 0
				for {
					select {
					case <-ch:
					case <-switchFlasg:
						return
					case <-releaseFlag:
					default:
						counter += 1
						time.Sleep(time.Duration(1) * time.Millisecond)
						if counter > 300 && udp_ints[2] != 0 && udp_ints[3] != 0 && float64(abs(udp_ints[0]))/float64(udp_ints[2]) < 0.5 && float64(abs(udp_ints[1]))/float64(udp_ints[3]) < 0.5 {
							counter = 0
							mk.ctrl.MouseBtnDown(MouseBtnLeft)
							time.Sleep(time.Duration(16) * time.Millisecond)
							mk.ctrl.MouseBtnUp(MouseBtnLeft)
							logger.Warnf("%v (%v)", udp_ints, udp_last)
						}
					}

				}

			} else {
				runningFlag = false
				switchFlasg <- false
			}

		},
	}

	macros["test_move_from_file"] = macro{
		Name:        "测试移动",
		Description: "测试移动",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			global_moved_x = 0
			global_moved_y = 0
			<-ch
			logger.Errorf("moved(%v,%v)", global_moved_x, global_moved_y)
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
	if err := mk.ctrl.MouseMove(dx, dy, Wheel); err != nil {
		return err
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

func (mk *macroMouseKeyboard) KeyDown(keyCode byte) error {
	return mk.ctrl.KeyDown(Linux2hid[keyCode])
}

func (mk *macroMouseKeyboard) KeyUp(keyCode byte) error {
	return mk.ctrl.KeyUp(Linux2hid[keyCode])
}
