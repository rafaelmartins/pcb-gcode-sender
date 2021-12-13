package response

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Setting struct {
	Key   uint8
	Value float64
}

type SettingHandler struct {
	Callback func(setting *Setting) error
}

func (*SettingHandler) Supports(data string) bool {
	return data[0] == '$' && strings.Contains(data, "=")
}

func (h *SettingHandler) Handle(data string) error {
	parts := strings.SplitN(data[1:], "=", 2)

	if len(parts) != 2 {
		return fmt.Errorf("setting: invalid message: %s", data)
	}

	if h.Callback == nil {
		return errors.New("setting: no callback defined")
	}

	key, err := strconv.ParseUint(parts[0], 10, 8)
	if err != nil {
		return err
	}

	value, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}

	return h.Callback(&Setting{
		Key:   uint8(key),
		Value: value,
	})
}
