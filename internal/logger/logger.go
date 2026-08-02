package logger

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	events chan string
	wg     sync.WaitGroup
	file   *os.File
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		events: make(chan string, 100),
		file:   file,
	}, nil
}

func (l *Logger) Start(ctx context.Context) {
	l.wg.Add(1)

	go func() {
		defer l.wg.Done()

		logger := log.New(l.file, "", log.LstdFlags)

		for {
			select {
			case <-ctx.Done():
				for {
					select {
					case event := <-l.events:
						time.Sleep(2 * time.Second)
						logger.Println(event)
					default:
						_ = l.file.Close()
						return
					}
				}

			case event := <-l.events:
				time.Sleep(2 * time.Second)
				logger.Println(event)
			}
		}
	}()
}

func (l *Logger) Log(event string) {
	select {
	case l.events <- event:
	default:
	}
}

func (l *Logger) Wait() {
	l.wg.Wait()
}
