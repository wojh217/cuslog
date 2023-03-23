package cuslog

import (
	"fmt"
	"strings"
	"time"
)

const (
	color_red     = uint8(iota + 91)
	color_green   //	绿
	color_yellow  //	黄
	color_blue    // 	蓝
	color_magenta //	洋红
)

type ColorFormatter struct {
	IgnoreBasicFields bool
}

func (f *ColorFormatter) Format(e *Entry) error {
	if !f.IgnoreBasicFields {
		f.WriteString(e, fmt.Sprintf("%s %s ", e.Time.Format(time.RFC3339), LevelNameMapping[e.Level]))
		if e.File != "" {
			short := e.File[strings.LastIndex(e.File, "/")+1:]

			f.WriteString(e, fmt.Sprintf("%s:%d", short, e.Line))
		}
		f.WriteString(e, " ")
	}

	switch e.Format {
	case FmtEmptySeparate:
		f.WriteString(e, fmt.Sprint(e.Args...))
	default:
		f.WriteString(e, fmt.Sprintf(e.Format, e.Args...))
	}
	f.WriteString(e, "\n")

	return nil
}

func (f *ColorFormatter) WriteString(e *Entry, s string) {
	switch e.Level {
	case DebugLevel:
		e.Buffer.WriteString(green(s))
	case InfoLevel:
		e.Buffer.WriteString(blue(s))
	case WarnLevel:
		e.Buffer.WriteString(yellow(s))
	case ErrorLevel:
		e.Buffer.WriteString(red(s))
	case PanicLevel:
		e.Buffer.WriteString(magenta(s))
	case FatalLevel:
		e.Buffer.WriteString(magenta(s))
	default:
		e.Buffer.WriteString(s)
	}
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}
func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}
func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}
