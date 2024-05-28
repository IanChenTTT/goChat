package main

import (
	"fmt"
	"log/slog"
	"os"
    "path/filepath"
)

func main() {
    logger := slog.Default()
    logger.Info("Hello world")
    ex, err := os.Executable()
    if err != nil{
        logger.Warn("some went wrong %v",err)
        panic(err)
    }
    exPath := filepath.Dir(ex)
    entries, err := os.ReadDir(exPath)
    if err != nil{
        logger.Warn("some went wrong %v",err)
        panic(err)
    }
    fmt.Println(exPath)
}
