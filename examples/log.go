package main

import (
	"log"

	"github.com/gggwvg/logrotate"
)

func main() {
	// default
	logger, err := logrotate.NewLogger()
	if err != nil {
		panic(err)
	}
	log.SetOutput(logger)
	log.Println("default")
	logger.Close()

	// specify a log file
	opts := []logrotate.OptionFunc{
		logrotate.Filename("/tmp/test.log"),
	}
	logger, err = logrotate.NewLogger(opts...)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logger)
	log.Println("log to test.log")
	logger.Close()

	// rotate via time period
	opts = append(opts, logrotate.RotatePeriod(logrotate.PeriodDaily))
	logger, err = logrotate.NewLogger(opts...)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logger)
	log.Println("rotate by daily")
	logger.Close()

	// rotate via file size and time period
	opts = append(opts, logrotate.RotateSize("100m"))
	logger, err = logrotate.NewLogger(opts...)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logger)
	log.Println("rotate by daily and file size 100m")
	logger.Close()
}
