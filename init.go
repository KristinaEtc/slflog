package slflog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KristinaEtc/config"
	"github.com/kardianos/osext"
	"github.com/ventu-io/slf"
	"github.com/ventu-io/slog"
)

var (
	bhDebug, bhInfo, bhError, bhDebugConsole, bhStdError *Handler
	logfileInfo, logfileDebug, logfileError              *os.File
	lf                                                   slog.LogFactory
	log                                                  slf.StructuredLogger
)

// Struct for log config.
type Config struct {
	StderrLvl string
	Logpath   string
	Filenames map[string]string
}

// ConfFile is a file with all program options
type ConfFile struct {
	Logs Config
}

var logConfig = ConfFile{
	Logs: Config{
		Filenames: map[string]string{"ERROR": "errors.log", "INFO": "info.log", "DEBUG": "debug.log"},
		StderrLvl: "DEBUG",
		Logpath:   "",
	},
}

//var configLogFile string = ""

// Searching configuration log file.
// Parsing configuration on it. If file doesn't exist, use default settings.
func init() {

	//var cf ConfFile = defaultConf
	config.ReadGlobalConfig(&logConfig, "Logs")
	initLoggers(logConfig.Logs)
}

// Init loggers: writers, log output, entry handlers.
func initLoggers(logC Config) {

	var logHandlers []slog.EntryHandler

	ConfigWriterOutput(&logHandlers, getLogLevel(logC.StderrLvl), os.Stderr)

	err := setLogOutput(&logHandlers, logC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Error init loggers: %s\n", err.Error())
	}

	lf = slog.New()

	lf.SetLevel(slf.LevelDebug)

	if config.CallerInfo == true {
		lf.SetCallerInfo(slf.CallerShort)
	}

	lf.SetEntryHandlers(logHandlers...)
	slf.Set(lf)
}

func setLogOutput(logHandlers *[]slog.EntryHandler, logC Config) error {

	pathForLogs, err := getPathForLogDir(logC.Logpath)
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

	logC.Logpath = pathForLogs

	ConfigFileOutput(logHandlers, slf.LevelDebug, filepath.Join(logC.Logpath, logC.Filenames["DEBUG"]))
	ConfigFileOutput(logHandlers, slf.LevelInfo, filepath.Join(logC.Logpath, logC.Filenames["INFO"]))
	ConfigFileOutput(logHandlers, slf.LevelError, filepath.Join(logC.Logpath, logC.Filenames["ERROR"]))

	return nil
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

//----------------------------------------------------------------------------------------//
// Common utils for logger. Do not move it to another library (like KristinaEtc/utils c:),
// because logger _must_ initialize first.
//----------------------------------------------------------------------------------------//

func getPathForLogDir(logpath string) (string, error) {

	if filepath.IsAbs(logpath) == true {
		return logpath, nil
	}
	filename, err := osext.Executable()
	if err != nil {
		return "", err
	}

	fpath := filepath.Dir(filename)
	fpath = filepath.Join(fpath, logpath)
	return fpath, nil

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

/*
// GetGlobalConf unmarshal json-object cf
// If parsing was not successuful, function return a structure with default options
func getFromGlobalConf(cf interface{}, defaultVal interface{}, whatParsed string) {

	file, e := ioutil.ReadFile(getConfigFilename())
	if e != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Error: %s\n", e.Error())
	}

	if err := json.Unmarshal([]byte(file), cf); err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Error parsing JSON : [%s]. For %s will be used defaulf options.\n",
			whatParsed, err.Error())
		cf = defaultVal
	} else {
		fmt.Fprintf(os.Stderr, "[slflog] Parsed [%s] configuration from [%s] file.\n", whatParsed, getConfigFilename())
	}
	//log.Debugf("%v", cf)
}

// GetConfigFilename is a function fot getting a name of a binary with full path to it
func getConfigFilename() string {
	binaryPath, err := osext.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[slflog] Error: could not get a path to binary file: %s\n", err.Error())
	}
	if runtime.GOOS == "windows" {
		// without ".exe"
		binaryPath = binaryPath[:len(binaryPath)-4]
		fmt.Fprintf(os.Stderr, "[slflog] Configfile for windows")
	}

	return binaryPath + ".config"
}

*/

/*
// Fill in the blank fields of config structure with default values from confDefault.
func fillConfig(userConfig *Config) {

	if userConfig.Logpath == "" {
		userConfig.Logpath = conf.Logpath
	} else if userConfig.StderrLvl == "" {
		userConfig.StderrLvl = conf.StderrLvl
	} else if _, exist := userConfig.Filenames["ERROR"]; !exist {
		userConfig.Filenames["ERROR"] = conf.Filenames["ERROR"]
	} else if _, exist := userConfig.Filenames["INFO"]; !exist {
		userConfig.Filenames["INFO"] = conf.Filenames["INFO"]
	} else if _, exist := userConfig.Filenames["DEBUG"]; !exist {
		userConfig.Filenames["DEBUG"] = conf.Filenames["DEBUG"]
	} else {
		log.WithCaller(slf.CallerShort).Warnf("Wrong config level: %s", userConfig.Logpath)
	}
	conf = userConfig
}*/
