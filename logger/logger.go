package logger

import (
	"fmt"
	"io"
	"log"

	"github.com/fatih/color"
)

// global variables holding loggers
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// Init of Logger
func Init(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
	green := color.New(color.FgGreen).SprintFunc()
	Info = log.New(infoHandle,
		fmt.Sprintf("%s ", green("Info:")),
		log.Ldate|log.Ltime|log.Lshortfile)

	yellow := color.New(color.FgYellow).SprintFunc()
	Warning = log.New(
		warningHandle,
		fmt.Sprintf("%s ", yellow("Warning:")),
		log.Ldate|log.Ltime|log.Lshortfile)

	red := color.New(color.FgRed).SprintFunc()
	Error = log.New(errorHandle,
		fmt.Sprintf("%s ", red("ERROR:")),
		log.Ldate|log.Ltime|log.Lshortfile)
}
