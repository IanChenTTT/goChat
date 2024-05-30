package server

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Log() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	l.Info("slog init")
	ex, err := os.Executable()
	Check(err, l)

	exPath := filepath.Dir(ex)
	ents, err := os.ReadDir(exPath)
	Check(err, l)

	var contain = false
	for _, ent := range ents {
		if strings.Contains(ent.Name(), "log") {
			contain = true
		}
	}
	if contain == false {
		l.Warn("log directory can't find")
		os.Exit(1)
	}
    fmt.Println("Executable path:",exPath)

	fs, err := os.OpenFile("log/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Check(err, l)

    l.Debug("Admim", slog.String("action","Create Log file"))
	defer func(fs *os.File) {
		err := fs.Close()
		Check(err, l)
	}(fs)
}

func Check(e error, l *slog.Logger) {
	if e != nil {
		l.Error("some went wrong ", slog.String("err", e.Error()))
		panic(e)
	}
}
