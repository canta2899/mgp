package main

import (
	"io"
	"log"
	"strings"

	"github.com/fatih/color"
)

type MessageHandler struct {
	ColoredOutput bool
	Logger        *log.Logger
	Ok            string
	Ko            string
	OkColor       func(a ...interface{}) string
	KoColor       func(a ...interface{}) string
}

func NewMessageHandler(coloredOutput bool, out io.Writer) *MessageHandler {
	return &MessageHandler{
		ColoredOutput: coloredOutput,
		Logger:        log.New(out, "", 0),
		Ok:            string("\u2713"),
		Ko:            string("\u00D7"),
		OkColor:       color.New(color.FgHiGreen).SprintFunc(),
		KoColor:       color.New(color.FgRed).SprintFunc(),
	}
}

func (m *MessageHandler) printNoColor(message ...string) {
	m.Logger.Println(strings.Join(message, " "))
}

func (m *MessageHandler) PrintSuccess(message ...string) {
	if m.ColoredOutput {
		m.Logger.Printf("%v %v\n", m.OkColor(m.Ok), strings.Join(message, " "))
	} else {
		m.printNoColor(message...)
	}
}

func (m *MessageHandler) PrintError(message ...string) {
	if m.ColoredOutput {
		m.Logger.Printf("%v %v\n", m.KoColor(m.Ko), strings.Join(message, " "))
	} else {
		m.printNoColor(message...)
	}
}

func (m *MessageHandler) PrintFatal(message ...string) {
	m.Logger.Fatal(strings.Join(message, " "))
}

func (m *MessageHandler) PrintInfo(message ...string) {
	m.printNoColor(message...)
}
