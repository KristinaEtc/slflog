package slflog

import (
	"github.com/kardianos/osext"
	"github.com/ventu-io/slf"
	"github.com/ventu-io/slog"
	"os"
	"path/filepath"
)

var (
	bhDebug, bhInfo, bhError, bhDebugConsole, bhStdError *Handler
	logfileInfo, logfileDebug, logfileError              *os.File
	lf                                                   slog.LogFactory
	log                                                  slf.StructuredLogger
)

const (
	errorFilename = "error.log"
	infoFilename  = "info.log"
	debugFilename = "debug.log"
)

// Init loggers
func InitLoggers(logpath string, loglvl string) {

	var logHandlers []slog.EntryHandler

	// optionally define the format (this here is the default one)
	//bhInfo.SetTemplate("{{.Time}} [\033[{{.Color}}m{{.Level}}\033[0m] {{.Context}}{{if .Caller}} ({{.Caller}}){{end}}: {{.Message}}{{if .Error}} (\033[31merror: {{.Error}}\033[0m){{end}} {{.Fields}}")

	ConfigWriterOutput(&logHandlers, getLogLevel(loglvl), os.Stderr)

	err := setLogOutput(&logHandlers, logpath)
	if err != nil {
		SafeLog("[go-stomp-server] Error init loggers: " + err.Error() + "\n")
	}

	lf = slog.New()
	lf.SetLevel(slf.LevelDebug)
	lf.SetEntryHandlers(logHandlers...)
	slf.Set(lf)
}

func setLogOutput(logHandlers *[]slog.EntryHandler, logpath string) error {

	pathForLogs, err := getPathForLogDir(logpath)
	if err != nil {
		return err
	}
	exist, err := exists(pathForLogs)
	if err != nil {
		return err
	}
	if !exist {
		err = os.Mkdir(pathForLogs, 0755)
		if err != nil {
			return err
		}
	}

	ConfigFileOutput(logHandlers, slf.LevelDebug, filepath.Join(pathForLogs, debugFilename))
	ConfigFileOutput(logHandlers, slf.LevelInfo, filepath.Join(pathForLogs, infoFilename))
	ConfigFileOutput(logHandlers, slf.LevelError, filepath.Join(pathForLogs, errorFilename))

	return nil
}

func getPathForLogDir(logpath string) (string, error) {

	if filepath.IsAbs(logpath) == true {
		return logpath, nil
	} else {
		filename, err := osext.Executable()
		if err != nil {
			return "", err
		}

		fpath := filepath.Dir(filename)
		fpath = filepath.Join(fpath, logpath)
		return fpath, nil
	}
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getLogLevel(lvl string) slf.Level {

	switch lvl {
	case slf.LevelDebug.String():
		return slf.LevelDebug

	case slf.LevelInfo.String():
		return slf.LevelInfo

	case slf.LevelWarn.String():
		return slf.LevelWarn

	case slf.LevelError.String():
		return slf.LevelError

	case slf.LevelFatal.String():
		return slf.LevelFatal
	case slf.LevelPanic.String():
		return slf.LevelPanic
	default:
		return slf.LevelDebug
	}
}
