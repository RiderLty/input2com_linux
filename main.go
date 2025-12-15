package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
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

type Config struct {
	Debug  bool `mapstructure:"debug"`
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Src struct {
		Inputs struct {
			Enabled bool   `mapstructure:"enabled"`
			HotPlug bool   `mapstructure:"hotPlug"`
			Pattern string `mapstructure:"pattern"`
		} `mapstructure:"inputs"`
		Makcu struct {
			Enabled  bool   `mapstructure:"enabled"`
			Baudrate int    `mapstructure:"baudrate"`
			TtyPath  string `mapstructure:"ttyPath"`
		} `mapstructure:"makcu"`
		UDP struct {
			Enabled bool `mapstructure:"enabled"`
			Port    int  `mapstructure:"port"`
		} `mapstructure:"udp"`
	} `mapstructure:"src"`
	Dst struct {
		Kcom5 struct {
			Baudrate int    `mapstructure:"baudrate"`
			TtyPath  string `mapstructure:"ttyPath"`
			Sbdesc   string `mapstructure:"sbdesc"`
			Csdesc   string `mapstructure:"csdesc"`
			Cpdesc   string `mapstructure:"cpdesc"`
			Xldesc   string `mapstructure:"xldesc"`
		} `mapstructure:"kcom5"`
		Makcu struct {
			Baudrate int    `mapstructure:"baudrate"`
			TtyPath  string `mapstructure:"ttyPath"`
		} `mapstructure:"makcu"`
		UDP struct {
			IP   string `mapstructure:"ip"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"udp"`
		UDS struct {
			Address string `mapstructure:"address"`
		} `mapstructure:"uds"`
		Uinput struct {
			DeviceName string `mapstructure:"deviceName"`
		} `mapstructure:"uinput"`
		UsbGadget struct {
			MouseFile    string `mapstructure:"mouseFile"`
			KeyboardFile string `mapstructure:"keyboardFile"`
		} `mapstructure:"usbgadget"`
	} `mapstructure:"dst"`
	UsingDst string `mapstructure:"usingDst"`
}

var Cfg *Config

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
	go serve(Cfg.Server.Port) //启动配置服务器
	if Cfg.Debug {
		logger.WithDebug()
	}
	var makcu *makcu
	var makcu_err error
	var macroKB *macroMouseKeyboard

	logger.Infof("使用输出接口 %s", Cfg.UsingDst)
	switch Cfg.UsingDst {

	case "kcom5":
		macroKB = NewMouseKeyboard_MacroInterceptor(
			NewMouseKeyboard_KCOM5(Cfg.Dst.Kcom5.TtyPath, Cfg.Dst.Kcom5.Baudrate, Cfg.Dst.Kcom5.Sbdesc, Cfg.Dst.Kcom5.Csdesc, Cfg.Dst.Kcom5.Cpdesc, Cfg.Dst.Kcom5.Xldesc),
		)
	case "makcu":
		makcu, makcu_err = getMackcuInstance(Cfg.Dst.Makcu.TtyPath, Cfg.Dst.Makcu.Baudrate)
		if makcu_err != nil {
			logger.Errorf("ERROR:%v", makcu_err)
			os.Exit(1)
		}
		macroKB = NewMouseKeyboard_MacroInterceptor(
			NewMouseKeyboard_MAKCU(*makcu),
		)
	case "udp":
		macroKB = NewMouseKeyboard_MacroInterceptor(
			NewMouseKeyboard_UDP(fmt.Sprintf("%s:%d", Cfg.Dst.UDP.IP, Cfg.Dst.UDP.Port)),
		)
	case "uds":
		macroKB = NewMouseKeyboard_MacroInterceptor(
			NewMouseKeyboard_UDS(Cfg.Dst.UDS.Address),
		)
	case "uinput":
		macroKB = NewMouseKeyboard_MacroInterceptor(
			NewMouseKeyboard_Uinput(Cfg.Dst.Uinput.DeviceName),
		)
	case "usbgadget":
		macroKB = NewMouseKeyboard_MacroInterceptor(
			NewMouseKeyboard_USBGadget(Cfg.Dst.UsbGadget.MouseFile, Cfg.Dst.UsbGadget.KeyboardFile),
		)

	default:
		logger.Fatalf("usingDst %s not support", Cfg.UsingDst)
		os.Exit(1)
	}

	if Cfg.Src.Inputs.Enabled {
		go initInputAdapter_LinuxInputs(macroKB, Cfg.Src.Inputs.HotPlug, Cfg.Src.Inputs.Pattern)
	}
	if Cfg.Src.Makcu.Enabled {
		if makcu == nil { //尝试复用makcu实例
			makcu, makcu_err = getMackcuInstance(Cfg.Src.Makcu.TtyPath, Cfg.Src.Makcu.Baudrate)
			if makcu_err != nil {
				logger.Errorf("ERROR:%v", makcu_err)
				os.Exit(1)
			}
		} else {
			if Cfg.Src.Makcu.TtyPath != Cfg.Dst.Makcu.TtyPath || Cfg.Src.Makcu.Baudrate != Cfg.Dst.Makcu.Baudrate {
				logger.Errorf("src makcu ttyPath %s baudrate %d not equal dst makcu ttyPath %s baudrate %d", Cfg.Src.Makcu.TtyPath, Cfg.Src.Makcu.Baudrate, Cfg.Dst.Makcu.TtyPath, Cfg.Dst.Makcu.Baudrate)
				os.Exit(1)
			}
		}
		go initInputAdapter_MakcuCOM(macroKB, *makcu)
	}
	if Cfg.Src.UDP.Enabled {
		go initInputAdapter_UDP(macroKB, Cfg.Src.UDP.Port)
	}
	go udp_listener(9321)
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-exitChan
	close(globalCloseSignal)
	logger.Info("已停止")
	time.Sleep(time.Millisecond * 40)
}
