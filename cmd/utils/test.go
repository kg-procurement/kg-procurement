package utils

type AdvancedLogger interface {
	Info(message string)
	Debug(message string)
	Error(message string)
	Panic(message string)
	Fatal(message string)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Info(message string) {
	InfoLogger.Println(message)
}

func (c ConsoleLogger) Debug(message string) {
	DebugLogger.Println(message)
}

func (c ConsoleLogger) Error(message string) {
	ErrorLogger.Println(message)
}

func (c ConsoleLogger) Panic(message string) {
	PanicLogger.Panicln(message)
}

func (c ConsoleLogger) Fatal(message string) {
	FatalLogger.Fatalln(message)
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}
