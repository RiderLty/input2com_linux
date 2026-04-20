package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type AimBotState uint32

const (
	Off               AimBotState = 0xffffffff //关闭
	Sticky            AimBotState = 0x00       //吸附
	StickyAndFullAuto AimBotState = 0x01       //吸附+全自检测并动开火
	StickyFireOnce    AimBotState = 0x02       //开火时仅触发一次
	AutoTrigger       AimBotState = 0x03       //自动扳机 适用于狙击枪 目标距离满足条件则按下一次
	AutoTriggerBow    AimBotState = 0x04       //自动扳机 适用于弓 目标距离满足条件则按松开并按下一次
)

var AimBotWorkingState AimBotState = Sticky

func initInputAdapter_AimBot(port int, aimBotResult *AimBotResult, aimBotNotify chan bool) {
	conn, err := CreateUDPReader(fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.Errorf("创建UDP监听接口失败: %v", err)
		return
	} else {
		interfaces, err := net.Interfaces()
		if err != nil {
			panic(err)
		}
		logger.Info("正在从以下地址接收识别结果:")
		for _, iface := range interfaces {
			if iface.Flags&net.FlagUp == 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}
				ipv4 := ipNet.IP.To4()
				if ipv4 == nil {
					continue // 跳过非 IPv4 地址
				}
				if !ipv4.IsLoopback() {
					logger.Infof("UDP://%s:%v", ipv4, port)
				}
			}
		}
	}
	eventsCh := conn.ReadChan()

	for {
		select {
		case <-globalCloseSignal:
			return
		case pack := <-eventsCh:
			aimBotResult.width = int32(binary.LittleEndian.Uint32(pack[0 : 0+4]))
			aimBotResult.height = int32(binary.LittleEndian.Uint32(pack[4 : 4+4]))
			aimBotResult.offsetX = int32(binary.LittleEndian.Uint32(pack[8 : 8+4]))
			aimBotResult.offsetY = int32(binary.LittleEndian.Uint32(pack[12 : 12+4]))
			select {
			case aimBotNotify <- true:
			default:
			}
		}
	}
}
