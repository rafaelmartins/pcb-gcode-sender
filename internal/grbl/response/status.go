package response

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

type StateType int

const (
	StateIdle StateType = iota
	StateRun
	StateHold
	StateJog
	StateAlarm
	StateDoor
	StateCheck
	StateHome
	StateSleep
)

var (
	stateMap = map[string]StateType{
		"Idle":  StateIdle,
		"Run":   StateRun,
		"Hold":  StateHold,
		"Jog":   StateJog,
		"Alarm": StateAlarm,
		"Door":  StateDoor,
		"Check": StateCheck,
		"Home":  StateHome,
		"Sleep": StateSleep,
	}
)

type Status struct {
	State     StateType
	StateName string
	WCO       *point.Point
	MPos      *point.Point
	WPos      *point.Point
	Other     map[string]string
}

type StatusHandler struct {
	wco      *point.Point
	Callback func(status *Status) error
}

func (*StatusHandler) Supports(data string) bool {
	return data[0] == '<' && data[len(data)-1] == '>'
}

func (h *StatusHandler) Handle(data string) error {
	msg := data[1 : len(data)-1]
	fields := strings.Split(msg, "|")
	if len(fields) == 0 {
		return fmt.Errorf("status: invalid status report: %s", data)
	}

	stateParts := strings.Split(fields[0], ":")
	// FIXME: handle sub-states

	state, ok := stateMap[stateParts[0]]
	if !ok {
		return fmt.Errorf("status: invalid state: %s", fields[0])
	}

	rv := &Status{
		State:     state,
		StateName: fields[0],
		Other:     map[string]string{},
	}

	if err := func() error {
		for _, field := range fields[1:] {
			parts := strings.Split(field, ":")
			if len(parts) != 2 {
				return fmt.Errorf("status: invalid data field: %s", field)
			}

			var err error

			switch parts[0] {
			case "WCO":
				h.wco, err = point.NewFromStringMM(parts[1])
				if err != nil {
					return err
				}

			case "MPos":
				rv.MPos, err = point.NewFromStringMM(parts[1])
				if err != nil {
					return err
				}

				defer func() {
					if h.wco != nil {
						rv.WPos = rv.MPos.Sub(h.wco)
						rv.WCO = h.wco.Copy()
					}
				}()

			case "WPos":
				rv.WPos, err = point.NewFromStringMM(parts[1])
				if err != nil {
					return err
				}

				defer func() {
					if h.wco != nil {
						rv.MPos = rv.MPos.Add(h.wco)
						rv.WCO = h.wco.Copy()
					}
				}()

			default:
				rv.Other[parts[0]] = strings.TrimSpace(parts[1])
			}
		}

		return nil
	}(); err != nil {
		return err
	}

	if h.Callback == nil {
		return errors.New("status: no callback defined")
	}

	return h.Callback(rv)
}
