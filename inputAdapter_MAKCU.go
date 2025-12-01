package main

import (
	"fmt"
	"strings"
	"time"

	"go.bug.st/serial"
)

type makcu struct {
	port   serial.Port
	writer chan []byte
	reader chan []byte
}

func getMakcuPort(portName string, baudRate int) (serial.Port, error) {
	conn, err := serial.Open(portName, &serial.Mode{
		BaudRate: baudRate,
	})
	if err != nil {
		return nil, err
	}
	if baudRate == 4000000 {
		// conn.Write([]byte("km.echo(1)\r"))
		conn.Write([]byte("km.version()\r"))
		ReadBuf := make([]byte, 2048)
		n, err := conn.Read(ReadBuf)
		if err != nil {
			logger.Errorf("MAKCU Port open failed with km.version() read error: %w", err)
			return nil, err
		} else {
			logger.Infof("MAKCU Port open success with km.version() response: [%v]", string(ReadBuf[:n]))
		}
		// conn.Write([]byte("km.echo(0)\r"))
		return conn, nil
	} else {
		n, err := conn.Write([]byte{0xDE, 0xAD, 0x05, 0x00, 0xA5, 0x00, 0x09, 0x3D, 0x00})
		if err != nil {
			return nil, err
		}
		if n != 9 {
			conn.Close()
			return nil, fmt.Errorf("返回数据不正确")
		}
		conn.Close()
		NewConn, err := serial.Open(portName, &serial.Mode{
			BaudRate: 4000000,
		})
		if err != nil {
			return nil, err
		}
		// return NewConn, nil
		time.Sleep(time.Duration(1) * time.Second)
		_, err = NewConn.Write([]byte("km.version()\r"))
		if err != nil {
			_ = NewConn.Close()
			return nil, fmt.Errorf("ChangeBaudRate: write error after reconnect: %w", err)
		}
		ReadBuf := make([]byte, 2048)
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
		logger.Infof("Successfully Changed Baud Rate To %d!\n", 4000000)
		return NewConn, nil
	}

	//=====================================================================================
	// NewConn.Write([]byte("km.lock_mx(1)\r"))
	// NewConn.Read(ReadBuf)
	// logger.Infof("km.lock_mx(1) response: [%v]", string(ReadBuf[:n]))

	// time.Sleep(3 * time.Second)

	// NewConn.Write([]byte("km.lock_mx(0)\r"))
	// NewConn.Read(ReadBuf)
	// logger.Infof("km.lock_mx(0) response: [%v]", string(ReadBuf[:n]))

	// 西巴，这玩意不能完全独占输入
	// 只能锁定按键然后捕获按键
	// 锁定轴后，不能捕获轴的移动啊！！！

}

func getMackcuInstance(portName string, baudRate int) (*makcu, error) {
	port, err := getMakcuPort(portName, baudRate)
	if err != nil {
		return nil, err
	}
	readerChan := make(chan []byte)
	writerChan := make(chan []byte)
	go (func() {
		for {
			select {
			case <-globalCloseSignal:
				return
			case data := <-writerChan:
				_, err := port.Write(data)
				if err != nil {
					logger.Errorf("MAKCU Port write error: %v", err)
				}
			}

		}
	})()
	go (func() {
		ReadBuf := make([]byte, 2048)
		for {
			n, err := port.Read(ReadBuf)
			if err != nil {
				logger.Errorf("MAKCU Port read error: %v", err)
			}
			select {
			case <-globalCloseSignal:
				return
			case readerChan <- ReadBuf[:n]:
			default: // 缓冲满，丢弃
			}
		}
	})()

	return &makcu{
		port:   port,
		writer: writerChan,
		reader: readerChan,
	}, nil
}

func initInputAdapter_MakcuCOM(mk mouseKeyboard, makcu makcu) {
	for {
		select {
		case <-globalCloseSignal:
			return
			// case data := <-makcu.reader:
			// 	logger.Infof("MAKCU Port read: %v", string(data))
		}
	}
}
