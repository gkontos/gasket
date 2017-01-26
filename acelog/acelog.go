package acelog

import (
	"net/http"
	"os"
	"sync"
	"time"

	logrus "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO "github.com/lestrrat/go-apache-logformat"

// TODO change log file configuration for test, dev, prod
// TODO use environment variable for config file location, ie os.Getenv("PROJECT_CFG")
var (
	logger      *logrus.Logger
	once        sync.Once
	accessLog   *logrus.Logger
	environment string
)

type LogConfig struct {
	format     string
	output     string
	level      string
	filename   string
	maxsize    int // megabytes
	maxbackups int
	maxage     int //days
}
type aceLoggingHandler struct {
	handler      http.Handler
	accessLogger *logrus.Logger
}

func init() {
	// pull configuration from config/log.config by default
	logger = initLogger("applicationLog")

}

func initLogger(configPrefix string) *logrus.Logger {

	logrus.Info("setting up logger for ", configPrefix)

	v := viper.New()
	v.SetConfigName("log") // no need to include file extension
	if _, err := os.Stat(os.Getenv("ACES_CFG")); err == nil {
		v.AddConfigPath(os.Getenv("ACES_CFG"))
	} else {
		v.AddConfigPath("./config") // set the path of your config file
		v.AddConfigPath("./../config")
	}
	err := v.ReadInConfig()
	var config *LogConfig
	if err != nil {
		config = &LogConfig{
			format: "text",
			output: "stdout",
			level:  "debug",
		}
		logrus.Warn("Unable to parse config file.  Setting logger to default configuration.")
		logrus.Warn("Using config file ", v.ConfigFileUsed())
		logrus.Error(err)
	} else {
		config = setConfiguration(v, configPrefix)
	}
	l := logrus.New()

	Configure(l, config)

	return l
}

// set the logger instance configuration
func Configure(l *logrus.Logger, config *LogConfig) {
	switch config.format {
	case "text":
		l.Formatter = &logrus.TextFormatter{}
	case "json":
		l.Formatter = &logrus.JSONFormatter{}
	default:
		l.Formatter = &logrus.TextFormatter{}
	}
	switch config.output {
	case "stdout":
		l.Out = os.Stdout
	case "file":
		l.Warn("Setting output file to ", config.filename)
		l.Out = &lumberjack.Logger{
			Filename:   config.filename,
			MaxSize:    config.maxsize, // megabytes
			MaxBackups: config.maxbackups,
			MaxAge:     config.maxage, //days
		}
	case "stderr":
		l.Out = os.Stderr
	default:
		l.Out = os.Stdout
	}

	switch config.level {
	case "debug":
		l.Level = logrus.DebugLevel
	case "info":
		l.Level = logrus.InfoLevel
	case "error":
		l.Level = logrus.ErrorLevel
	case "warn":
		l.Level = logrus.WarnLevel
	case "fatal":
		l.Level = logrus.FatalLevel
	default:
		// Only log the warning severity or above.
		l.Level = logrus.DebugLevel
	}
}

func setConfiguration(v *viper.Viper, prefix string) *LogConfig {
	config := &LogConfig{
		format:     v.GetString(prefix + ".format"),
		output:     v.GetString(prefix + ".output"),
		level:      v.GetString(prefix + ".level"),
		filename:   v.GetString(prefix + ".filename"),
		maxsize:    v.GetInt(prefix + ".maxsize"),
		maxbackups: v.GetInt(prefix + ".maxbackups"),
		maxage:     v.GetInt(prefix + ".maxage"),
	}

	return config
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func RequestLogHandler(h http.Handler) http.Handler {

	once.Do(func() {
		accessLog = initLogger("accessLog")

	})
	return aceLoggingHandler{handler: h, accessLogger: accessLog}
}

// TODO add httpstatus, add response size
func (h aceLoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	h.handler.ServeHTTP(w, r)

	latency := time.Since(start)

	//		status := fn.status
	//		if status == 0 {
	//			status = 200
	//		}

	fields := logrus.Fields{
		//			"status":     status,
		"method":   r.Method,
		"request":  r.RequestURI,
		"remote":   r.RemoteAddr,
		"duration": float64(latency.Nanoseconds()) / float64(1000),
		//			"size":       fn.size,
		"referer":    r.Referer(),
		"user-agent": r.UserAgent(),
	}

	h.accessLogger.WithFields(fields).Info("Request Complete")

}
