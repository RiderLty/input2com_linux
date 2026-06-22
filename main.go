package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

var globalCloseSignal = make(chan bool) //仅会在程序退出时关闭  不用于其他用途

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
	UsingMouseDst    string `mapstructure:"usingMouseDst"`
	UsingKeyboardDst string `mapstructure:"usingKeyboardDst"`
}

var Cfg *Config

func main() {
	usingName, err := os.ReadFile("./config/using.txt")
	if err != nil {
		panic(fmt.Errorf("读取 using.txt 失败: %w", err))
	}
	profileName := strings.TrimSpace(string(usingName))
	viper.SetConfigFile(filepath.Join("./config/profiles", profileName+".yaml"))
	err = viper.ReadInConfig()
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
	var makcuInst *makcu
	var macroKB *macroMouseKeyboard

	// makcu 实例缓存，鼠标和键盘可共享同一个串口
	getOrCreateMakcu := func(ttyPath string, baudrate int) *makcu {
		if makcuInst != nil {
			return makcuInst
		}
		var err error
		makcuInst, err = getMackcuInstance(ttyPath, baudrate)
		if err != nil {
			logger.Errorf("ERROR:%v", err)
			os.Exit(1)
		}
		return makcuInst
	}

	// 创建输出控制器的工厂函数
	createController := func(dstName string) mouseKeyboard {
		switch dstName {
		case "kcom5":
			return NewMouseKeyboard_KCOM5(Cfg.Dst.Kcom5.TtyPath, Cfg.Dst.Kcom5.Baudrate, Cfg.Dst.Kcom5.Sbdesc, Cfg.Dst.Kcom5.Csdesc, Cfg.Dst.Kcom5.Cpdesc, Cfg.Dst.Kcom5.Xldesc)
		case "makcu":
			return NewMouseKeyboard_MAKCU(*getOrCreateMakcu(Cfg.Dst.Makcu.TtyPath, Cfg.Dst.Makcu.Baudrate))
		case "udp":
			return NewMouseKeyboard_UDP(fmt.Sprintf("%s:%d", Cfg.Dst.UDP.IP, Cfg.Dst.UDP.Port))
		case "uds":
			return NewMouseKeyboard_UDS(Cfg.Dst.UDS.Address)
		case "uinput":
			return NewMouseKeyboard_Uinput(Cfg.Dst.Uinput.DeviceName)
		case "usbgadget":
			return NewMouseKeyboard_USBGadget(Cfg.Dst.UsbGadget.MouseFile, Cfg.Dst.UsbGadget.KeyboardFile)
		default:
			logger.Fatalf("输出接口 %s 不支持", dstName)
			os.Exit(1)
			return nil
		}
	}

	logger.Infof("鼠标输出接口: %s, 键盘输出接口: %s", Cfg.UsingMouseDst, Cfg.UsingKeyboardDst)
	mouseCtrl := createController(Cfg.UsingMouseDst)
	keyboardCtrl := createController(Cfg.UsingKeyboardDst)
	macroKB = NewMouseKeyboard_MacroInterceptor(mouseCtrl, keyboardCtrl)

	if Cfg.Src.Inputs.Enabled {
		go initInputAdapter_LinuxInputs_MouseKeyboard(macroKB, Cfg.Src.Inputs.HotPlug, Cfg.Src.Inputs.Pattern)
		go initInputAdapter_LinuxInputs_Joystick(macroKB, Cfg.Src.Inputs.HotPlug, Cfg.Src.Inputs.Pattern)
		go initInputAdapter_LinuxInputs_Touchpad(macroKB, Cfg.Src.Inputs.HotPlug, Cfg.Src.Inputs.Pattern)
	}
	if Cfg.Src.Makcu.Enabled {
		getOrCreateMakcu(Cfg.Src.Makcu.TtyPath, Cfg.Src.Makcu.Baudrate)
		// 校验 src makcu 和已创建的 dst makcu 实例参数一致
		if Cfg.Src.Makcu.TtyPath != Cfg.Dst.Makcu.TtyPath || Cfg.Src.Makcu.Baudrate != Cfg.Dst.Makcu.Baudrate {
			logger.Errorf("src makcu ttyPath %s baudrate %d 与 dst makcu ttyPath %s baudrate %d 不一致", Cfg.Src.Makcu.TtyPath, Cfg.Src.Makcu.Baudrate, Cfg.Dst.Makcu.TtyPath, Cfg.Dst.Makcu.Baudrate)
			os.Exit(1)
		}
		go initInputAdapter_MakcuCOM(macroKB, *makcuInst)
	}
	if Cfg.Src.UDP.Enabled {
		go initInputAdapter_UDP(macroKB, Cfg.Src.UDP.Port)
	}
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-exitChan
	close(globalCloseSignal)
	logger.Info("已停止")
	time.Sleep(time.Millisecond * 40)
}
