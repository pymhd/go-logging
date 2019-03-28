package logger 


import (
	"fmt"
	"bytes"
	"time"
	"sync"
	"runtime"
	"strconv"
	"strings"
	"github.com/pymhd/go-logging/handlers"
)

const (
	defaultLayout = "02-01-2006 15:04:05"
)

const (
	OTIME = 1 << iota
	OLEVEL
	OFILE
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

var severityLeveles = map[int]string{DEBUG: "DEBUG", INFO: "INFO", WARNING: "WARNING", ERROR: "ERROR"}
var registeredLoggers = make(map[string]*Logger)


type buffer struct {
	bytes.Buffer
	next *buffer
	tmp  [64]byte
l}

type Logger struct {
	mu      sync.Mutex
	flags   int
	level   int
	handler handlers.Handler

	// I am sorry for stealing from Google's glog package =(
	freeList   *buffer
	freeListMu sync.Mutex
}

func (l *Logger) getBuffer() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		// point Logger's next buffer to next avail after this one
		l.freeList = b.next
	}
	l.freeListMu.Unlock()
	if b == nil {
		b = new(buffer)
	} else {
		// to reset buffer and disconnect from next buffer in list (isolate anf flush)
		b.next = nil
		b.Reset()
	}
	return b
}

func (l *Logger) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		// let for GC
		return
	}
	l.freeListMu.Lock()
	//to insert buffer back in chain, after getBuffer gets it it will isolate it again and flush
	b.next = l.freeList
	l.freeList = b
	l.freeListMu.Unlock()
}

func (l *Logger) flushBuffer(b *buffer) {
	l.handler.Write(b.Bytes())
}

func (l *Logger) writeHeader(level int, buf *buffer) {
	if OTIME&l.flags > 0 {
		now := time.Now().Format(defaultLayout)
		buf.WriteString(now)
		buf.WriteString(" ")
	}
	if OLEVEL&l.flags > 0 {
		severity := severityLeveles[level]
		buf.WriteString(severity)
		buf.WriteString(" ")
	}
	if OFILE&l.flags > 0 {
		_, file, line, ok := runtime.Caller(3)
		if !ok {
			file = "????"
			line = 0
		} else {
			slash := strings.LastIndex(file, "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
		}

		buf.WriteString(file)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(line))
		buf.WriteString(" ")
	}
	buf.WriteString("  ")
}

func (l *Logger) print(level int, v ...interface{}) {
	if level >= l.level {
		b := l.getBuffer()
		defer l.putBuffer(b)

		l.writeHeader(level, b)
		fmt.Fprintln(b, v...)

		l.flushBuffer(b)
	}
}

func (l *Logger) printf(level int, format string, v ...interface{}) {
	if level >= l.level {
		b := l.getBuffer()
		defer l.putBuffer(b)

		l.writeHeader(level, b)
		fmt.Fprintf(b, format, v...)

		l.flushBuffer(b)
	}
}


// EXPORTED METHODS
func (l *Logger) Debug(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.print(DEBUG, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.printf(DEBUG, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.print(INFO, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.printf(INFO, format, v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.print(WARNING, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.printf(WARNING, format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.print(ERROR, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.printf(ERROR, format, v...)
}

func New(name, string, h handlers.Handler, level, flags int) *Logger {
	existingLogger, ok := registeredLoggers[name]
	if ok {
		return existingLogger
	}
	l := new(Logger)
	l.flags = flags
	l.level = level
	l.handler = h
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			l.mu.Lock()
			l.handler.Flush()
			l.mu.Unlock()
		}
	}()
	existingLogger[name] = l
	return l
}
