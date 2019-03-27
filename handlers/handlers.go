packah main

import (
        "os"
)

// Null (debug logs)
type NullHandler struct {}


func (nh NullHandler) Write(p []byte) (n int, err error) {
        return 0, nil
}

func (nh NullHandler) Close() error {
        return nil 
}


//File
type FileHandler struct {
        fd 	*os.File
}

func (fh FileHandler) Write(p []byte) (n int, err error) {
        return fh.fd.Write(p)
}

func (fh FileHandler) Close() error {
        return fh.fd.Close()
}

