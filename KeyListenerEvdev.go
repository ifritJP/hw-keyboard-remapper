// +build linux,!logger

package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	"github.com/sirupsen/logrus"
)

const (
	device_glob = "/dev/input/event*"
)

func ListInputDevices() ([]string, error) {
	devices, err := evdev.ListInputDevices(device_glob)
	if err != nil {
		return nil, err
	}
	list := make([]string, len(devices))
	for index, dev := range devices {
		list[index] = dev.Name
	}
	return list, nil
}

func select_device(devName string) (*evdev.InputDevice, error) {
	devices, _ := evdev.ListInputDevices(device_glob)

	lines := make([]string, 0)
	max := 0
	if len(devices) > 0 {
		for i := range devices {
			dev := devices[i]
			str := fmt.Sprintf("%-3d %-20s %-35s %s", i, dev.Fn, dev.Name, dev.Phys)
			if len(str) > max {
				max = len(str)
			}
			lines = append(lines, str)
			if devName != "" && dev.Name == devName {
				return dev, nil
			}
		}
		if devName != "" {
			return nil, nil
		}

		fmt.Printf("%-3s %-20s %-35s %s\n", "ID", "Device", "Name", "Phys")
		fmt.Printf(strings.Repeat("-", max) + "\n")
		fmt.Printf(strings.Join(lines, "\n") + "\n")

		var choice int
		choice_max := len(lines) - 1

		for {
			fmt.Printf("Select device [0-%d]: ", choice_max)
			_, err := fmt.Scan(&choice)
			if err != nil {
				return nil, err
			}
			if choice <= choice_max && choice >= 0 {
				return devices[choice], nil
			}
		}
	}

	errmsg := fmt.Sprintf("no accessible input devices found by %s", device_glob)
	return nil, errors.New(errmsg)
}

func format_event(ev *evdev.InputEvent) (KeyEvent, bool) {
	var code_name string

	code := int(ev.Code)

	switch ev.Type {
	case evdev.EV_KEY:
		val, haskey := evdev.KEY[code]
		if haskey {
			code_name = val
		} else {
			val, haskey := evdev.BTN[code]
			if haskey {
				code_name = val
			} else {
				code_name = "?"
			}
		}

		keyEvent := KeyEvent{
			Code: uint8(code), Pressed: ev.Value > 0, Name: code_name}
		logrus.Tracef("KeyEvent = %v", ev)
		return keyEvent, true
	}
	return KeyEvent{}, false
}

func SetKeyListener(
	keyboardName string, listener func(keyEvent KeyEvent)) error {
	var dev *evdev.InputDevice
	var events []evdev.InputEvent
	var err error
	for {
		dev, err = select_device(keyboardName)
		if err != nil {
			return err
		}
		if dev != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	logrus.Infof("ready keyboardName = %s", dev.Name)

	dev.Grab()
	defer dev.Release()

	for {
		events, err = dev.Read()
		if err != nil {
			return err
		}
		for i := range events {
			keyEvent, ok := format_event(&events[i])
			if ok {
				listener(keyEvent)
			}
		}
	}
}
