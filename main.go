/*
Simple fmt wrapper to log with custom formatter
No allocation optimizations
*/
package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	Fatal   = "\033[0;41mFATAL\033[0m"
	Error   = "\033[0;31mERROR\033[0m"
	Warning = "\033[0;33mWARNING\033[0m"
	Info    = "\033[0;34mINFO\033[0m"
	Debug   = "\033[0;37mDEBUG\033[0m"
)

type logger struct {
	mu  sync.Mutex
	out io.Writer
}

var (
	verbose bool
	logging *logger
)

func init() {
	logging = new(logger)
	logging.out = os.Stdout
}

/*
Debugln acts  just like fmt.Println but with custom prefix  formatter
*/
func Debugln(v ...interface{}) {
	if !verbose {
		return
	}

	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Debug)
	fmt.Fprintln(logging.out, v...)
}

/*
Debugf acts  just like fmt.Printf but with custom prefix  formatter
*/
func Debugf(s string, v ...interface{}) {
	if !verbose {
		return
	}

	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Debug)
	fmt.Fprintf(logging.out, s, v...)
}

func Infoln(v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Info)
	fmt.Fprintln(logging.out, v...)
}

func Infof(s string, v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Info)
	fmt.Fprintf(logging.out, s, v...)
}

func Warningln(v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Warning)
	fmt.Fprintln(logging.out, v...)
}

func Warningf(s string, v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Warning)
	fmt.Fprintf(logging.out, s, v...)
}

func Errorln(v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Error)
	fmt.Fprintln(logging.out, v...)
}

func Errorf(s string, v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Error)
	fmt.Fprintf(logging.out, s, v...)
}

func Fatalln(v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Fatal)
	fmt.Fprintln(logging.out, v...)
}

func Fatalf(s string, v ...interface{}) {
	logging.mu.Lock()
	defer logging.mu.Unlock()

	printHeader(Fatal)
	fmt.Fprintf(logging.out, s, v...)
}

func SetOutput(o io.Writer) {
	logging.out = o
}

func EnableDebug() {
	verbose = true
}

func printHeader(s string) {
	// glog mostly
	now := time.Now().Format(time.RFC822)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "????"
		line = 0
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	fmt.Fprintf(logging.out, "%s  %s:%d  %s: ", now, file, line, s)
}
