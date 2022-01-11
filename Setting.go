package main

import (
	"encoding/json"
	"os"
)

type SettingSwitchKey struct {
	// 有効かどうか。 nil の場合は有効。
	On *bool
	// 置き換え元コード
	Src byte
	// 置き換え先コード
	Dst byte
}

type Setting struct {
	InputKeyboardName *string
	SwitchKeys        []SettingSwitchKey
	ConvKeyMap        map[string][]ConvKeyInfo
}

func load(path string) (*Setting, error) {
	fileObj, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if fileInfo, err := os.Stat(path); err != nil {
		return nil, err
	} else {
		buf := make([]byte, fileInfo.Size())
		if _, err := fileObj.Read(buf); err != nil {
			return nil, err
		}
		var setting Setting
		err := json.Unmarshal(buf, &setting)
		return &setting, err
	}
}
