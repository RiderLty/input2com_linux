package main

import (
	"fmt"

	"strconv"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"
)

type MakcuHandle struct {
	PortName string
	Port     serial.Port
	// 新增字段：用于按键回调
	buttonCallback    func(MouseButton, bool) // 回调函数，参数为按键枚举和状态 (true=按下)
	lastButtonMask    byte                    // 上一次收到的按键状态
	currentButtonMask byte                    // 当前按键状态（可用于查询）
	listenerRunning   bool                    // 标记监听协程是否在运行
	stopListener      chan struct{}           // 用于通知监听协程停止
	mu                sync.Mutex
	mouseButtonByte   byte
	keyBytes          []byte
}

// Make a connection to the COM port where our MAKCU was found.
func Connect(portName string, baudRate int) (*MakcuHandle, error) {
	port, err := serial.Open(portName, &serial.Mode{
		BaudRate: baudRate,
	})
	if err != nil {
		return nil, err
	}
	return &MakcuHandle{Port: port, PortName: portName}, nil
}

// Close the connection to the MAKCU
func (m *MakcuHandle) Close() error {
	if m == nil {
		return fmt.Errorf("Close: MakcuHandle is nil (no device connected)")
	}

	err := m.Port.Close()
	if err != nil {
		return fmt.Errorf("Close: failed to close handle: %w", err)
	}

	return nil
}

// Sends the bytes needed to change the Baud Rate of the MAKCU to 4m and then returns a new Connection object with the new baud rate
// Note: This is NOT a permanent change and will reset back to the default 115200 baud rate after the MAKCU powers off and then back on again.
func ChangeBaudRate(m *MakcuHandle) (*MakcuHandle, error) {
	if m == nil {
		return nil, fmt.Errorf("ChangeBaudRate: MakcuHandle is nil (no device connected)")
	}

	n, err := m.Write([]byte{0xDE, 0xAD, 0x05, 0x00, 0xA5, 0x00, 0x09, 0x3D, 0x00})
	if err != nil {
		// Always try to close the handle on error
		_ = m.Close()
		return nil, fmt.Errorf("ChangeBaudRate: write error: %w", err)
	}

	if n != 9 {
		_ = m.Close()
		return nil, fmt.Errorf("ChangeBaudRate: wrong number of bytes written (got %d, want 9)", n)
	}

	if err := m.Close(); err != nil {
		logger.Logger.Errorf("ChangeBaudRate: failed to close old connection: %v", err)
		// Continue, but log the error
	}

	NewConn, err := Connect(m.PortName, 4000000)
	if err != nil {
		return nil, fmt.Errorf("ChangeBaudRate: connect error: %w", err)
	}

	time.Sleep(1 * time.Second)

	_, err = NewConn.Write([]byte("km.version()\r"))
	if err != nil {
		_ = NewConn.Close()
		return nil, fmt.Errorf("ChangeBaudRate: write error after reconnect: %w", err)
	}

	ReadBuf := make([]byte, 32)
	n, err = NewConn.Read(ReadBuf)
	if err != nil {
		_ = NewConn.Close()
		return nil, fmt.Errorf("ChangeBaudRate: read error after reconnect: %w", err)
	}

	if !strings.Contains(string(ReadBuf[:n]), "MAKCU") {
		_ = NewConn.Close()
		return nil, fmt.Errorf("ChangeBaudRate: did not receive expected response, got: %q", string(ReadBuf[:n]))
	}

	time.Sleep(1 * time.Second)

	logger.Logger.Infof("Successfully Changed Baud Rate To %d!\n", 4000000)

	return NewConn, nil
}

// Sends the given bytes to the MAKCU and returns the number of bytes written.
func (m *MakcuHandle) Write(data []byte) (int, error) {
	if m == nil {
		return -1, fmt.Errorf("write: MakcuHandle is nil (no device connected)")
	}
	return m.Port.Write(data)
}

// Reads data from the MAKCU and saves it to a given buffer then returns the number of bytes read.
func (m *MakcuHandle) Read(buffer []byte) (int, error) {
	return m.Port.Read(buffer)
}

func (m *MakcuHandle) MouseBtnDown(keyCode byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mouseButtonByte |= keyCode
	_, err := m.Write([]byte(MouseKeyDown[keyCode]))
	if err != nil {
		return err
	}
	return nil
}

func (m *MakcuHandle) MouseBtnUp(keyCode byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mouseButtonByte &^= keyCode

	_, err := m.Write([]byte(MouseKeyUp[keyCode]))
	if err != nil {
		return err
	}
	return nil
}
func (m *MakcuHandle) LeftClick() error {
	if m == nil {
		return fmt.Errorf("LeftClick: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.left(1)\r km.left(0)\r"))
	if err != nil {
		logger.Logger.Infof("Failed to click mouse: %v", err)
		return err
	}

	return nil
}

func (m *MakcuHandle) RightDown() error {
	if m == nil {
		return fmt.Errorf("RightDown: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.right(1)\r"))
	if err != nil {
		logger.Logger.Infof("Failed to press mouse: Write Error: %v", err)
		return err
	}

	return nil
}

func (m *MakcuHandle) RightUp() error {
	if m == nil {
		return fmt.Errorf("RightUp: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.right(0)\r"))
	if err != nil {
		logger.Logger.Infof("Failed to release mouse: Write Error: %v", err)
		return err
	}

	return nil
}

func (m *MakcuHandle) RightClick() error {
	if m == nil {
		return fmt.Errorf("RightClick: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.right(1)\r km.right(0)\r"))
	if err != nil {
		logger.Logger.Infof("Failed to right click mouse: %v", err)
		return err
	}

	return nil
}

func (m *MakcuHandle) MiddleDown() error {
	if m == nil {
		return fmt.Errorf("MiddleDown: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.middle(1)\r"))
	if err != nil {
		logger.Logger.Infof("Failed to press middle mouse button: Write Error: %v", err)
		return err
	}

	return nil
}

func (m *MakcuHandle) MiddleUp() error {
	if m == nil {
		return fmt.Errorf("MiddleUp: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.middle(0)\r"))
	if err != nil {
		logger.Logger.Infof("Failed to release middle mouse button: Write Error: %v", err)
		return err
	}

	return nil
}

func (m *MakcuHandle) MiddleClick() error {
	if m == nil {
		return fmt.Errorf("MiddleClick: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.middle(1)\r km.middle(0)\r"))
	if err != nil {
		logger.Logger.Infof("Failed to middle click mouse: %v", err)
		return err
	}

	return nil
}

const (
	MOUSE_BUTTON_LEFT   = 0
	MOUSE_BUTTON_RIGHT  = 1
	MOUSE_BUTTON_MIDDLE = 2
	MOUSE_BUTTON_SIDE1  = 3
	MOUSE_BUTTON_SIDE2  = 4
	MOUSE_X             = 5
	MOUSE_Y             = 6
)

func (m *MakcuHandle) Click(i int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, err := m.Write([]byte(fmt.Sprintf("km.click(%d,1)\r", i)))
	if err != nil {
		logger.Logger.Infof("Failed to middle click mouse: %v", err)
		return err
	}

	return nil
}

// 启用某按键的连发模式
func (m *MakcuHandle) Turbo(i int) error {
	_, err := m.Write([]byte(fmt.Sprintf("km.turbo(%d)\r", i)))
	if err != nil {
		logger.Logger.Infof("Failed to middle click mouse: %v", err)
		return err
	}

	return nil
}
func (m *MakcuHandle) LockMouse(Button int, lock int) error {
	if m == nil {
		return fmt.Errorf("MouseLock: MakcuHandle is nil (no device connected)")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if lock != 1 && lock != 0 {
		return fmt.Errorf("MouseLock: lock must be 1(lock) or 0(unlock)")
	}

	switch Button {
	case MOUSE_BUTTON_LEFT:
		_, err := m.Write([]byte("km.lock_ml(" + strconv.Itoa(lock) + ")\r"))
		if err != nil {
			logger.Logger.Infof("Failed to lock mouse: Write Error: %v", err)
			return err
		}
	case MOUSE_BUTTON_RIGHT:
		_, err := m.Write([]byte("km.lock_mr(" + strconv.Itoa(lock) + ")\r"))
		if err != nil {
			logger.Logger.Infof("Failed to lock mouse: Write Error: %v", err)
			return err
		}
	case MOUSE_BUTTON_MIDDLE:
		_, err := m.Write([]byte("km.lock_mm(" + strconv.Itoa(lock) + ")\r"))
		if err != nil {
			logger.Logger.Infof("Failed to lock mouse: Write Error: %v", err)
			return err
		}
	case MOUSE_BUTTON_SIDE1:
		_, err := m.Write([]byte("km.lock_ms1(" + strconv.Itoa(lock) + ")\r"))
		if err != nil {
			logger.Logger.Infof("Failed to lock mouse: Write Error: %v", err)
			return err
		}
	case MOUSE_BUTTON_SIDE2:
		_, err := m.Write([]byte("km.lock_ms2(" + strconv.Itoa(lock) + ")\r"))
		if err != nil {
			logger.Logger.Infof("Failed to lock mouse: Write Error: %v", err)
			return err
		}
	case MOUSE_X:
		_, err := m.Write([]byte("km.lock_mx(" + strconv.Itoa(lock) + ")\r"))
		if err != nil {
			logger.Logger.Infof("Failed to lock mouse: Write Error: %v", err)
			return err
		}
	case MOUSE_Y:
		_, err := m.Write([]byte("km.lock_my(" + strconv.Itoa(lock) + ")\r"))
		if err != nil {
			logger.Logger.Infof("Failed to lock mouse: Write Error: %v", err)
			return err
		}
	default:
		return fmt.Errorf("invalid mouse button: %d", Button)

	}

	return nil
}

func (m *MakcuHandle) GetButtonStatus() (int, error) {
	if m == nil {
		return -1, fmt.Errorf("MAKCU not connected")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	_, err := m.Write([]byte("km.buttons()\r"))
	if err != nil {
		return -1, err
	}

	buf := make([]byte, 128)
	n, err := m.Read(buf)
	if err != nil {
		return -1, err
	}

	response := string(buf[:n])
	if strings.Contains(response, "1") {
		return 1, nil
	} else if strings.Contains(response, "0") {
		return 0, nil
	}

	return -1, fmt.Errorf("could not parse buttons status")
}

func (m *MakcuHandle) SetButtonStatus(enable bool) error {
	if m == nil {
		return fmt.Errorf("MAKCU not connected")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	var cmd string
	if enable {
		cmd = "km.buttons(1,1)\r"
	} else {
		cmd = "km.buttons(0)\r"
	}

	_, err := m.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to set button status: %v", err)
	}

	logger.Logger.Infof("Button echo %s\n", map[bool]string{true: "enabled", false: "disabled"}[enable])
	return nil
}

func (m *MakcuHandle) MouseMove(dx, dy, wheel int32) error {
	//m.mu.Lock()
	//defer m.mu.Unlock()
	if wheel != 0 {
		_, err := m.Write([]byte(fmt.Sprintf("km.wheel(%d)\r", wheel)))
		if err != nil {
			logger.Logger.Infof("Failed to scroll mouse: Write Error: %v", err)
			return err
		}
	} else {
		_, err := m.Write([]byte(fmt.Sprintf("km.move(%d, %d)\r", dx, dy)))
		if err != nil {
			logger.Logger.Infof("Failed to move mouse: Write Error: %v", err)
			return err
		}
	}
	return nil
}

// use a curve with the built in curve functionality from MAKCU... i THINK this is only on fw v3+ ??? idk don't care to fact check it rn either :)
// "It is common sense that the higher the number of the third parameter, the smoother the curve will be fitted" - from MAKCU/km box docs
func (m *MakcuHandle) MoveMouseWithCurve(x, y int, params ...int) error {
	if m == nil {
		return fmt.Errorf("MoveMouseWithCurve: MakcuHandle is nil (no device connected)")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	var cmd string
	switch len(params) {
	case 0:
		cmd = fmt.Sprintf("km.move(%d, %d)\r", x, y)
	case 1:
		cmd = fmt.Sprintf("km.move(%d, %d, %d)\r", x, y, params[0])
	case 3:
		cmd = fmt.Sprintf("km.move(%d, %d, %d, %d, %d)\r", x, y, params[0], params[1], params[2])
	default:
		logger.Logger.Infof("Invalid number of parameters")
		return fmt.Errorf("invalid number of parameters")
	}

	_, err := m.Write([]byte(cmd))
	if err != nil {
		logger.Logger.Infof("Failed to move mouse with curve: Write Error: %v", err)
		return err
	}

	return nil
}

// ----------------------- 新增：按键名称和枚举映射 -----------------------

// 按键名称映射，对应 Python 的 BUTTON_MAP
var ButtonNameMap = []string{
	"left",   // bit 0
	"right",  // bit 1
	"middle", // bit 2
	"mouse4", // bit 3
	"mouse5", // bit 4
}

type MouseButton int

// 按键枚举映射，对应 Python 的 BUTTON_ENUM_MAP
var ButtonEnumMap = []MouseButton{
	MOUSE_BUTTON_LEFT,   // bit 0
	MOUSE_BUTTON_RIGHT,  // bit 1
	MOUSE_BUTTON_MIDDLE, // bit 2
	MOUSE_BUTTON_SIDE1,  // bit 3 (mouse4)
	MOUSE_BUTTON_SIDE2,  // bit 4 (mouse5)
}

// MouseButton 的 String 方法，用于打印按钮名称
func (b MouseButton) String() string {
	switch b {
	case MOUSE_BUTTON_LEFT:
		return "LEFT"
	case MOUSE_BUTTON_RIGHT:
		return "RIGHT"
	case MOUSE_BUTTON_MIDDLE:
		return "MIDDLE"
	case MOUSE_BUTTON_SIDE1:
		return "MOUSE4"
	case MOUSE_BUTTON_SIDE2:
		return "MOUSE5"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", b)
	}
}

// ----------------------- 修改：listenLoop 函数 -----------------------

func (m *MakcuHandle) ListenLoop() {
	// 使用一个足够大的缓冲区来读取数据
	readBuf := make([]byte, 4096)
	lineBuf := make([]byte, 0, 256) // 用于拼接一行文本
	expectingTextMode := false
	lastByte := byte(0)

	for {
		// 检查是否收到停止信号
		select {
		case <-m.stopListener:
			logger.Logger.Infof("listenLoop received stop signal")
			return
		default:
		}

		// 尝试读取数据
		n, err := m.Port.Read(readBuf)
		if err != nil {
			logger.Logger.Errorf("listenLoop: serial read error: %v", err)
			// 在错误后尝试重新连接或退出，这里简单处理为退出
			m.listenerRunning = false
			return
		}

		if n == 0 {
			time.Sleep(1 * time.Millisecond) // 避免空转
			continue
		}

		// 处理读取到的每个字节
		for i := 0; i < n; i++ {
			byteVal := readBuf[i]

			// 处理回车符 (CR, 0x0D)
			if byteVal == 0x0D {
				lastByte = byteVal
				// 在文本模式下，CR 可能是 CRLF 的一部分，先缓存
				if expectingTextMode && len(lineBuf) < cap(lineBuf) {
					lineBuf = append(lineBuf, byteVal)
				}
				continue
			}

			// 处理换行符 (LF, 0x0A)
			if byteVal == 0x0A {
				// 情况1: 前面是 CR，构成 CRLF，完成文本行
				if lastByte == 0x0D {
					// 将 CR 也加入行缓冲区
					if len(lineBuf) < cap(lineBuf) {
						lineBuf = append(lineBuf, 0x0D)
					}
					// 完成行
					if len(lineBuf) > 0 {
						m.processTextLine(lineBuf)
						lineBuf = lineBuf[:0] // 清空
					}
					expectingTextMode = false
				} else if expectingTextMode || len(lineBuf) > 0 {
					// 情况2: 独立的 LF，且在文本模式或有缓存，也视为文本行结束
					if len(lineBuf) > 0 {
						m.processTextLine(lineBuf)
						lineBuf = lineBuf[:0]
					}
					expectingTextMode = false
				} else {
					// 情况3: 独立的 LF，且不在文本模式、无缓存 -> 可能是按键状态数据
					m.handleButtonData(byteVal)
					expectingTextMode = false
				}
				lastByte = byteVal
				continue
			}

			// 处理其他字节
			// 如果是可打印字符或制表符，则认为是文本开始
			if byteVal >= 32 || byteVal == 0x09 {
				expectingTextMode = true
				if len(lineBuf) < cap(lineBuf) {
					lineBuf = append(lineBuf, byteVal)
				}
			} else {
				// 非打印字符，可能是按键状态
				// 特别处理延迟的 CR
				if lastByte == 0x0D {
					m.handleButtonData(0x0D)
				}
				m.handleButtonData(byteVal)
				expectingTextMode = false
				lineBuf = lineBuf[:0] // 清空可能的错误缓存
			}

			lastByte = byteVal
		}
	}
}

// processTextLine 处理接收到的一行文本（如命令响应）。
// 这个方法可以用来解析 ">>> OK" 或其他命令的返回值。
func (m *MakcuHandle) processTextLine(line []byte) {
	// 去除可能的 CRLF
	content := strings.TrimRight(string(line), "\r\n")
	// 这里可以添加逻辑来匹配和处理带 ID 的命令响应
	// 例如，检查 content 是否包含 "#" 来识别是哪个命令的响应
	// 但为了简化，我们只打印
	logger.Logger.Debugf("Received text: %s", content)
}

// ----------------------- 修改：handleButtonData 函数 -----------------------

// handleButtonData 处理一个代表按键状态的字节。
func (m *MakcuHandle) handleButtonData(buttonMask byte) {
	// 如果状态没有变化，直接返回
	if buttonMask == m.lastButtonMask {
		return
	}

	// 计算变化的位
	changedBits := buttonMask ^ m.lastButtonMask
	logger.Logger.Debugf("Button state changed: 0x%02X -> 0x%02X", m.lastButtonMask, buttonMask)

	// 遍历每一位，检查哪些按键状态发生了变化
	for bit := uint(0); bit < 8; bit++ {
		if changedBits&(1<<bit) != 0 {
			isPressed := (buttonMask & (1 << bit)) != 0

			// 根据位索引获取对应的按钮名称和枚举
			var buttonName string
			var button MouseButton

			if int(bit) < len(ButtonNameMap) {
				buttonName = ButtonNameMap[bit]
			} else {
				buttonName = fmt.Sprintf("bit%d", bit)
			}

			if int(bit) < len(ButtonEnumMap) {
				button = ButtonEnumMap[bit]
			} else {
				// 对于超出映射范围的位，可以忽略或使用一个占位符
				logger.Logger.Infof("Unknown button bit: %d (%s)", bit, buttonName)
				continue
			}

			// 更新内部状态
			if isPressed {
				m.currentButtonMask |= 1 << bit
			} else {
				m.currentButtonMask &^= 1 << bit
			}

			// 记录按钮状态变化
			logger.Logger.Debugf("Button %s (%v): %s", buttonName, byte(button),
				map[bool]string{true: "PRESSED", false: "RELEASED"}[isPressed])

			// 调用用户设置的回调函数
			if m.buttonCallback != nil {
				// 在回调中使用 goroutine 避免阻塞监听循环
				go func(btn MouseButton, pressed bool) {
					defer func() {
						if r := recover(); r != nil {
							logger.Logger.Errorf("Button callback panicked: %v", r)
						}
					}()
					m.buttonCallback(btn, pressed)
				}(button, isPressed)
			}
		}
	}

	// 更新上一次的状态
	m.lastButtonMask = buttonMask
}

// ----------------------- 新增：GetButtonStates 方法 -----------------------

// SetButtonCallback 设置一个回调函数，当鼠标按键状态改变时会被调用。
// 回调函数接收两个参数：MouseButton（哪个按键）和 bool（true=按下，false=释放）。
// 传入 nil 可以取消回调。
func (m *MakcuHandle) SetButtonCallback(callback func(MouseButton, bool)) {
	if m == nil {
		return
	}
	m.buttonCallback = callback
	logger.Logger.Infof("Button callback %s", map[bool]string{true: "set", false: "cleared"}[callback != nil])
}

// GetButtonStates 返回一个映射，显示当前所有按钮的按下状态
// 对应 Python 的 get_button_states() 方法
func (m *MakcuHandle) GetButtonStates() map[string]bool {
	states := make(map[string]bool)
	for bit := uint(0); bit < uint(len(ButtonNameMap)); bit++ {
		if m.currentButtonMask&(1<<bit) != 0 {
			states[ButtonNameMap[bit]] = true
		} else {
			states[ButtonNameMap[bit]] = false
		}
	}
	return states
}
func (m *MakcuHandle) IsMouseBtnPressed(keyCode byte) bool {
	return (m.currentButtonMask & keyCode) != 0
}
func (m *MakcuHandle) KeyDown(keyCode byte) error {
	return nil
}
func (m *MakcuHandle) KeyUp(keyCode byte) error {
	return nil
}

// ----------------------- 新增：GetButtonMask 方法 -----------------------

// GetButtonMask 返回当前的按钮状态掩码
// 对应 Python 的 get_button_mask() 方法
func (m *MakcuHandle) GetButtonMask() byte {
	return m.currentButtonMask
}
