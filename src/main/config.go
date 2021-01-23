package main

var config struct {
	Mongodb struct {
		Host       string
		Database   string
		Collection string
	}

	/*
		panic (zerolog.PanicLevel, 5)
		fatal (zerolog.FatalLevel, 4)
		error (zerolog.ErrorLevel, 3)
		warn (zerolog.WarnLevel, 2)
		info (zerolog.InfoLevel, 1)
		debug (zerolog.DebugLevel, 0)
		trace (zerolog.TraceLevel, -1)
	*/
	LogLevel string `envconfig:"default=info"`
}
