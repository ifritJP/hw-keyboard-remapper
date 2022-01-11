package main

type KeyEvent struct {
	Code    uint8
	Pressed bool
	Name    string
}

func (event *KeyEvent) KeyString() string {
	return event.Name
}

func (event *KeyEvent) KeyPress() bool {
	return event.Pressed
}

func (event *KeyEvent) KeyRelease() bool {
	return !event.Pressed
}
