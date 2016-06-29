package slflog

import (
	"encoding/json"
	"fmt"
	"github.com/kardianos/osext"
	"github.com/ventu-io/slf"
	"github.com/ventu-io/slog"
	"io/ioutil"
	"os"
	"path"
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

type Config struct {
	ConsoleLvl string `json:ConsoleLvl`
	Logpath    string `json:Logpath`
	Filenames  fNames `json:Filenames`
}

type fNames struct {
	Errors string `json:Errors`
	Info   string `json:Info`
	Debug  string `json:Debug`
}

var conf = &Config{
	Filenames: fNames{Errors: "errors.log", Info: "info.log", Debug: "debug.log"},
}

func init() {

	filename, err := osext.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Error: could not get a path to binary file for creating logdir: %s\n", err.Error())
	}
	filepath := path.Dir(filename)
	fmt.Fprintf(os.Stderr, "[slflog] Configlog.json will be founded on %s directory\n", filepath)

	file, e := ioutil.ReadFile("configlog.json")
	if e != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Config logfile error: %s\n", e.Error())
	}

	if err := json.Unmarshal([]byte(file), &conf); err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Config logfile error: %s\n", err.Error())
	}
	fmt.Printf("Results: %v\n", conf)

	initLoggers(conf.Logpath, conf.ConsoleLvl)
}

// Init loggers
func initLoggers(logpath string, loglvl string) {

	var logHandlers []slog.EntryHandler

	// optionally define the format (this here is the default one)
	//bhInfo.SetTemplate("{{.Time}} [\033[{{.Color}}m{{.Level}}\033[0m] {{.Context}}{{if .Caller}} ({{.Caller}}){{end}}: {{.Message}}{{if .Error}} (\033[31merror: {{.Error}}\033[0m){{end}} {{.Fields}}")

	ConfigWriterOutput(&logHandlers, getLogLevel(loglvl), os.Stderr)

	err := setLogOutput(&logHandlers, logpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Error init loggers: %s\n", err.Error())
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

	conf.Logpath = pathForLogs

	ConfigFileOutput(logHandlers, slf.LevelDebug, filepath.Join(conf.Logpath, conf.Filenames.Debug))
	ConfigFileOutput(logHandlers, slf.LevelInfo, filepath.Join(conf.Logpath, conf.Filenames.Info))
	ConfigFileOutput(logHandlers, slf.LevelError, filepath.Join(conf.Logpath, conf.Filenames.Errors))

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
