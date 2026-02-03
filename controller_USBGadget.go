package main

import (
	"os"
	"sync"
)

type USBGadgetMouseKeyboard struct {
	mouseFile       *os.File
	keyboardFile    *os.File
	mouseButtonByte byte
	modKeyState     byte
	pressedKeys     map[byte]bool
	mu              sync.Mutex
}

func NewMouseKeyboard_USBGadget(mousDevPath string, keyboardDevPath string) *USBGadgetMouseKeyboard {
	logger.Errorf("USBGadgetMouseKeyboard mousDevPath %s keyboardDevPath %s", mousDevPath, keyboardDevPath)

	// 打开鼠标设备文件
	mouseFile, err := os.OpenFile(mousDevPath, os.O_WRONLY, 0666)
	if err != nil {
		logger.Errorf("Failed to open mouse device: %v", err)
		return nil
	}

	// 打开键盘设备文件
	keyboardFile, err := os.OpenFile(keyboardDevPath, os.O_WRONLY, 0666)
	if err != nil {
		logger.Errorf("Failed to open keyboard device: %v", err)
		mouseFile.Close()
		return nil
	}

	return &USBGadgetMouseKeyboard{
		mouseFile:       mouseFile,
		keyboardFile:    keyboardFile,
		mouseButtonByte: 0x00,
		modKeyState:     0x00,
		pressedKeys:     make(map[byte]bool),
	}
}

func (mk *USBGadgetMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()

	// 构建鼠标报告数据
	// 格式: [按钮状态(1字节), X移动(2字节), Y移动(2字节), 滚轮移动(2字节)]
	report := make([]byte, 7)
	report[0] = mk.mouseButtonByte

	// 写入X移动（2字节，little-endian，有符号）
	report[1] = byte(dx & 0xFF)
	report[2] = byte((dx >> 8) & 0xFF)

	// 写入Y移动（2字节，little-endian，有符号）
	report[3] = byte(dy & 0xFF)
	report[4] = byte((dy >> 8) & 0xFF)

	// 写入滚轮移动（2字节，little-endian，有符号）
	report[5] = byte(Wheel & 0xFF)
	report[6] = byte((Wheel >> 8) & 0xFF)

	// 写入鼠标设备文件
	_, err := mk.mouseFile.Write(report)
	if err != nil {
		logger.Errorf("Failed to write mouse move: %v", err)
		return err
	}

	return nil
}
func (mk *USBGadgetMouseKeyboard) MouseBtnDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()

	// 更新鼠标按钮状态
	mk.mouseButtonByte |= keyCode

	// 构建鼠标报告数据
	// 格式: [按钮状态(1字节), X移动(2字节), Y移动(2字节), 滚轮移动(2字节)]
	report := make([]byte, 7)
	report[0] = mk.mouseButtonByte
	// X、Y、滚轮移动都为0
	report[1] = 0
	report[2] = 0
	report[3] = 0
	report[4] = 0
	report[5] = 0
	report[6] = 0

	// 写入鼠标设备文件
	_, err := mk.mouseFile.Write(report)
	if err != nil {
		logger.Errorf("Failed to write mouse button down: %v", err)
		return err
	}

	return nil
}

func (mk *USBGadgetMouseKeyboard) MouseBtnUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()

	// 更新鼠标按钮状态
	mk.mouseButtonByte &^= keyCode

	// 构建鼠标报告数据
	// 格式: [按钮状态(1字节), X移动(2字节), Y移动(2字节), 滚轮移动(2字节)]
	report := make([]byte, 7)
	report[0] = mk.mouseButtonByte
	// X、Y、滚轮移动都为0
	report[1] = 0
	report[2] = 0
	report[3] = 0
	report[4] = 0
	report[5] = 0
	report[6] = 0

	// 写入鼠标设备文件
	_, err := mk.mouseFile.Write(report)
	if err != nil {
		logger.Errorf("Failed to write mouse button up: %v", err)
		return err
	}

	return nil
}
func (mk *USBGadgetMouseKeyboard) KeyDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()

	// 检查是否是修饰键
	if modKey, ok := SpecialKeysMap[keyCode]; ok {
		// 更新修饰键状态
		mk.modKeyState |= modKey
	} else {
		// 检查是否是普通键
		if keyCode <= 0x81 {
			// 检查是否已经达到最大按键数（6个）
			if len(mk.pressedKeys) < 6 {
				mk.pressedKeys[keyCode] = true
			} else {
				return nil
			}
		}
	}

	// 构建键盘报告数据
	// 格式: [修饰键状态(1字节), 保留(1字节), 按键码1(1字节), 按键码2(1字节), 按键码3(1字节), 按键码4(1字节), 按键码5(1字节), 按键码6(1字节)]
	report := make([]byte, 8)
	report[0] = mk.modKeyState
	report[1] = 0x00 // 保留字节

	// 添加普通按键
	i := 2
	for key := range mk.pressedKeys {
		if i < 8 {
			report[i] = key
			i++
		}
	}

	// 剩余位置填充0
	for ; i < 8; i++ {
		report[i] = 0x00
	}

	// 写入键盘设备文件
	_, err := mk.keyboardFile.Write(report)
	if err != nil {
		logger.Errorf("Failed to write key down: %v", err)
		return err
	}

	return nil
}

func (mk *USBGadgetMouseKeyboard) KeyUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()

	// 检查是否是修饰键
	if modKey, ok := SpecialKeysMap[keyCode]; ok {
		// 清除修饰键状态
		mk.modKeyState &^= modKey
	} else {
		// 检查是否是普通键
		if keyCode <= 0x81 {
			// 从pressedKeys中移除
			delete(mk.pressedKeys, keyCode)
		}
	}

	// 构建键盘报告数据
	report := make([]byte, 8)
	report[0] = mk.modKeyState
	report[1] = 0x00 // 保留字节

	// 添加普通按键
	i := 2
	for key := range mk.pressedKeys {
		if i < 8 {
			report[i] = key
			i++
		}
	}

	// 剩余位置填充0
	for ; i < 8; i++ {
		report[i] = 0x00
	}

	// 写入键盘设备文件
	_, err := mk.keyboardFile.Write(report)
	if err != nil {
		logger.Errorf("Failed to write key up: %v", err)
		return err
	}

	return nil
}
