package server

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Log() {
	logger := slog.Default()
	logger.Info("slog init")
	ex, err := os.Executable()
	if err != nil {
		logger.Warn("some went wrong %v", err)
		panic(err)
	}
	exPath := filepath.Dir(ex)
	ents, err := os.ReadDir(exPath)
	if err != nil {
		logger.Warn("some went wrong %v", err)
		panic(err)
	}
  var contain = false
	for _, ent := range ents {
    if strings.Contains(ent.Name(), "log"){
      contain = true
    }
	}
  if contain == false{
    logger.Warn("log directory can't find")
    os.Exit(1);
  }
	fmt.Println(exPath)
}
