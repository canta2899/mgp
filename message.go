package main

import (
	"io"
	"log"
	"strings"

	"github.com/fatih/color"
)

type MessageHandler struct {
	coloredOutput bool
	logger        *log.Logger
	Ok            string
	Ko            string
	OkColor       func(a ...interface{}) string
	KoColor       func(a ...interface{}) string
}

func NewMessageHandler(coloredOutput bool, out io.Writer) *MessageHandler {
	return &MessageHandler{
		coloredOutput: coloredOutput,
		logger:        log.New(out, "", 0),
		Ok:            string("\u2713"),
		Ko:            string("\u00D7"),
		OkColor:       color.New(color.FgHiGreen).SprintFunc(),
		KoColor:       color.New(color.FgRed).SprintFunc(),
	}
}

func printNoColor(message ...string) {
	log.Println(message)
}

func (m *MessageHandler) PrintSuccess(message ...string) {
	if m.coloredOutput {
		log.Printf("%v %v\n", m.OkColor(m.Ok), strings.Join(message, " "))
	} else {
		printNoColor(message...)
	}
}

func (m *MessageHandler) PrintError(message ...string) {
	if m.coloredOutput {
		log.Printf("%v %v\n", m.KoColor(m.Ko), strings.Join(message, " "))
	} else {
		printNoColor(message...)
	}
}

func (m *MessageHandler) PrintFatal(message ...string) {
	log.Fatal(strings.Join(message, " "))
}

func (m *MessageHandler) PrintInfo(message ...string) {
	printNoColor(message...)
}
