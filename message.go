package main

import (
	"io"
	"log"

	"github.com/fatih/color"
)

var green = color.New(color.FgHiGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

// Runes for emoji
const OK string = string('\u2713')
const KO string = string('\u00D7')

type MessageHandler struct {
	coloredOutput bool
}

func NewMessageHandler(coloredOutput bool, out io.Writer) *MessageHandler {

	// Configuring logger
	log.SetFlags(0)
	log.SetOutput(out)

	return &MessageHandler{
		coloredOutput: coloredOutput,
	}
}

func (m *MessageHandler) printSuccess(message string) {
	if m.coloredOutput {
		log.Printf("%v %v\n", green(OK), message)
	} else {
		log.Printf("%v\n", message)
	}
}

func (m *MessageHandler) printError(message string) {
	if m.coloredOutput {
		log.Printf("%v %v\n", red(KO), message)
	}
}

func (m *MessageHandler) printFatal(message string) {
	log.Fatal(message)
}

func (m *MessageHandler) printInfo(message string) {
	log.Println(message)
}
