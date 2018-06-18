package logging

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"strings"
	"runtime"
)

var VERBOSE bool

func init() {
	flag.BoolVar(&VERBOSE, "v", false, "verbose")
}

func New(out io.Writer) *Logger {
	return &Logger{Out: out}
}

type Logger struct {
	sync.Mutex
	Verbose bool
	Out     io.Writer
}

func (l *Logger) formatHeader(severity string) {
	// glog mostly
	now := time.Now().Format(time.RFC822)
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "????"
		line = 0
	} else{
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	fmt.Fprintf(l.Out, "%s  %s:%d  %s: ", now, file, line, severity)
}

func (l *Logger) Infoln(v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("\033[0;34mINFO\033[0m")
	fmt.Fprintln(l.Out, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("INFO")
	fmt.Fprintf(l.Out, format, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("\033[0;31mERROR\033[0m")
	fmt.Fprintln(l.Out, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("\033[0;31mERROR\033[0m")
	fmt.Fprintf(l.Out, format, v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	if VERBOSE {
		l.Lock()
		defer l.Unlock()
		l.formatHeader("\033[0;37mDEBUG\033[0m")
		fmt.Fprintln(l.Out, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if VERBOSE {
		l.Lock()
		defer l.Unlock()
		l.formatHeader("\033[0;37mDEBUG\033[0m")
		fmt.Fprintf(l.Out, format, v...)
	}
}

func (l *Logger) Warningln(v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("WARNING")
	fmt.Fprintln(l.Out, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("WARNING")
	fmt.Fprintf(l.Out, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("FATAL")
	fmt.Fprintf(l.Out, format, v...)
	fmt.Fprintln(l.Out, "Exit with status 1")
	os.Exit(1)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.formatHeader("FATAL")
	fmt.Fprintln(l.Out, v...)
	fmt.Fprintln(l.Out, "Exit with status 1")
	os.Exit(1)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Debugln(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.Infoln(v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.Warningln(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Errorln(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Fatalln(v...)
}
