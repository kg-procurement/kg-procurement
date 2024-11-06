package utils

type AdvancedLogger interface {
	Info(message string)
	Infof(format string, v ...any)
	Debug(message string)
	Debugf(format string, v ...any)
	Error(message string)
	Errorf(format string, v ...any)
	Fatal(message string)
	Fatalf(format string, v ...any)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Info(message string) {
	infoLogger.Println(message)
}

func (c ConsoleLogger) Infof(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func (c ConsoleLogger) Debug(message string) {
	debugLogger.Println(message)
}

func (c ConsoleLogger) Debugf(format string, v ...any) {
	debugLogger.Printf(format, v...)
}

func (c ConsoleLogger) Error(message string) {
	errorLogger.Println(message)
}

func (c ConsoleLogger) Errorf(format string, v ...any) {
	errorLogger.Printf(format, v...)
}

func (c ConsoleLogger) Fatal(message string) {
	fatalLogger.Fatalln(message)
}

func (c ConsoleLogger) Fatalf(format string, v ...any) {
	fatalLogger.Printf(format, v...)
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

var Logger AdvancedLogger = NewConsoleLogger()
