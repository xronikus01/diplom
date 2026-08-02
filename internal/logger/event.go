package logger

type EventLogger interface {
	Log(event string)
}
