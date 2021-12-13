package grbl

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/gcode"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/grbl/response"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/grbl/usbserial"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

type Grbl struct {
	serial   *usbserial.UsbSerial
	handlers response.ResponseHandlers
	ignore   []*gcode.Field

	State     response.StateType
	StateName string
	WCO       *point.Point
	MPos      *point.Point
	WPos      *point.Point
	Settings  map[uint8]float64

	Version    string
	LastProbe  *point.Point
	LastAlarm  *response.Alarm
	GCodeState *GCodeStates
}

func NewGrbl(device string) (*Grbl, error) {
	serial, err := usbserial.Open(device)
	if err != nil {
		return nil, err
	}

	if err := serial.Flush(); err != nil {
		return nil, err
	}

	rv := &Grbl{
		serial: serial,
		ignore: []*gcode.Field{
			{
				Letter: 'M',
				Value:  0,
			},
			{
				Letter: 'M',
				Value:  1,
			},
		},

		Settings: map[uint8]float64{},
	}
	rv.handlers = []response.ResponseHandler{
		&response.StatusHandler{
			Callback: rv.StatusHandler,
		},
		&response.MessageHandler{
			Callback: rv.MessageHandler,
		},
		&response.AlarmHandler{
			Callback: rv.AlarmHandler,
		},
		&response.BannerHandler{
			Callback: rv.BannerHandler,
		},
		&response.SettingHandler{
			Callback: rv.SettingHandler,
		},
	}

	// it takes some status calls (or a reset) for grbl to retrieve the WCO
	for rv.WCO == nil {
		if err := rv.SendCommands(context.Background(), "?"); err != nil {
			return nil, err
		}
	}

	// populate gcode states
	if err := rv.SendCommands(context.Background(), "$G"); err != nil {
		return nil, err
	}

	return rv, nil
}

func (g *Grbl) Close() error {
	if g.serial == nil {
		return nil
	}
	return g.serial.Close()
}

func (g *Grbl) StatusHandler(status *response.Status) error {
	g.State = status.State
	g.StateName = status.StateName

	if status.WCO != nil {
		g.WCO = status.WCO.Copy()
	}

	if status.MPos != nil {
		g.MPos = status.MPos.Copy()
	}

	if status.WPos != nil {
		g.WPos = status.WPos.Copy()
	}

	return nil
}

func (g *Grbl) MessageHandler(msg *response.Message) error {
	switch msg.Type {
	case "MSG":
		log.Print("message: ", msg.Content)

	case "GC":
		if g.GCodeState == nil {
			var err error
			g.GCodeState, err = NewGCodeStates(msg.Content)
			if err != nil {
				return nil
			}
			log.Print("gcode state: ", *g.GCodeState)
		}

	case "PRB":
		parts := strings.Split(msg.Content, ":")
		if len(parts) != 2 {
			return fmt.Errorf("grbl: failed to parse probe data: %s", msg.Content)
		}

		g.LastProbe = nil
		if parts[1] == "1" {
			pb, err := point.NewFromStringMM(parts[0])
			if err != nil {
				return err
			}
			g.LastProbe = pb
		}

		log.Print("probe point: ", g.LastProbe)

	default:
		log.Print("no handler found for message: ", *msg)
	}

	return nil
}

func (g *Grbl) AlarmHandler(alarm *response.Alarm) error {
	g.LastAlarm = alarm
	log.Print("alarm: ", *alarm)

	return nil
}

func (g *Grbl) BannerHandler(banner *response.Banner) error {
	g.Version = banner.Version
	log.Print("banner: ", *banner)

	return nil
}

func (g *Grbl) SettingHandler(setting *response.Setting) error {
	g.Settings[setting.Key] = setting.Value
	log.Print("setting: ", *setting)

	return nil
}

func (g *Grbl) SendJob(ctx context.Context, j gcode.Job) error {
	for _, l := range j {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := g.SendLine(l); err != nil {
			return err
		}
	}
	return nil
}

func (g *Grbl) SendLine(l gcode.Line) error {
	if len(l) > 0 {
		for _, ign := range g.ignore {
			if l[0].String() == ign.String() {
				log.Printf("grbl: ignoring g-code: %s", l)
				return nil
			}
		}
	}

	if err := g.send(l.String(), true); err != nil {
		return err
	}
	if g.GCodeState != nil {
		g.GCodeState.ProcessLine(l)
	}
	return nil
}

func (g *Grbl) SendRTCommand(cmd string) error {
	if strings.ContainsAny(cmd, "\r\n") {
		return errors.New("grbl: realtime commands should not have newlines")
	}

	// we just write command. response will be catched by the next streaming read.
	return g.serial.WriteLine(strings.TrimSpace(cmd), false)
}

func (g *Grbl) SendCommands(ctx context.Context, cmds string) error {
	scanner := bufio.NewScanner(strings.NewReader(cmds))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := g.send(scanner.Text(), true); err != nil {
			return err
		}
	}

	return nil
}

func (g *Grbl) SendGCodeInline(ctx context.Context, data string) error {
	j, err := gcode.NewJobFromData(data)
	if err != nil {
		return err
	}

	return g.SendJob(ctx, j)
}

func (g *Grbl) send(data string, nl bool) error {
	if err := g.serial.WriteLine(data, nl); err != nil {
		return err
	}

	for {
		line, err := g.serial.ReadLine()
		if err != nil {
			return err
		}

		if line == "" {
			continue
		} else if line == "ok" {
			return nil
		} else if strings.HasPrefix(line, "error:") {
			id, _ := strconv.Atoi(line[6:])
			return Error(id)
		}

		if handler := g.handlers.Lookup(line); handler != nil {
			if err := handler.Handle(line); err != nil {
				log.Printf("error: grbl: %s", err)
			}
		} else {
			log.Printf("warning: grbl: no handler found for response: %s", line)
		}
	}
}
