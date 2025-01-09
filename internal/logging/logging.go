package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

func LogSetup(lvl string) {
	l, err := log.ParseLevel(lvl)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(
		&log.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
			},
		},
	)

	if l == log.DebugLevel {
		log.SetLevel(l)
		//log.SetReportCaller(true)
	} else {
		log.SetLevel(l)
	}

	file, err := os.OpenFile("/home/kaw/parserService/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//log.SetOutput(os.Stdout)
	log.SetOutput(file)
}
