package handlers

import (
        "os"
)

//Console
type StreamHandler struct {}

func(ch StreamHandler) Write(p []byte) (n int, err error) {
        return os.Stdout.Write(p)
}


func(ch StreamHandler) Close() error {
        return nil
}
