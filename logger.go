package main

import (
	"bytes"
	"io"
	"fmt"
	"time"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"go-logging/handlers"
)

const (
	OLEVEL = 1 << iota
	OFILE
	OTIME
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

var severityLeveles = map[int]string{DEBUG: "DEBUG", INFO: "INFO", WARNING: "WARNING", ERROR: "ERROR"}

type buffer struct {
	bytes.Buffer
	next *buffer
	tmp  [64]byte
}

type logger struct {
	mu      sync.Mutex
	flags   int
	level   int
	handler io.WriteCloser

	// I am sorry for stealing from Google's glog package =(
	freeList   *buffer
	freeListMu sync.Mutex
}

func (l *logger) getBuffer() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		// point logger's next buffer to next avail after this one
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

func (l *logger) putBuffer(b *buffer) {
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

func (l *logger) flushBuffer(b *buffer) {
	l.handler.Write(b.Bytes())
}

func (l *logger) writeHeader(level int, buf *buffer) {
	severity := severityLeveles[level]
	now := time.Now().Format("02-01-20016 15:04:05")
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
	buf.WriteString(now)
	buf.WriteString(" ")
	buf.WriteString(severity)
	buf.WriteString(" ")
	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString("   ")
}

func (l *logger) print(level int, v ...interface{}) {
	if level >= l.level {
		b := l.getBuffer()
		defer l.putBuffer(b)

		l.writeHeader(level, b)
		fmt.Fprintln(b, v...)

		l.flushBuffer(b)
	}
}

func (l *logger) printf(level int, format string, v ...interface{}) {
	if level >= l.level {
		b := l.getBuffer()
		defer l.putBuffer(b)

		l.writeHeader(level, b)
		fmt.Fprintf(b, format, v...)

		l.flushBuffer(b)
	}
}

// EXPORTED METHODS
func (l *logger) Debug(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.print(DEBUG, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.printf(DEBUG, format, v...)
}

func (l *logger) Info(v ...interface{}) {
        l.mu.Lock()
        defer l.mu.Unlock()

        l.print(INFO, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
        l.mu.Lock()
        defer l.mu.Unlock()

        l.printf(INFO, format, v...)
}

func (l *logger) Warning(v ...interface{}) {
        l.mu.Lock()
        defer l.mu.Unlock()

        l.print(WARNING, v...)
}

func (l *logger) Warningf(format string, v ...interface{}) {
        l.mu.Lock()
        defer l.mu.Unlock()

        l.printf(WARNING, format, v...)
}


func (l *logger) Error(v ...interface{}) {
        l.mu.Lock()
        defer l.mu.Unlock()

        l.print(ERROR, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
        l.mu.Lock()
        defer l.mu.Unlock()

        l.printf(ERROR, format, v...)
}





func New(h io.WriteCloser, level, flags int) *logger {
	l := new(logger)
	l.flags = flags
	l.level = level
	l.handler = h
	return l
}

func main() {
	//log := logger.New(handlers.StreamHandler{}, logger.DEBUG, logger.OLEVEL|logger.OFILE|logger.OTIME)
	log := New(handlers.StreamHandler{}, INFO, OLEVEL|OFILE|OTIME)
	go log.Debug("test from debug")
	go log.Info("test from info")
	go log.Warning("Test from warning")
	go log.Error("Test From error")
	time.Sleep(1 * time.Second)
}
