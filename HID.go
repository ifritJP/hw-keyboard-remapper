// -*- coding:utf-8; -*-

package main

const L_SHIFTBIT = uint8(1 << 1)
const LR_SHIFTBIT = uint8(L_SHIFTBIT | (1 << 5))

type ConvKeyInfo struct {
	On *bool
	// HID の modifier の一致条件。
	// HID の modifier と以下の式が成立する時に、このキーに置き換える。
	// (modifier & condModifierMask) == condModifierResult
	CondModifierMask   byte `json:"modMask"`
	CondModifierResult byte `json:"modResult"`
	// HID コード
	Code byte
	// modifier に XOR する値
	ModifierXor byte `json:"modXor"`
}

// HID のキーの状態
type HIDKeyInfo struct {
	// HID キーコード
	OrgCode byte
	// キーコード名
	Name string
	// modifier キーかどうか
	IsModifier bool
	// 押されているかどうか
	Pressed bool
	// ConvKeyInfo
	convKeyInfoList []*ConvKeyInfo
}

// HID の modifier bit を返す
func (info *HIDKeyInfo) GetModifierBit() uint8 {
	if info.Pressed && info.IsModifier {
		return 1 << (info.OrgCode - 0xe0)
	}
	return 0
}

// 置き換えを処理する。
//
// @param modifierFlag 置き換え前の modifierFlag
// @return byte 置き換え後の HID キーコード
// @return byte 置き換え後の modifierFlag
func (info *HIDKeyInfo) process(modifierFlag byte) (byte, byte) {
	for _, convKey := range info.convKeyInfoList {
		// 置き換え情報を処理する
		if (modifierFlag & convKey.CondModifierMask) == convKey.CondModifierResult {
			modifierFlag = modifierFlag ^ convKey.ModifierXor
			return convKey.Code, modifierFlag
		}
	}
	if !info.IsModifier {
		return info.OrgCode, modifierFlag
	}
	return 0, modifierFlag
}

type HIDKeyboard struct {
	// HID のキーパケット 8 バイト
	// byte0: modifier
	// byte1: reservede
	// byte2: pressed-key1
	// byte3: pressed-key2
	// byte4: pressed-key3
	// byte5: pressed-key4
	// byte6: pressed-key5
	// byte7: pressed-key6
	data []byte
	// HID キーコード → HIDKeyInfo
	keyInfoMap map[uint8]*HIDKeyInfo
}

func NewHIDKeyInfo(code byte, name string, modifier bool) *HIDKeyInfo {
	return &HIDKeyInfo{code, name, modifier, false, []*ConvKeyInfo{}}
}

func NewHIDKeyboard() *HIDKeyboard {
	keyInfoMap := map[uint8]*HIDKeyInfo{
		0x00: NewHIDKeyInfo(0x00, "Reserved", false),
		0x01: NewHIDKeyInfo(0x01, "Keyboard ErrorRollOver", false),
		0x02: NewHIDKeyInfo(0x02, "Keyboard POSTFail", false),
		0x03: NewHIDKeyInfo(0x03, "Keyboard ErrorUndefined", false),
		0x04: NewHIDKeyInfo(0x04, "Keyboard a and A", false),
		0x05: NewHIDKeyInfo(0x05, "Keyboard b and B", false),
		0x06: NewHIDKeyInfo(0x06, "Keyboard c and C", false),
		0x07: NewHIDKeyInfo(0x07, "Keyboard d and D", false),
		0x08: NewHIDKeyInfo(0x08, "Keyboard e and E", false),
		0x09: NewHIDKeyInfo(0x09, "Keyboard f and F", false),
		0x0A: NewHIDKeyInfo(0x0A, "Keyboard g and G", false),
		0x0B: NewHIDKeyInfo(0x0B, "Keyboard h and H", false),
		0x0C: NewHIDKeyInfo(0x0C, "Keyboard i and I", false),
		0x0D: NewHIDKeyInfo(0x0D, "Keyboard j and J", false),
		0x0E: NewHIDKeyInfo(0x0E, "Keyboard k and K", false),
		0x0F: NewHIDKeyInfo(0x0F, "Keyboard l and L", false),
		0x10: NewHIDKeyInfo(0x10, "Keyboard m and M", false),
		0x11: NewHIDKeyInfo(0x11, "Keyboard n and N", false),
		0x12: NewHIDKeyInfo(0x12, "Keyboard o and O", false),
		0x13: NewHIDKeyInfo(0x13, "Keyboard p and P", false),
		0x14: NewHIDKeyInfo(0x14, "Keyboard q and Q", false),
		0x15: NewHIDKeyInfo(0x15, "Keyboard r and R", false),
		0x16: NewHIDKeyInfo(0x16, "Keyboard s and S", false),
		0x17: NewHIDKeyInfo(0x17, "Keyboard t and T", false),
		0x18: NewHIDKeyInfo(0x18, "Keyboard u and U", false),
		0x19: NewHIDKeyInfo(0x19, "Keyboard v and V", false),
		0x1A: NewHIDKeyInfo(0x1A, "Keyboard w and W", false),
		0x1B: NewHIDKeyInfo(0x1B, "Keyboard x and X", false),
		0x1C: NewHIDKeyInfo(0x1C, "Keyboard y and Y", false),
		0x1D: NewHIDKeyInfo(0x1D, "Keyboard z and Z", false),
		0x1E: NewHIDKeyInfo(0x1E, "Keyboard 1 and !", false),
		0x1F: NewHIDKeyInfo(0x1F, "Keyboard 2 and @", false),
		0x20: NewHIDKeyInfo(0x20, "Keyboard 3 and #", false),
		0x21: NewHIDKeyInfo(0x21, "Keyboard 4 and $", false),
		0x22: NewHIDKeyInfo(0x22, "Keyboard 5 and %", false),
		0x23: NewHIDKeyInfo(0x23, "Keyboard 6 and ∧", false),
		0x24: NewHIDKeyInfo(0x24, "Keyboard 7 and &", false),
		0x25: NewHIDKeyInfo(0x25, "Keyboard 8 and *", false),
		0x26: NewHIDKeyInfo(0x26, "Keyboard 9 and (", false),
		0x27: NewHIDKeyInfo(0x27, "Keyboard 0 and )", false),
		0x28: NewHIDKeyInfo(0x28, "Keyboard Return (ENTER)", false),
		0x29: NewHIDKeyInfo(0x29, "Keyboard ESCAPE", false),
		0x2A: NewHIDKeyInfo(0x2A, "Keyboard DELETE (Backspace)", false),
		0x2B: NewHIDKeyInfo(0x2B, "Keyboard Tab", false),
		0x2C: NewHIDKeyInfo(0x2C, "Keyboard Spacebar", false),
		0x2D: NewHIDKeyInfo(0x2D, "Keyboard - and (underscore)", false),
		0x2E: NewHIDKeyInfo(0x2E, "Keyboard = and +", false),
		0x2F: NewHIDKeyInfo(0x2F, "Keyboard [ and {", false),
		0x30: NewHIDKeyInfo(0x30, "Keyboard ] and }", false),
		0x31: NewHIDKeyInfo(0x31, "Keyboard \\ and |", false),
		0x32: NewHIDKeyInfo(0x32, "Keyboard Non-US # and `", false),
		0x33: NewHIDKeyInfo(0x33, "Keyboard ; and :", false),
		0x34: NewHIDKeyInfo(0x34, "Keyboard ' and \"", false),
		0x35: NewHIDKeyInfo(0x35, "Keyboard Grave Accent and Tilde", false),
		0x36: NewHIDKeyInfo(0x36, "Keyboard , and <", false),
		0x37: NewHIDKeyInfo(0x37, "Keyboard . and >", false),
		0x38: NewHIDKeyInfo(0x38, "Keyboard / and ?", false),
		0x39: NewHIDKeyInfo(0x39, "Keyboard Caps Lock", false),
		0x3A: NewHIDKeyInfo(0x3A, "Keyboard F1", false),
		0x3B: NewHIDKeyInfo(0x3B, "Keyboard F2", false),
		0x3C: NewHIDKeyInfo(0x3C, "Keyboard F3", false),
		0x3D: NewHIDKeyInfo(0x3D, "Keyboard F4", false),
		0x3E: NewHIDKeyInfo(0x3E, "Keyboard F5", false),
		0x3F: NewHIDKeyInfo(0x3F, "Keyboard F6", false),
		0x40: NewHIDKeyInfo(0x40, "Keyboard F7", false),
		0x41: NewHIDKeyInfo(0x41, "Keyboard F8", false),
		0x42: NewHIDKeyInfo(0x42, "Keyboard F9", false),
		0x43: NewHIDKeyInfo(0x43, "Keyboard F10", false),
		0x44: NewHIDKeyInfo(0x44, "Keyboard F11", false),
		0x45: NewHIDKeyInfo(0x45, "Keyboard F12", false),
		0x46: NewHIDKeyInfo(0x46, "Keyboard PrintScreen", false),
		0x47: NewHIDKeyInfo(0x47, "Keyboard Scroll Lock", false),
		0x48: NewHIDKeyInfo(0x48, "Keyboard Pause", false),
		0x49: NewHIDKeyInfo(0x49, "Keyboard Insert", false),
		0x4A: NewHIDKeyInfo(0x4A, "Keyboard Home", false),
		0x4B: NewHIDKeyInfo(0x4B, "Keyboard PageUp", false),
		0x4C: NewHIDKeyInfo(0x4C, "Keyboard Delete Forward", false),
		0x4D: NewHIDKeyInfo(0x4D, "Keyboard End", false),
		0x4E: NewHIDKeyInfo(0x4E, "Keyboard PageDown", false),
		0x4F: NewHIDKeyInfo(0x4F, "Keyboard RightArrow", false),
		0x50: NewHIDKeyInfo(0x50, "Keyboard LeftArrow", false),
		0x51: NewHIDKeyInfo(0x51, "Keyboard DownArrow", false),
		0x52: NewHIDKeyInfo(0x52, "Keyboard UpArrow", false),
		0x53: NewHIDKeyInfo(0x53, "Keypad Num Lock and Clear", false),
		0x54: NewHIDKeyInfo(0x54, "Keypad /", false),
		0x55: NewHIDKeyInfo(0x55, "Keypad *", false),
		0x56: NewHIDKeyInfo(0x56, "Keypad -", false),
		0x57: NewHIDKeyInfo(0x57, "Keypad +", false),
		0x58: NewHIDKeyInfo(0x58, "Keypad ENTER", false),
		0x59: NewHIDKeyInfo(0x59, "Keypad 1 and End", false),
		0x5A: NewHIDKeyInfo(0x5A, "Keypad 2 and Down Arrow", false),
		0x5B: NewHIDKeyInfo(0x5B, "Keypad 3 and PageDn", false),
		0x5C: NewHIDKeyInfo(0x5C, "Keypad 4 and Left Arrow", false),
		0x5D: NewHIDKeyInfo(0x5D, "Keypad 5", false),
		0x5E: NewHIDKeyInfo(0x5E, "Keypad 6 and Right Arrow", false),
		0x5F: NewHIDKeyInfo(0x5F, "Keypad 7 and Home", false),
		0x60: NewHIDKeyInfo(0x60, "Keypad 8 and Up Arrow", false),
		0x61: NewHIDKeyInfo(0x61, "Keypad 9 and PageUp", false),
		0x62: NewHIDKeyInfo(0x62, "Keypad 0 and Insert", false),
		0x63: NewHIDKeyInfo(0x63, "Keypad . and Delete", false),
		0x64: NewHIDKeyInfo(0x64, "Keyboard Non-US \\ and |", false),
		0x65: NewHIDKeyInfo(0x65, "Keyboard Application", false),
		0x66: NewHIDKeyInfo(0x66, "Keyboard Power", false),
		0x67: NewHIDKeyInfo(0x67, "Keypad =", false),
		0x68: NewHIDKeyInfo(0x68, "Keyboard F13", false),
		0x69: NewHIDKeyInfo(0x69, "Keyboard F14", false),
		0x6A: NewHIDKeyInfo(0x6A, "Keyboard F15", false),
		0x6B: NewHIDKeyInfo(0x6B, "Keyboard F16 ", false),
		0x6C: NewHIDKeyInfo(0x6C, "Keyboard F17 ", false),
		0x6D: NewHIDKeyInfo(0x6D, "Keyboard F18 ", false),
		0x6E: NewHIDKeyInfo(0x6E, "Keyboard F19 ", false),
		0x6F: NewHIDKeyInfo(0x6F, "Keyboard F20 ", false),
		0x70: NewHIDKeyInfo(0x70, "Keyboard F21 ", false),
		0x71: NewHIDKeyInfo(0x71, "Keyboard F22 ", false),
		0x72: NewHIDKeyInfo(0x72, "Keyboard F23 ", false),
		0x73: NewHIDKeyInfo(0x73, "Keyboard F24 ", false),
		0x74: NewHIDKeyInfo(0x74, "Keyboard Execute", false),
		0x75: NewHIDKeyInfo(0x75, "Keyboard Help", false),
		0x76: NewHIDKeyInfo(0x76, "Keyboard Menu", false),
		0x77: NewHIDKeyInfo(0x77, "Keyboard Select", false),
		0x78: NewHIDKeyInfo(0x78, "Keyboard Stop", false),
		0x79: NewHIDKeyInfo(0x79, "Keyboard Again", false),
		0x7A: NewHIDKeyInfo(0x7A, "Keyboard Undo", false),
		0x7B: NewHIDKeyInfo(0x7B, "Keyboard Cut", false),
		0x7C: NewHIDKeyInfo(0x7C, "Keyboard Copy", false),
		0x7D: NewHIDKeyInfo(0x7D, "Keyboard Paste", false),
		0x7E: NewHIDKeyInfo(0x7E, "Keyboard Find", false),
		0x7F: NewHIDKeyInfo(0x7F, "Keyboard Mute", false),
		0x80: NewHIDKeyInfo(0x80, "Keyboard Volume Up", false),
		0x81: NewHIDKeyInfo(0x81, "Keyboard Volume Down", false),
		0x82: NewHIDKeyInfo(0x82, "Keyboard Locking Caps Lock", false),
		0x83: NewHIDKeyInfo(0x83, "Keyboard Locking Num Lock", false),
		0x84: NewHIDKeyInfo(0x84, "Keyboard Locking Scroll Lock", false),
		0x85: NewHIDKeyInfo(0x85, "Keypad Comma", false),
		0x86: NewHIDKeyInfo(0x86, "Keypad Equal Sign", false),
		0x87: NewHIDKeyInfo(0x87, "Keyboard International1", false),
		0x88: NewHIDKeyInfo(0x88, "Keyboard International2 katakana-hiragana", false),
		0x89: NewHIDKeyInfo(0x89, "Keyboard International3 ", false),
		0x8A: NewHIDKeyInfo(0x8A, "Keyboard International4 henkan", false),
		0x8B: NewHIDKeyInfo(0x8B, "Keyboard International5 muhenkan", false),
		0x8C: NewHIDKeyInfo(0x8C, "Keyboard International6 ", false),
		0x8D: NewHIDKeyInfo(0x8D, "Keyboard International7 ", false),
		0x8E: NewHIDKeyInfo(0x8E, "Keyboard International8 ", false),
		0x8F: NewHIDKeyInfo(0x8F, "Keyboard International9 ", false),
		0x90: NewHIDKeyInfo(0x90, "Keyboard LANG1", false),
		0x91: NewHIDKeyInfo(0x91, "Keyboard LANG2", false),
		0x92: NewHIDKeyInfo(0x92, "Keyboard LANG3", false),
		0x93: NewHIDKeyInfo(0x93, "Keyboard LANG4", false),
		0x94: NewHIDKeyInfo(0x94, "Keyboard LANG5", false),
		0x95: NewHIDKeyInfo(0x95, "Keyboard LANG6", false),
		0x96: NewHIDKeyInfo(0x96, "Keyboard LANG7", false),
		0x97: NewHIDKeyInfo(0x97, "Keyboard LANG8", false),
		0x98: NewHIDKeyInfo(0x98, "Keyboard LANG9", false),
		0x99: NewHIDKeyInfo(0x99, "Keyboard Alternate Erase", false),
		0x9A: NewHIDKeyInfo(0x9A, "Keyboard SysReq/Attention", false),
		0x9B: NewHIDKeyInfo(0x9B, "Keyboard Cancel ", false),
		0x9C: NewHIDKeyInfo(0x9C, "Keyboard Clear ", false),
		0x9D: NewHIDKeyInfo(0x9D, "Keyboard Prior ", false),
		0x9E: NewHIDKeyInfo(0x9E, "Keyboard Return ", false),
		0x9F: NewHIDKeyInfo(0x9F, "Keyboard Separator ", false),
		0xA0: NewHIDKeyInfo(0xA0, "Keyboard Out ", false),
		0xA1: NewHIDKeyInfo(0xA1, "Keyboard Oper ", false),
		0xA2: NewHIDKeyInfo(0xA2, "Keyboard Clear/Again ", false),
		0xA3: NewHIDKeyInfo(0xA3, "Keyboard CrSel/Props ", false),
		0xA4: NewHIDKeyInfo(0xA4, "Keyboard ExSel ", false),
		// 0xA5-AF Reserved
		0xB0: NewHIDKeyInfo(0xB0, "Keypad 00 ", false),
		0xB1: NewHIDKeyInfo(0xB1, "Keypad 000 ", false),
		0xB2: NewHIDKeyInfo(0xB2, "Thousands Separator", false),
		0xB3: NewHIDKeyInfo(0xB3, "Decimal Separator", false),
		0xB4: NewHIDKeyInfo(0xB4, "Currency Unit", false),
		0xB5: NewHIDKeyInfo(0xB5, "Currency Sub-unit", false),
		0xB6: NewHIDKeyInfo(0xB6, "Keypad ( ", false),
		0xB7: NewHIDKeyInfo(0xB7, "Keypad ) ", false),
		0xB8: NewHIDKeyInfo(0xB8, "Keypad { ", false),
		0xB9: NewHIDKeyInfo(0xB9, "Keypad } ", false),
		0xBA: NewHIDKeyInfo(0xBA, "Keypad Tab ", false),
		0xBB: NewHIDKeyInfo(0xBB, "Keypad Backspace ", false),
		0xBC: NewHIDKeyInfo(0xBC, "Keypad A ", false),
		0xBD: NewHIDKeyInfo(0xBD, "Keypad B ", false),
		0xBE: NewHIDKeyInfo(0xBE, "Keypad C ", false),
		0xBF: NewHIDKeyInfo(0xBF, "Keypad D ", false),
		0xC0: NewHIDKeyInfo(0xC0, "Keypad E ", false),
		0xC1: NewHIDKeyInfo(0xC1, "Keypad F ", false),
		0xC2: NewHIDKeyInfo(0xC2, "Keypad XOR ", false),
		0xC3: NewHIDKeyInfo(0xC3, "Keypad ∧ ", false),
		0xC4: NewHIDKeyInfo(0xC4, "Keypad % ", false),
		0xC5: NewHIDKeyInfo(0xC5, "Keypad < ", false),
		0xC6: NewHIDKeyInfo(0xC6, "Keypad > ", false),
		0xC7: NewHIDKeyInfo(0xC7, "Keypad & ", false),
		0xC8: NewHIDKeyInfo(0xC8, "Keypad && ", false),
		0xC9: NewHIDKeyInfo(0xC9, "Keypad | ", false),
		0xCA: NewHIDKeyInfo(0xCA, "Keypad || ", false),
		0xCB: NewHIDKeyInfo(0xCB, "Keypad : ", false),
		0xCC: NewHIDKeyInfo(0xCC, "Keypad # ", false),
		0xCD: NewHIDKeyInfo(0xCD, "Keypad Space ", false),
		0xCE: NewHIDKeyInfo(0xCE, "Keypad @ ", false),
		0xCF: NewHIDKeyInfo(0xCF, "Keypad ! ", false),
		0xD0: NewHIDKeyInfo(0xD0, "Keypad Memory Store ", false),
		0xD1: NewHIDKeyInfo(0xD1, "Keypad Memory Recall ", false),
		0xD2: NewHIDKeyInfo(0xD2, "Keypad Memory Clear ", false),
		0xD3: NewHIDKeyInfo(0xD3, "Keypad Memory Add ", false),
		0xD4: NewHIDKeyInfo(0xD4, "Keypad Memory Subtract ", false),
		0xD5: NewHIDKeyInfo(0xD5, "Keypad Memory Multiply ", false),
		0xD6: NewHIDKeyInfo(0xD6, "Keypad Memory Divide ", false),
		0xD7: NewHIDKeyInfo(0xD7, "Keypad +/- ", false),
		0xD8: NewHIDKeyInfo(0xD8, "Keypad Clear ", false),
		0xD9: NewHIDKeyInfo(0xD9, "Keypad Clear Entry ", false),
		0xDA: NewHIDKeyInfo(0xDA, "Keypad Binary ", false),
		0xDB: NewHIDKeyInfo(0xDB, "Keypad Octal ", false),
		0xDC: NewHIDKeyInfo(0xDC, "Keypad Decimal ", false),
		0xDD: NewHIDKeyInfo(0xDD, "Keypad Hexadecimal ", false),
		// 0xDE-DF Reserved
		0xE0: NewHIDKeyInfo(0xE0, "Keyboard LeftControl", true),
		0xE1: NewHIDKeyInfo(0xE1, "Keyboard LeftShift", true),
		0xE2: NewHIDKeyInfo(0xE2, "Keyboard LeftAlt", true),
		0xE3: NewHIDKeyInfo(0xE3, "Keyboard Left GUI", true),
		0xE4: NewHIDKeyInfo(0xE4, "Keyboard RightControl", true),
		0xE5: NewHIDKeyInfo(0xE5, "Keyboard RightShift", true),
		0xE6: NewHIDKeyInfo(0xE6, "Keyboard RightAlt", true),
		0xE7: NewHIDKeyInfo(0xE7, "Keyboard Right GUI", true),
	}
	return &HIDKeyboard{make([]byte, 8), keyInfoMap}
}

func (keyboard *HIDKeyboard) PressKey(code uint8) {
	keyInfo := keyboard.keyInfoMap[code]
	keyInfo.Pressed = true
}

func (keyboard *HIDKeyboard) ReleaseKey(code uint8) {
	keyInfo := keyboard.keyInfoMap[code]
	keyInfo.Pressed = false
}

func (keyboard *HIDKeyboard) ReleaseAllKeys() {
	for _, keyInfo := range keyboard.keyInfoMap {
		keyInfo.Pressed = false
	}
}

func (keyboard *HIDKeyboard) GetKeyInfo(code uint8) *HIDKeyInfo {
	return keyboard.keyInfoMap[code]
}

func (keyboard *HIDKeyboard) AddConvKey(code byte, convKey *ConvKeyInfo) {
	keyInfo := keyboard.GetKeyInfo(code)
	keyInfo.convKeyInfoList = append(keyInfo.convKeyInfoList, convKey)
}

func (keyboard *HIDKeyboard) SetupHidPackat() []byte {
	orgModifierFlag := uint8(0)
	// 一旦 data をクリアする
	for index := range keyboard.data {
		keyboard.data[index] = 0
	}
	// modifier をセットする
	index := 2
	for _, val := range keyboard.keyInfoMap {
		if val.Pressed && val.IsModifier {
			orgModifierFlag = orgModifierFlag | val.GetModifierBit()
		}
	}
	// キーの置き換え等を処理する
	modifierFlag := orgModifierFlag
	for _, keyInfo := range keyboard.keyInfoMap {
		if keyInfo.Pressed {
			code := byte(0)
			code, modifierFlag = keyInfo.process(modifierFlag)
			if code > 0 {
				keyboard.data[index] = code
				index++
				if index >= len(keyboard.data) {
					break
				}
			}
		}
	}
	keyboard.data[0] = modifierFlag
	return keyboard.data
}
