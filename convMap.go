// -*- coding:utf-8; -*-

package main

import "github.com/sirupsen/logrus"

type Code2HidCode struct {
	// linux のキーコード → HID のキーコード
	code2HidCode map[uint8]uint8
	// HID の remap コード
	remapHIDCode map[uint8]uint8
	// 処理を終了させるキーキーケンス
	exitKeySequence []uint8
	// 現在の キーシーケンスの位置
	keySequencePos     int
	exitKeySequenceTxt string
}

func NewCode2HidCode(exitKeySequenceTxt string) *Code2HidCode {
	code := Code2HidCode{}
	code.exitKeySequenceTxt = exitKeySequenceTxt
	code.exitKeySequence = make([]uint8, len(exitKeySequenceTxt))
	for index, oneChar := range exitKeySequenceTxt {
		code.exitKeySequence[index] = uint8(oneChar) - uint8('a') + uint8(0x04)
	}

	code.remapHIDCode = map[uint8]uint8{}
	code.code2HidCode = map[uint8]uint8{
		1:  0x29, // "ESC"
		2:  0x1E, // "1"
		3:  0x1F, // "2"
		4:  0x20, // "3"
		5:  0x21, // "4"
		6:  0x22, // "5"
		7:  0x23, // "6"
		8:  0x24, // "7"
		9:  0x25, // "8"
		10: 0x26, // "9"
		11: 0x27, // "0"
		12: 0x2D, // "-"
		13: 0x2E, // "="
		14: 0x2A, // "BS"
		15: 0x2B, // "TAB"
		16: 0x14, // "Q"
		17: 0x1A, // "W"
		18: 0x08, // "E"
		19: 0x15, // "R"
		20: 0x17, // "T"
		21: 0x1C, // "Y"
		22: 0x18, // "U"
		23: 0x0C, // "I"
		24: 0x12, // "O"
		25: 0x13, // "P"
		26: 0x2F, // "["
		27: 0x30, // "]"
		28: 0x28, // "ENTER"
		29: 0xE0, // "L_CTRL"
		30: 0x04, // "A"
		31: 0x16, // "S"
		32: 0x07, // "D"
		33: 0x09, // "F"
		34: 0x0A, // "G"
		35: 0x0B, // "H"
		36: 0x0D, // "J"
		37: 0x0E, // "K"
		38: 0x0F, // "L"
		39: 0x33, // ";"
		40: 0x34, // "'"
		41: 0x35, // "`"
		42: 0xE1, // "L_SHIFT"
		43: 0x31, // "\\"
		44: 0x1D, // "Z"
		45: 0x1B, // "X"
		46: 0x06, // "C"
		47: 0x19, // "V"
		48: 0x05, // "B"
		49: 0x11, // "N"
		50: 0x10, // "M"
		51: 0x36, // ","
		52: 0x37, // "."
		53: 0x38, // "/"
		54: 0xe5, // "R_SHIFT"
		55: 0x55, // "*"
		56: 0xe2, // "L_ALT"
		57: 0x2C, // "SPACE"
		58: 0x39, // "CAPS_LOCK"
		59: 0x3A, // "F1"
		60: 0x3B, // "F2"
		61: 0x3C, // "F3"
		62: 0x3D, // "F4"
		63: 0x3E, // "F5"
		64: 0x3F, // "F6"
		65: 0x40, // "F7"
		66: 0x41, // "F8"
		67: 0x42, // "F9"
		68: 0x43, // "F10"
		69: 0x83, // "NUM_LOCK"
		70: 0x84, // "SCROLL_LOCK"
		71: 0x5F, // "keypad HOME_7"
		72: 0x60, // "keypad UP_8"
		73: 0x61, // "keypad PGUP_9"
		74: 0x56, // "keypad -"
		75: 0x5C, // "keypad LEFT_4"
		76: 0x5D, // "keypad 5"
		77: 0x5E, // "keypad RT_ARROW_6"
		78: 0x57, // "keypad +"
		79: 0x59, // "keypad END_1"
		80: 0x5A, // "keypad DOWN_2"
		81: 0x5B, // "keypad PGDN_3"
		82: 0x62, // "keypad INS_0"
		83: 0x63, // "keypad DEL_."
		//	84: , // ""
		//	85: , // "全角半角"
		//	86: , // ""
		87: 0x44, // "F11"
		88: 0x45, // "F12"
		//	89: , // ""
		//	90: , // "カタナカ"
		//	91: , // "ひらがな"
		92: 0x8A, // "変換"
		93: 0x88, // "カタカナひらがな"
		94: 0x8B, // "無変換"
		//	95: , // "keypad comma"
		96:  0x58, // "keypad R_ENTER"
		97:  0xE4, // "R_CTRL"
		98:  0x54, // "/"
		99:  0x46, // "PRT_SCR"
		100: 0xE6, // "R_ALT"
		//	101:, // ""
		102: 0x4A, // "Home"
		103: 0x52, // "Up"
		104: 0x4B, // "PgUp"
		105: 0x50, // "Left"
		106: 0x4F, // "Right"
		107: 0x4D, // "End"
		108: 0x51, // "Down"
		109: 0x4E, // "PgDn"
		110: 0x49, // "Insert"
		111: 0x4C, // "Del"
		//	112: , // "macro"
		//	113: , // "mute"
		//	114: , // "volue down"
		//	115: , // "volue up"
		//	116: , // "power"
		//	117: , // "keypad ="
		//	118: , // "keypad +-"
		119: 0x48, // "Pause"
		124: 0x89, // "YEN"
		125: 0xE3, // win
	}
	return &code
}

func (conv *Code2HidCode) GetExitKeySequenceTxt() string {
	return conv.exitKeySequenceTxt
}

// HID コードの remap を設定
func (conv *Code2HidCode) SetHIDRemap(oldCode uint8, newCode uint8) {
	conv.remapHIDCode[oldCode] = newCode
}

func (conv *Code2HidCode) GetHIDKeyCode(code uint8) uint8 {
	// linux のコードから HID のコードに置き換える
	hidCode, has := conv.code2HidCode[code]
	if !has {
		return 0
	}
	// HID から、 remap 用 HID コードに置き換え
	convCode, has := conv.remapHIDCode[hidCode]
	if has {
		return convCode
	}
	return hidCode
}

func (conv *Code2HidCode) ProcessKeyEvent(keyboard *HIDKeyboard, keyEvent KeyEvent) ([]byte, bool, int) {
	hidCode := conv.GetHIDKeyCode(keyEvent.Code)
	hidKeyInfo := keyboard.GetKeyInfo(hidCode)

	eventTxt := ""
	// if the state of key is pressed
	if keyEvent.KeyPress() {
		keyboard.PressKey(hidCode)
		eventTxt = "press"

		// キーシーケンスのチェック
		if conv.keySequencePos < len(conv.exitKeySequence) {
			if conv.exitKeySequence[conv.keySequencePos] == hidCode {
				conv.keySequencePos++
			} else {
				if conv.exitKeySequence[0] == hidCode {
					conv.keySequencePos = 1
				} else {
					conv.keySequencePos = 0
				}
			}
		}
	}

	// if the state of key is released
	if keyEvent.KeyRelease() {
		keyboard.ReleaseKey(hidCode)
		eventTxt = "release"
	}

	logrus.Debugf(
		"[event] %s key %d(0x%x) %v -> %v",
		eventTxt, keyEvent.Code, keyEvent.Code, keyEvent.KeyString(), hidKeyInfo.Name)

	return keyboard.SetupHidPackat(), len(conv.exitKeySequence) == conv.keySequencePos, conv.keySequencePos
}
