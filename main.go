package main

import (
	"log"
	"os"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/grbl"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/shell"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("serial device required")
	}

	if len(os.Args) > 2 {
		if err := os.Chdir(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	}

	g, err := grbl.NewGrbl(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	a := &actions.Actions{
		Grbl: g,
	}
	if err := shell.Run(a); err != nil {
		log.Fatal(err)
	}
}
