package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
)

var globalCloseSignal = make(chan bool) //仅会在程序退出时关闭  不用于其他用途

var udp_ints [6]int32
var udp_last byte

func udp_listener(port int) {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	if err != nil {
		logger.Errorf("udp error : %v", err)
		return
	}
	defer listen.Close()

	recv_ch := make(chan []byte)
	go func() {
		for {
			var buf [1024]byte
			n, _, err := listen.ReadFromUDP(buf[:])
			if err != nil {
				break
			}
			recv_ch <- buf[:n]
		}
	}()
	logger.Infof("已准备接收远程事件: 0.0.0.0:%d", port)
	for {
		select {
		case <-globalCloseSignal:
			return
		case pack := <-recv_ch:

			for i := 0; i < 6; i++ {
				start := i * 4
				udp_ints[i] = int32(binary.LittleEndian.Uint32(pack[start : start+4]))
			}
			udp_last = pack[24]
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

	var auto_detect = parser.Flag("a", "auto-detect", &argparse.Options{
		Required: false,
		Default:  false,
		Help:     "关闭自动检测设备 默认开启",
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

	var patern = parser.String("", "pattern", &argparse.Options{
		Required: false,
		Default:  ".*",
		Help:     "捕获设备名称的通配符模式",
	})

	var port = parser.Int("p", "port", &argparse.Options{
		Required: false,
		Help:     "端口",
		Default:  9264,
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	go serve(*port) //启动配置服务器

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
	logger.Infof("%v,%v,%v,%v,%v,%v", devpath, *badurate, *sbDesc, *csDesc, *cpDesc, *xlDesc)

	makcu, err := getMackcuInstance("/dev/ttyACM0", 4000000)
	if err != nil {
		logger.Errorf("ERROR:%v", err)
		return
	}

	makcuKB := NewMouseKeyboard_MAKCU(*makcu)
	// comKB := NewMouseKeyboard_KCOM5(devpath, *badurate, *sbDesc, *csDesc, *cpDesc, *xlDesc)
	macroKB := NewMouseKeyboard_MacroInterceptor(makcuKB)

	go initInputAdapter_LinuxInputs(macroKB, *auto_detect, *patern)

	go initInputAdapter_MakcuCOM(macroKB, *makcu)
	// comKB := NewMouseKeyboard_KCOM5(devpath, *badurate, *sbDesc, *csDesc, *cpDesc, *xlDesc)
	// macroKB := NewMouseKeyboard_MacroInterceptor(comKB)
	// makcuKB, err := Connect(devpath, 4000000)
	// if makcuKB == nil {
	// 	panic(err)
	// }
	// makcuKB.SetButtonStatus(true)
	// makcuKB, err = ChangeBaudRate(makcuKB)
	// if makcuKB == nil {
	// 	panic(err)
	// }
	// go makcuKB.ListenLoop()
	// defer makcuKB.Close()

	// handelMakcuEvent := func(btn MouseButton, pressed bool) {
	// 	logger.Logger.Infof("btn %d, pressed: %t", btn, pressed)
	// 	// if pressed {
	// 	// 	macroKB.BtnDown(byte(1<<btn), "makcu")
	// 	// } else {
	// 	// 	macroKB.BtnUp(byte(1<<btn), "makcu")
	// 	// }
	// }
	// makcuKB.SetButtonCallback(handelMakcuEvent)

	go udp_listener(9321)

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-exitChan
	close(globalCloseSignal)
	logger.Info("已停止")
	time.Sleep(time.Millisecond * 40)
}
