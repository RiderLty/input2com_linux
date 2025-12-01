package main

import (
	"bytes"
	"sync"
	"time"

	"go.bug.st/serial"
)

func OpenSerialWritePipe(portName string, baudRate int) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}
	return port, nil
}

func intToByte(value int32) byte {
	if value < -128 || value > 127 {
		logger.Error("Value must be in the range of -128 to 127")
		return 0x00 // Return a default value if out of range
	}
	if value >= 0 {
		return byte(value)
	}
	return byte(0x100 + value)
}

type comMouseKeyboard struct {
	port            serial.Port
	mouseButtonByte byte
	keyBytes        []byte
	mu              sync.Mutex
}

func string2bytes(s string) []byte {
	var buf bytes.Buffer
	for _, r := range s {
		if r > 127 {
			logger.Warnf("Warn: Skipping non-ASCII character %U", r)
			continue
		}
		buf.WriteByte(byte(r))
	}
	return buf.Bytes()
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

func NewMouseKeyboard_KCOM5(portName string, baudRate int, sbDesc string, csDesc string, cpDesc string, xlDesc string) *comMouseKeyboard {
	//初始化KCOM5鼠标键盘
	// portName: 串口端口名
	// baudRate: 波特率
	// sbDesc: 设备描述符
	// csDesc: 厂商描述符
	// cpDesc: 产品描述符
	// xlDesc: 序列号描述符
	// 如果描述符为空 则不设置
	port, err := OpenSerialWritePipe(portName, baudRate)
	if err != nil {
		logger.Error("Failed to open serial port")
		return nil
	}
	port.Write([]byte{0x57, 0xAB, 0x02, 0x00, 0x00, 0x00, 0x00})
	port.Write([]byte{0x57, 0xAB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	port.Write([]byte{0x57, 0xAB, 0x50})
	resp := make([]byte, 21)
	port.Read(resp)
	logger.Infof("设备描述符: %v", resp)

	resp = make([]byte, 1024)
	port.Write([]byte{0x57, 0xAB, 0x51})
	port.Read(resp)
	logger.Infof("厂商描述符:%s", resp[3:resp[3]+4])
	resp = make([]byte, 1024)
	port.Write([]byte{0x57, 0xAB, 0x52})
	port.Read(resp)
	logger.Infof("产品描述符:%s", resp[3:resp[3]+4])

	resp = make([]byte, 1024)
	port.Write([]byte{0x57, 0xAB, 0x53})
	port.Read(resp)
	logger.Infof("序列号描述符:%s", resp[3:resp[3]+4])

	if sbDesc != "" {
		wb := []byte{0x12, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x40, 0x6D, 0x04, 0x90, 0xC0, 0x02, 0x22, 0x01, 0x02, 0x03, 0x01} //
		if len(wb) > 18 {
			logger.Warn("设备描述符不得超过20字节")
		} else {
			if len(wb) < 18 { //用0填充到18字节
				padding := make([]byte, 18-len(wb))
				wb = append(wb, padding...)
			}
			descCmd := append([]byte{0x57, 0xAB, 0xA0}, wb...)
			port.Write(descCmd)
			port.Read(resp) //57CDA0
			logger.Infof("设置设备描述符: %v", descCmd)
		}
	}

	if csDesc != "" {
		wb := string2bytes(csDesc)
		if len(wb) > 40 {
			logger.Warn("厂商描述符不得超过40字节")
		} else {
			descCmd := append([]byte{0x57, 0xAB, 0xA1, byte(len(wb))}, wb...)
			port.Write(descCmd)
			port.Read(resp) //57CDA01
		}
	}

	if cpDesc != "" {
		wb := string2bytes(cpDesc)
		if len(wb) > 40 {
			logger.Warn("产品描述符不得超过40字节")
		} else {
			descCmd := append([]byte{0x57, 0xAB, 0xA2, byte(len(wb))}, wb...)
			port.Write(descCmd)
			port.Read(resp) //57CDA01
		}
	}

	if xlDesc != "" {
		wb := string2bytes(xlDesc)
		if len(wb) > 40 {
			logger.Warn("序列号描述符不得超过40字节")
		} else {
			descCmd := append([]byte{0x57, 0xAB, 0xA3, byte(len(wb))}, wb...)
			port.Write(descCmd)
			port.Read(resp) //57CDA01
		}
	}

	if sbDesc != "" || csDesc != "" || cpDesc != "" || xlDesc != "" {
		port.Write([]byte{0x57, 0xAB, 0xAA})
	}

	return &comMouseKeyboard{port: port, mouseButtonByte: 0x00, keyBytes: []byte{0x57, 0xAB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
}

func (mk *comMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()

	counter := 0
	for dx != 0 || dy != 0 || Wheel != 0 {
		if counter != 0 {
			time.Sleep(time.Duration(1) * time.Millisecond)
		}
		stepDx := clamp(dx, -127, 127)
		stepDy := clamp(dy, -127, 127)
		stepWheel := clamp(Wheel, -127, 127)
		_, err := mk.port.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, intToByte(stepDx), intToByte(stepDy), intToByte(stepWheel)})
		if err != nil {
			return err
		}
		dx -= stepDx
		dy -= stepDy
		Wheel -= stepWheel
		counter++
	}
	if counter != 1 {
		time.Sleep(time.Duration(1) * time.Millisecond)
	}
	return nil
}

func (mk *comMouseKeyboard) MouseBtnDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.mouseButtonByte |= keyCode
	_, err := mk.port.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, 0x00, 0x00, 0x00})

	if err != nil {
		return err
	}
	return nil
}

func (mk *comMouseKeyboard) MouseBtnUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.mouseButtonByte &^= keyCode
	_, err := mk.port.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, 0x00, 0x00, 0x00})
	if err != nil {
		return err
	}
	return nil
}

func (mk *comMouseKeyboard) KeyDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	if keyCode >= KeyLeftCtrl && keyCode <= KeyRightGui {
		mk.keyBytes[3] |= SpecialKeysMap[keyCode]
	} else {
		for i := 0; i < 7; i++ {
			if i == 6 {
				return nil // No space to add new key, ignore
			}
			if mk.keyBytes[i+5] == 0x00 {
				mk.keyBytes[i+5] = keyCode
				break
			}
		}
	}
	_, err := mk.port.Write(mk.keyBytes)
	if err != nil {
		return err
	}
	return nil
}

func (mk *comMouseKeyboard) KeyUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	if keyCode >= KeyLeftCtrl && keyCode <= KeyRightGui {
		mk.keyBytes[3] &^= SpecialKeysMap[keyCode]
	} else {
		for i := 0; i < 7; i++ {
			if i == 6 {
				return nil // No space to add new key, ignore
			}
			if mk.keyBytes[i+5] == keyCode {
				mk.keyBytes[i+5] = 0x00
				break
			}
		}
	}
	_, err := mk.port.Write(mk.keyBytes)
	if err != nil {
		return err
	}
	return nil
}
