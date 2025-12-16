package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tautcony/qart/routers"
)

func main() {
	prod := flag.Bool("prod", false, "run in production mode")
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	// Configure zerolog with custom format: [QART-{level}] time | source | message
	if *prod {
		gin.SetMode(gin.ReleaseMode)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		// Development mode: custom format without colors
		output := zerolog.ConsoleWriter{
			Out:     os.Stdout,
			NoColor: false,
			PartsOrder: []string{
				zerolog.LevelFieldName,
				zerolog.TimestampFieldName,
				zerolog.CallerFieldName,
				zerolog.MessageFieldName,
			},
			FormatLevel: func(i any) string {
				if ll, ok := i.(string); ok {
					if ll == "info" {
						return "[QART]"
					}
					return fmt.Sprintf("[QART-%s]", strings.ToLower(ll))
				}
				return ""
			},
			FormatTimestamp: func(i any) string {
				if ts, ok := i.(string); ok {
					// Parse RFC3339 format and reformat
					if t, err := time.Parse(time.RFC3339, ts); err == nil {
						return fmt.Sprintf("%s |", t.Format("2006/01/02 - 15:04:05"))
					}
				}
				return fmt.Sprintf(" %s |", i)
			},
			FormatCaller: func(i any) string {
				if caller, ok := i.(string); ok {
					// Extract just the filename
					parts := strings.Split(caller, ":")
					if len(parts) > 0 {
						filename := filepath.Base(parts[0])
						return fmt.Sprintf("%s |", filename)
					}
				}
				return ""
			},
			FormatMessage: func(i any) string {
				return fmt.Sprintf("%s", i)
			},
			FormatFieldName: func(i any) string {
				return fmt.Sprintf("%s=", i)
			},
			FormatFieldValue: func(i any) string {
				return fmt.Sprintf("%s", i)
			},
		}

		log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	r := routers.SetupRouter()

	if err := r.Run(":" + *port); err != nil {
		panic(err)
	}
}
