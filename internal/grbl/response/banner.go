package response

import (
	"errors"
	"fmt"
	"strings"
)

type Banner struct {
	Version string
}

type BannerHandler struct {
	Callback func(banner *Banner) error
}

func (*BannerHandler) Supports(data string) bool {
	return strings.HasPrefix(data, "Grbl")
}

func (h *BannerHandler) Handle(data string) error {
	parts := strings.Split(data, " ")
	if len(parts) < 2 {
		return fmt.Errorf("banner: invalid message: %s", data)
	}

	if h.Callback == nil {
		return errors.New("banner: no callback defined")
	}

	return h.Callback(&Banner{
		Version: parts[1],
	})
}
