package logging

// Logger is an interface that could be implemented by other loggers to be used in the parser
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
