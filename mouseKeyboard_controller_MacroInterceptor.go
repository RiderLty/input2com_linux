package main

// 宏拦截器，实现了键盘鼠标接口类，接收另一个键盘鼠标接口作为控制器
// 实现无感的替换动作与控制
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const IterationDataFilePath = "./config/迭代录制压枪数据.txt"

func writeIterationData(data [][3]int32) error {
	// 创建/清空文件。O_TRUNC 用于清空内容。0644 是文件权限
	f, err := os.OpenFile(IterationDataFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// 格式化数据并写入文件
	// 格式要求：movex,movey，空格隔开的多行文本
	for _, move := range data {
		// move[0] 是 moveX, move[1] 是 moveY
		// move[2] 是 index，不写入文件
		// 使用 fmt.Fprintf 写入文件
		_, err := fmt.Fprintf(f, "%d %d\n", move[0], move[1])
		if err != nil {
			return err
		}
	}
	return nil
}

var global_moved_x int64 = 0
var global_moved_y int64 = 0

type macroMouseKeyboard struct {
	mouseBtnArgs map[byte]chan bool
	keyArgs      map[byte]chan bool
	ctrl         mouseKeyboard
	macros       map[string]macro // 存储宏函数
	iterLock     sync.Mutex       // 迭代压枪锁
	iterState    bool             // 迭代压枪状态
	iterChan     chan [2]int32    // 迭代压枪数据通道
	iterLast     [][3]int32       // 迭代压枪数据通道，记录上一次的移动数据
}

func abs(x int32) int32 {
	if x < 0 {
		return x * -1
	} else {
		return x
	}
}

func configInit() {
	preConfigDict["清空"] = [2]map[byte]string{} // server设置的时候，都是重置然后一条一条设置。
	preConfigDict["弓箭自动扳机"] = [2]map[byte]string{
		{
			MouseBtnBack:    "ai_triger_auto",
			MouseBtnLeft:    "ai_triger_MOUSE_LEFT",
			MouseBtnForward: "btn_left",
		},
		{},
	}
	preConfigDict["手游_PKM"] = [2]map[byte]string{
		{
			MouseBtnLeft:    "手游_PKM",
			MouseBtnForward: "btn_left",
		},
		{},
	}
	preConfigDict["手游_M7"] = [2]map[byte]string{
		{
			MouseBtnLeft:    "手游_M7",
			MouseBtnForward: "btn_left",
		},
		{},
	}

	preConfigDict["PC_M7"] = [2]map[byte]string{
		{
			MouseBtnLeft:    "PC_M7",
			MouseBtnForward: "btn_left",
		},
		{},
	}

	preConfigDict["测量用"] = [2]map[byte]string{
		{
			MouseBtnMiddle: "test_move",
		},
		{},
	}

	preConfigDict["迭代测量压枪"] = [2]map[byte]string{
		{
			MouseBtnLeft:    "rec_down_drag_iter",
			MouseBtnForward: "btn_left",
			MouseBtnBack:    "downdrag_args_from_file",
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

func downDragMacroFactory10ms(path string) func(mk *macroMouseKeyboard, ch chan bool) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var result [][2]int32 // 存储结果的二维数组
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line) // 按空格分割每行
		if len(fields) != 2 {
			logger.Errorf("跳过无效行: %s (需要2个数字)\n", line)
			continue
		}
		var arr [2]int32
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
		counter := 0
		clock := time.NewTicker(10 * time.Millisecond)
		mk.ctrl.MouseBtnDown(MouseBtnLeft)
		for {
			select {
			case <-ch:
				mk.ctrl.MouseBtnUp(MouseBtnLeft)
				return
			case <-clock.C:
				if counter < len(result) {
					mk.ctrl.MouseMove(result[counter][0], result[counter][1], 0)
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

	// macros["K437_downdrag"] = macro{
	// 	Name:        "K437压枪",
	// 	Description: "K437盲人镜，站立模式下压枪",
	// 	fn:          downDragMacroFactory("./config/K437.txt"),
	// }

	// macros["老王的PKM"] = macro{
	// 	Name:        "老王的PKM",
	// 	Description: "老王给的红点PKM 站立压枪",
	// 	fn:          downDragMacroFactory("./config/老王的PKM.txt"),
	// }

	macros["手游_PKM"] = macro{
		Name:        "老王的PKM",
		Description: "老王给的红点PKM 站立压枪",
		fn:          downDragMacroFactory10ms("./config/手游_PKM.txt"),
	}
	macros["手游_M7"] = macro{
		Name:        "手游_M7",
		Description: "手游_M7 站立压枪",
		fn:          downDragMacroFactory10ms("./config/手游_M7.txt"),
	}
	macros["PC_M7"] = macro{
		Name:        "PC_M7",
		Description: "PC_M7 站立压枪",
		fn:          downDragMacroFactory10ms("./config/PC_M7.txt"),
	}

	macros["btn_left"] = macro{
		Name:        "左键",
		Description: "普通的左键功能，用于其他按键映射",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			<-ch
			mk.ctrl.MouseBtnUp(MouseBtnLeft)
		},
	}

	macros["downdrag_args_from_file"] = macro{
		Name:        "从文件读取压枪数据",
		Description: "读取'./config/迭代录制压枪数据.txt'中压枪数据进行压枪",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			logger.Infof("从文件读取压枪数据: %s", "./config/迭代录制压枪数据.txt")
			file, err := os.Open("./config/迭代录制压枪数据.txt")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			var result [][2]int32 // 存储结果的二维数组
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				fields := strings.Fields(line) // 按空格分割每行
				if len(fields) != 2 {
					logger.Errorf("跳过无效行: %s (需要2个数字)\n", line)
					continue
				}
				var arr [2]int32
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
			counter := 0
			clock := time.NewTicker(10 * time.Millisecond)
			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			for {
				select {
				case <-ch:
					mk.ctrl.MouseBtnUp(MouseBtnLeft)
					return
				case <-clock.C:
					if counter < len(result) {
						mk.ctrl.MouseMove(result[counter][0], result[counter][1], 0)
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
						if counter > 420 && udp_ints[2] != 0 && udp_ints[3] != 0 && float64(abs(udp_ints[0]))/float64(udp_ints[2]) < 0.5 && float64(abs(udp_ints[1]))/float64(udp_ints[3]) < 0.5 {
							counter = 0
							drop_move = true
							// mk.move_once_auto()
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

	macros["test_move"] = macro{
		Name:        "测试移动",
		Description: "测试移动",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			global_moved_x = 0
			global_moved_y = 0
			<-ch
			logger.Errorf("moved(%v,%v)", global_moved_x, global_moved_y)
		},
	}

	macros["rec_down_drag_iter"] = macro{
		Name:        "迭代测试压枪数据(10ms切片版)",
		Description: "基于10ms时间片叠加：新数据 = 旧脚本回放 + 手动修正，并保存到文件。",
		fn: func(mk *macroMouseKeyboard, ch chan bool) {
			_, err := os.Stat(IterationDataFilePath)
			if os.IsNotExist(err) {
				logger.Infof("数据文件 %s 不存在。检测为首次运行，已清空上次压枪数据 (mk.iterLast)。", IterationDataFilePath)
				mk.iterLast = nil // 清空 mk.iterLast 确保从零开始录制
			} else if err != nil {
				logger.Errorf("检查压枪数据文件 %s 状态时发生错误: %v", IterationDataFilePath, err)
			}
			newIterData := make([][3]int32, 0, 10000)

			mk.ctrl.MouseBtnDown(MouseBtnLeft)
			mk.iterState = true

			// 2. 核心计时器：固定 10ms 心跳
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop() // 确保函数退出时停止计时器

			// 3. 状态变量
			var tickIndex int = 0            // 当前时间片序号
			var userAccumX, userAccumY int32 // 当前10ms内，用户手动的累积位移

			// 缓存旧数据的长度，避免越界
			lastLen := len(mk.iterLast)

		LOOP:
			for {
				select {
				// --- A. 停止信号 ---
				case <-ch:
					mk.iterState = false
					mk.ctrl.MouseBtnUp(MouseBtnLeft)

					// 1. 优雅地排空剩余通道数据（非阻塞）
				DRAIN:
					for {
						select {
						case <-mk.iterChan:
						default:
							break DRAIN
						}
					}

					// 2. 保存本次叠加后的结果到 mk.iterLast
					if len(newIterData) > 0 {
						mk.iterLast = newIterData
						logger.Infof("迭代录制结束，生成 %d 帧(10ms)数据", len(newIterData))

						// 3. 将新数据写入文件 (新增逻辑)
						if err := writeIterationData(newIterData); err != nil {
							logger.Errorf("写入文件失败: %v", err)
						} else {
							logger.Infof("成功写入 %d 帧(10ms)数据到文件 %s", len(newIterData), IterationDataFilePath)
							logger.Info("移除此文件以开始新的录制")
						}
					} else {
						logger.Warn("本次未录制到任何数据，未更新和保存文件。")
					}

					break LOOP

				// --- B. 处理用户手动修正 (实时响应) ---
				case move := <-mk.iterChan:
					// 1. 实时透传：保证手感无延迟
					mk.ctrl.MouseMove(move[0], move[1], 0)

					// 2. 累积到当前时间片
					userAccumX += move[0]
					userAccumY += move[1]

				// --- C. 时间片心跳 (回放旧数据 + 结算录制) ---
				case <-ticker.C:
					// 1. 获取旧脚本在当前帧的数据 (如果存在)
					var scriptX, scriptY int32 = 0, 0
					if tickIndex < lastLen {
						scriptX = mk.iterLast[tickIndex][0]
						scriptY = mk.iterLast[tickIndex][1]
					}

					// 2. 执行旧脚本回放
					if scriptX != 0 || scriptY != 0 {
						mk.ctrl.MouseMove(scriptX, scriptY, 0)
					}

					// 3. 计算新一帧的数据：新数据 = 旧脚本 + 用户修正
					totalX := scriptX + userAccumX
					totalY := scriptY + userAccumY

					// 4. 存入新数组 [TotalX, TotalY, Index]
					// 即使是 0,0 也要存，保持时间轴对齐
					newIterData = append(newIterData, [3]int32{totalX, totalY, int32(tickIndex)})

					// 5. 重置累积器，进入下一帧
					userAccumX = 0
					userAccumY = 0
					tickIndex++
				}
			}
		},
	}
	return &macroMouseKeyboard{
		mouseBtnArgs: mouseBtnArgs,
		keyArgs:      keyArgs,
		ctrl:         controler,
		macros:       macros,
		iterLock:     sync.Mutex{},
		iterState:    false,
		iterChan:     make(chan [2]int32),
	}

}

func (mk *macroMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	if mk.iterState {
		mk.iterChan <- [2]int32{int32(dx), int32(dy)}
		return nil
	} else {
		global_moved_x += int64(dx)
		global_moved_y += int64(dy)
		if err := mk.ctrl.MouseMove(dx, dy, Wheel); err != nil {
			return err
		}
		return nil
	}
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
