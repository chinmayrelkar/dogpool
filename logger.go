package dogpool

import "log"

type Logger interface {
	Error(string, ...any)
	Warn(string, ...any)
	Info(string, ...any)
	Debug(string, ...any)
	Trace(string, ...any)
}

var _ Logger = &logger{}

type logger struct {
	log log.Logger
}

func (l *logger) Error(msg string, args ...any) { l.log.Println("Error:", msg, args) }
func (l *logger) Warn(msg string, args ...any)  { l.log.Println("Warn:", msg, args) }
func (l *logger) Info(msg string, args ...any)  { l.log.Println("Info:", msg, args) }
func (l *logger) Debug(msg string, args ...any) { l.log.Println("Debug:", msg, args) }
func (l *logger) Trace(msg string, args ...any) { l.log.Println("Trace:", msg, args) }

func NewLogger() Logger {
	return &logger{
		log: *log.Default(),
	}
}
