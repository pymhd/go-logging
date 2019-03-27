package handlers

import (
        "os"
)

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

