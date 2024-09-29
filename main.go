package main

import (
	"RISK-CodeConflict/ai"
	"RISK-CodeConflict/core"
	"RISK-CodeConflict/gui"
	"RISK-CodeConflict/remote"
	"flag"
	"fmt"
	"image/color"
	"os"
	"time"
)

const VERSION = "1.0"

func main() {
	// print program name and version
	programName := fmt.Sprintf("RISK: Code Conflict v%s", VERSION)
	println(programName)
	println()

	// server flags
	var host string
	var port string
	var aiPlayer int
	var remotePlayer int
	var humanPlayer int
	var noLog bool
	var autoRedraw bool

	// parse
	flag.StringVar(&host, "host", "localhost", "Server host")
	flag.StringVar(&port, "port", "1234", "Server port")
	flag.IntVar(&aiPlayer, "ai", 0, "add RandomAI players")
	flag.IntVar(&remotePlayer, "remote", 0, "waiting for remote client-AI players")
	flag.IntVar(&humanPlayer, "human", 0, "add human players (control via the server gui)")
	flag.BoolVar(&noLog, "noLog", false, "disables combat output in the server log")
	flag.BoolVar(&autoRedraw, "autoRedraw", false, "forces the gui to redraw every frame")
	flag.Parse()

	// player, host and port
	if aiPlayer+remotePlayer+humanPlayer < 1 || host == "" || port == "" {
		flag.Usage()
		os.Exit(6)
	}

	//---------------------------------------------------------------------------------------------------

	// run program
	colorIndex := 0
	playerColors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 0, G: 255, B: 0, A: 255},     // Green
		{R: 0, G: 0, B: 255, A: 255},     // Blue
		{R: 255, G: 255, B: 0, A: 255},   // Yellow
		{R: 255, G: 0, B: 255, A: 255},   // Magenta
		{R: 254, G: 255, B: 255, A: 255}, // White
		{R: 1, G: 0, B: 0, A: 255},       // dummy
		{R: 2, G: 0, B: 0, A: 255},       // dummy
		{R: 3, G: 0, B: 0, A: 255},       // dummy
		{R: 4, G: 0, B: 0, A: 255},       // dummy
		{R: 5, G: 0, B: 0, A: 255},       // dummy
		{R: 6, G: 0, B: 0, A: 255},       // dummy
		{R: 7, G: 0, B: 0, A: 255},       // dummy
		{R: 8, G: 0, B: 0, A: 255},       // dummy
		{R: 9, G: 0, B: 0, A: 255},       // dummy
		{R: 10, G: 0, B: 0, A: 255},      // dummy
		{R: 11, G: 0, B: 0, A: 255},      // dummy
	}

	//-----------------------------------------------

	// new world
	w := core.NewWorld()
	w.NoLog = noLog

	// add human player
	for i := 0; i < humanPlayer; i++ {
		name := fmt.Sprintf("Human %d", i+1)
		if err := w.AddPlayer(name, playerColors[colorIndex]); err != nil {
			panic(err)
		}
		colorIndex++
	}

	// start server
	go remote.RunServer(host, port, w, aiPlayer+remotePlayer+humanPlayer)
	time.Sleep(200 * time.Millisecond)

	// add local AIs
	for i := 0; i < aiPlayer; i++ {
		name := fmt.Sprintf("RandomAI %d", i+1)
		go ai.Play(host, port, name, playerColors[colorIndex])
		colorIndex++
	}

	// run gui (blocking)
	if err := gui.RunGUI(1778, 1000, programName, w, autoRedraw); err != nil {
		panic(err)
	}
}
