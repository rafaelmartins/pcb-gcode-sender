package response

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	AlarmHardLimitError     = 1
	AlarmSoftLimitError     = 2
	AlarmAbortCycle         = 3
	AlarmProbeFailInitial   = 4
	AlarmProbeFailContact   = 5
	AlarmHomingFailReset    = 6
	AlarmHomingFailDoor     = 7
	AlarmHomingFailPulloff  = 8
	AlarmHomingFailApproach = 9
)

var (
	alarmMap = map[uint8]string{
		AlarmHardLimitError:     "Hard limit has been triggered. Machine position is likely lost due to sudden halt. Re-homing is highly recommended.",
		AlarmSoftLimitError:     "Soft limit alarm. G-code motion target exceeds machine travel. Machine position retained. Alarm may be safely unlocked.",
		AlarmAbortCycle:         "Reset while in motion. Machine position is likely lost due to sudden halt. Re-homing is highly recommended.",
		AlarmProbeFailInitial:   "Probe fail. Probe is not in the expected initial state before starting probe cycle when G38.2 and G38.3 is not triggered and G38.4 and G38.5 is triggered.",
		AlarmProbeFailContact:   "Probe fail. Probe did not contact the workpiece within the programmed travel for G38.2 and G38.4.",
		AlarmHomingFailReset:    "Homing fail. The active homing cycle was reset.",
		AlarmHomingFailDoor:     "Homing fail. Safety door was opened during homing cycle.",
		AlarmHomingFailPulloff:  "Homing fail. Pull off travel failed to clear limit switch. Try increasing pull-off setting or check wiring.",
		AlarmHomingFailApproach: "Homing fail. Could not find limit switch within search distances. Try increasing max travel, decreasing pull-off distance, or check wiring.",
	}
)

type Alarm struct {
	Code    uint8
	Message string
}

func (a *Alarm) String() string {
	return a.Message
}

type AlarmHandler struct {
	Callback func(alarm *Alarm) error
}

func (*AlarmHandler) Supports(data string) bool {
	return strings.HasPrefix(data, "ALARM:")
}

func (h *AlarmHandler) Handle(data string) error {
	code64, err := strconv.ParseUint(data[6:], 10, 8)
	if err != nil {
		return err
	}
	code := uint8(code64)

	msg, ok := alarmMap[code]
	if !ok {
		return fmt.Errorf("alarm: invalid error code: %d", code)
	}

	if h.Callback == nil {
		return errors.New("alarm: no callback defined")
	}

	return h.Callback(&Alarm{
		Code:    code,
		Message: msg,
	})
}
