package cmdpkg

import (
	"fmt"
	"log"
	"os"
)

type Printer struct {
	format_type string

	WarnLogger *log.Logger
	ErrLogger  *log.Logger

	Tab string
}

func NewPrinter() *Printer {
	p := Printer{
		format_type: "default",
		WarnLogger:  log.New(os.Stderr, ParseColors(".BOLD.PINKwarning.RESET: "), log.Lmsgprefix),
		ErrLogger:   log.New(os.Stderr, ParseColors(".BOLD.REDerror.RESET: "), log.Lmsgprefix),
	}
	return &p
}

func (p *Printer) SetFormat(s string) { p.format_type = s }
func (p Printer) Format() string      { return p.format_type }

func (p Printer) Fatal(format string, v ...interface{}) {
	p.ErrLogger.Fatalf(format, v...)
}

func (p Printer) Error(format string, v ...interface{}) {
	p.ErrLogger.Printf(format, v...)
}

func (p Printer) Warning(format string, v ...interface{}) {
	p.WarnLogger.Printf(format, v...)
}

type Printable interface {
	JSON() string
	CSV() string
	String() string
}

func (p Printer) Print(obj Printable) {
	switch p.format_type {
	case "json":
		fmt.Println(ParseColors(obj.JSON()))

	case "csv":
		fmt.Println(ParseColors(obj.CSV()))

	default:
		fmt.Println(ParseColors(obj.String()))
	}
}
