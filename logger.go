package log


import (
        "io"
        "sync"
)

var (
        l *logger
)

type logger struct {
        mu	sync.Mutex
        handler	io.WriteCloser
}


func Info(v ...interface{}) (n int, err error) {
        l.mu.Lock()
        defer l.mu.Unlock()
        
        return fmt.Fprintln(l.handler, v...)
}

func Infof(s string, v ...interface{}) (n int, err error) {
        l.mu.Lock()
        defer l.mu.Unlock()

        return fmt.Fprintf(l.handler, s, v...)
}


func SetHandler(h Handler) {
        l.handler = h
}


func init() {
        l = new(logger)
        l.SetHandler(ConsoleHandler{})
}
