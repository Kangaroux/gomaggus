package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/realmd/server"
	_ "github.com/lib/pq"
	"github.com/phuslu/log"
)

var (
	flagLogDisable bool
	flagLogVerbose bool
	flagLogLevel   int
)

func init() {
	flag.BoolVar(&flagLogDisable, "nolog", false, "disable all logging")
	flag.BoolVar(&flagLogVerbose, "verbose", false, "use verbose logs (longer timestamp, filename)")
	flag.IntVar(&flagLogLevel, "loglevel", int(log.InfoLevel),
		fmt.Sprintf("minimum error level to log (%d-%d)", log.TraceLevel, log.PanicLevel))
	flag.Parse()

	if flagLogLevel < int(log.TraceLevel) || flagLogLevel > int(log.PanicLevel) {
		fmt.Println("error: invalid log level")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	if flagLogDisable {
		log.DefaultLogger.Writer = log.WriterFunc(func(e *log.Entry) (int, error) { return 0, nil })
	} else {
		log.DefaultLogger = log.Logger{
			Level: log.Level(flagLogLevel),
			Writer: &log.ConsoleWriter{
				ColorOutput: true,
			},
		}

		if flagLogVerbose {
			log.DefaultLogger.Caller = 1
			log.DefaultLogger.TimeFormat = "2006-01-02 15:04:05.000000"
		} else {
			log.DefaultLogger.TimeFormat = "15:04:05.000000"
		}
	}

	db, err := sqlx.Connect(
		"postgres",
		"postgres://gomaggus:password@localhost:5432/gomaggus?sslmode=disable",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	server := server.New(db, server.DefaultListenAddr)
	server.Start()
}
