package logging

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logger *zap.Logger
)

func Setup() {
	checkLogFile()
}

func createLogger(filePath string) {
	var err error

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{filePath}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = eSTimeEncoder
	//cfg.EncoderConfig = zapcore.EncoderConfig{
	//	TimeKey:        "timestamp",
	//	LineEnding:     zapcore.DefaultLineEnding,
	//	EncodeTime:     eSTimeEncoder,
	//	EncodeDuration: zapcore.SecondsDurationEncoder,
	//}
	logger, err = cfg.Build()
	if err != nil {
		log.Fatalf("logger.Setup err: %v", err)
	}
}

func eSTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05-0700"))
}

func Info(v interface{}) {
	checkLogFile()
	logger.Info("", zap.Any("data", v))
}

func Error(v interface{}) {
	checkLogFile()
	logger.Error("", zap.Any("data", v))
}

func checkLogFile() {
	fp := "./storage/logs/"
	fn := time.Now().Format("2006-01-02") + ".log"
	filePath := fp + fn

	var err error
	f, err := file.MustOpen(fn, fp)
	if err != nil {
		log.Fatalf("logger.Setup err: %v", err)
	}
	_ = f.Close()

	createLogger(filePath)
	go func() {
		c := make(chan os.Signal, 1)
		for {
			signal.Notify(c, syscall.SIGUSR1)
			s := <-c
			logger = nil
			createLogger(filePath)
			fmt.Println("Got signal:", s)
		}
	}()
}
