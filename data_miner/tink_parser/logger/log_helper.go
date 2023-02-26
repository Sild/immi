package logger

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func prefix() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return fmt.Sprintf("%s:%s", filename, fn)
}

func Debug(format string, v ...interface{}) {
	log.Debug().Msgf(fmt.Sprintf("[%s] %s", prefix(), format), v...)

}

func Info(format string, v ...interface{}) {
	log.Info().Msgf(fmt.Sprintf("[%s] %s", prefix(), format), v...)

}

func Warn(format string, v ...interface{}) {
	log.Warn().Msgf(fmt.Sprintf("[%s] %s", prefix(), format), v...)
}

func Error(format string, v ...interface{}) {
	log.Error().Msgf(fmt.Sprintf("[%s] %s", prefix(), format), v...)

}
func Fatal(format string, v ...interface{}) {
	// log.Fatal().Msgf(format, v);
	log.Fatal().Msgf(fmt.Sprintf("[%s] %s", prefix(), format), v...)
}
