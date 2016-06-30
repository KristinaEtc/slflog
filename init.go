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

/*const (
	errorFilename = "error.log"
	infoFilename  = "info.log"
	debugFilename = "debug.log"
)*/

// Struct for log config.
type Config struct {
	StderrLvl string            `json:StderrLvl`
	Logpath   string            `json:Logpath`
	Filenames map[string]string `json:Filenames`
}

var conf = &Config{
	Filenames: map[string]string{"ERROR": "errorrrrrs.log", "INFO": "info.log", "DEBUG": "debug.log"},
	StderrLvl: "DEBUG",
}

const configLogfile string = "configlog.json"

// Searching configuration log file.
// Parsing configuration on it. If file doesn't exist, use default settings.
func init() {

	filename, err := osext.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Error: could not get a path to binary file for creating logdir: %s\n", err.Error())
	}
	filepath := path.Dir(filename)
	fmt.Fprintf(os.Stderr, "[slflog] Configlog.json will be founded on %s directory\n", filepath)

	// Default logpath - aDirectoryWithBinaryFile/logs.
	fpath := filepath + string(os.PathSeparator) + "logs"
	conf.Logpath = fpath

	// Parsing configlog.json
	file, err := ioutil.ReadFile(configLogfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Config logfile error: %s\n", err.Error())
	}
	var userConfig = &Config{}
	err = json.Unmarshal(file, &userConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Config logfile error: %s\n", err.Error())
	}

	initLoggers(conf.Logpath, conf.StderrLvl)
}

// Init loggers: writers, log output, entry handlers.
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

	ConfigFileOutput(logHandlers, slf.LevelDebug, filepath.Join(conf.Logpath, conf.Filenames["DEBUG"]))
	ConfigFileOutput(logHandlers, slf.LevelInfo, filepath.Join(conf.Logpath, conf.Filenames["INFO"]))
	ConfigFileOutput(logHandlers, slf.LevelError, filepath.Join(conf.Logpath, conf.Filenames["ERROR"]))

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

// Exists returns whether the given file or directory exists or not.
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

// Format string to slf.Level.
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

// Fill in the blank fields of config structure with default values from confDefault.
func fillConfig(userConfig *Config) {

	if userConfig.Logpath == "" {
		userConfig.Logpath = conf.Logpath
	}
	if userConfig.StderrLvl == "" {
		userConfig.StderrLvl = conf.StderrLvl
	}
	if _, exist := userConfig.Filenames["ERROR"]; !exist {
		userConfig.Filenames["ERROR"] = conf.Filenames["ERROR"]
	}
	if _, exist := userConfig.Filenames["INFO"]; !exist {
		userConfig.Filenames["INFO"] = conf.Filenames["INFO"]
	}
	if _, exist := userConfig.Filenames["DEBUG"]; !exist {
		userConfig.Filenames["DEBUG"] = conf.Filenames["DEBUG"]
	}

	conf = userConfig
}
