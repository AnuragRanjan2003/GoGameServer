package logs

import (
	"context"
	"log"
	"os"
)

const DUMP_FILE = "debug/logdump.txt"

type Logger struct {
	logChannel     chan *LogMessage
	logFile        *os.File
	current_offset int64
	context        context.Context
}

func NewLogger(ctx context.Context) *Logger {
	_file, err := initLogFile(DUMP_FILE)
	if err != nil {
		return nil
	}

	return &Logger{
		logChannel:     make(chan *LogMessage, 100),
		logFile:        _file,
		current_offset: 0,
		context:        ctx,
	}
}

func (l Logger) Start() {
	l.runDumper(l.context)
}

func initLogFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
}

func (l *Logger) PushLog(m *LogMessage) {
	go func() {
		l.logChannel <- m
	}()
}

func (l *Logger) runDumper(ctx context.Context) {
	_cache := make([]LogMessage, 109)
	for {
		select {
		case <-ctx.Done():
			log.Println("Application closed")
			return
		default:
			m := <-l.logChannel
			_cache = append(_cache, *m)

			if len(_cache) == 100 {
				go l.dumpLogsAndClear(_cache, ctx)
			}
		}
	}
}

func (l *Logger) dumpLogsAndClear(list []LogMessage, ctx context.Context) {
	for _, m := range list {
		select {
		case <-ctx.Done():
			log.Println("Application closed")
			return
		default:
			l.logFile.WriteAt(m.Bytes(), l.current_offset)
			l.current_offset++
		}
	}
}
