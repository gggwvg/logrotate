package main

import (
	"github.com/gggwvg/logrotate"
)

func main() {
	rotateBy10k()
	rotateByHourly()
}

func rotateBy10k() {
	logger, err := logrotate.NewLogger(
		logrotate.File("/tmp/test_log/rotate_size_10k.log"),
		logrotate.RotateSize("10k"),
		logrotate.Compress(true),
	)
	if err != nil {
		panic(err)
	}
	defer logger.Close()
	for i := 0; i < 10000; i++ {
		_, err = logger.Write([]byte("This is a test message."))
		if err != nil {
			panic(err)
		}
	}
}

func rotateByHourly() {
	// TODO
}
