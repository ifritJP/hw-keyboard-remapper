package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strconv"

	"github.com/sirupsen/logrus"
)

func setSignal(callback func()) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		callback()
		fmt.Printf("\nsignal\n")
		os.Exit(1)
	}()
}

func main() {

	var cmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	help := cmd.Bool("help", false, "display help message")
	cmd.Usage = func() {
		fmt.Fprintf(cmd.Output(), "\nUsage: %s options\n\n", os.Args[0])
		fmt.Fprintf(cmd.Output(), " options:\n\n")
		cmd.PrintDefaults()
		os.Exit(1)
	}

	verboseMode := cmd.Bool("v", false, "verbose")
	configPath := cmd.String("conf", "", "config file path")
	keyboardOp := cmd.String("kb", "", "keyboard name")
	logLevel := cmd.Int(
		"log", int(logrus.DebugLevel),
		fmt.Sprintf("log level %d - %d", logrus.FatalLevel, logrus.TraceLevel))

	opMode := cmd.String("mode", "remap", "operation mode. [remap,list,scan]")

	if len(os.Args) <= 1 {
		cmd.Usage()
	}
	cmd.Parse(os.Args[1:])
	if *help {
		cmd.Usage()
	}
	if *opMode == "list" {
		if list, err := ListInputDevices(); err != nil {
			logrus.Error(err)
		} else {
			if len(list) == 0 {
				fmt.Printf("It maybe need to use sudo command.\n")
			} else {
				for _, name := range list {
					fmt.Printf("'%s'\n", name)
				}
			}
			os.Exit(1)
		}
	}

	if *verboseMode {
		if logLevel != nil {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.Level(*logLevel))
		}
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	convCode := NewCode2HidCode("qweqweqweqwe")
	hidKeyboard := NewHIDKeyboard()

	keyboardName := ""
	logrus.Infof("configPath = %v", configPath)
	if *configPath != "" {
		if setting, err := load(*configPath); err != nil {
			logrus.Error(err)
			os.Exit(1)
		} else {
			logrus.Infof("config.json = %v", setting)
			if setting.InputKeyboardName != nil {
				keyboardName = *setting.InputKeyboardName
			}
			for _, switchKey := range setting.SwitchKeys {
				if switchKey.On == nil || *switchKey.On {
					convCode.SetHIDRemap(switchKey.Src, switchKey.Dst)
				}
			}
			for codeTxt, convKeyList := range setting.ConvKeyMap {
				for _, convKey := range convKeyList {
					if code, err := strconv.ParseUint(codeTxt, 0, 8); err != nil {
						logrus.Error(err)
					} else {
						if convKey.On == nil || *convKey.On {
							// AddConvKey() する ConvKeyInfo 情報のオブジェクトを
							// 別々にするため、 cloneConvKey を作る。
							cloneConvKey := convKey
							hidKeyboard.AddConvKey(byte(code), &cloneConvKey)
						}
					}
				}
			}
		}
	}
	if *keyboardOp != "" {
		keyboardName = *keyboardOp
	}

	if *opMode == "scan" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Infof("Detecting keyboard = %s", keyboardName)
		logrus.Infof(
			"Enter '%s', if you want to exit from this program.",
			convCode.GetExitKeySequenceTxt())
		SetKeyListener(keyboardName, func(keyEvent KeyEvent) {
			data, _, _ := convCode.ProcessKeyEvent(hidKeyboard, keyEvent)
			logrus.Printf("data %v", data)
		})
		os.Exit(0)
	}

	if keyboardName == "" {
		fmt.Printf("keyboard isn't set. Please set -kb option or set config.\n")
		os.Exit(1)
	}

	hidOut, err := os.OpenFile("/dev/hidg0", os.O_RDWR, os.ModeCharDevice)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	logrus.Infof("keyboardName = %s", keyboardName)
	setSignal(func() {
		// 強制停止の時に、変な data を送信したままにしないように
		// 全 0 のデータでクリアする
		zeroData := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		hidOut.Write(zeroData)
		hidOut.Write(zeroData)
	})
	for {
		hidKeyboard.ReleaseAllKeys()
		logrus.Infof("Detecting keyboard = %s", keyboardName)
		logrus.Infof(
			"Enter '%s', if you want to exit from this program.",
			convCode.GetExitKeySequenceTxt())
		SetKeyListener(keyboardName, func(keyEvent KeyEvent) {
			data, matchkeySeq, keySeqPos :=
				convCode.ProcessKeyEvent(hidKeyboard, keyEvent)
			logrus.Debugf("data %v, %d", data, keySeqPos)
			if matchkeySeq {
				logrus.Printf("match key sequence")
				zeroData := []byte{0, 0, 0, 0, 0, 0, 0, 0}
				hidOut.Write(zeroData)
				os.Exit(0)
			}
			hidOut.Write(data)
		})
		time.Sleep(1 * time.Second)
	}
}
