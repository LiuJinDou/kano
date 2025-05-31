package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

const CONFIG_ENV_NAME = "GO_CONFIG_ENV"

type Configuration struct {
	APP struct {
		Name    string `mapstructure:"name"`
		AppEnv  string `mapstructure:"app_env"`
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
	Logger           *Logger    `mapstructure:"logger"`
	DBS              []DB       `mapstructure:"dbs"`
	TencentYun       TencentYun `mapstructure:"tencentyun"`
	ApplicationCodes []string   `mapstructure:"application_codes"` // List of allowed application codes
}

type DB struct {
	ConnectionName string `mapstructure:"connection_name"`
	Driver         string `mapstructure:"driver"`
	Host           string `mapstructure:"host"`
	Database       string `mapstructure:"database"`
	Port           int    `mapstructure:"port"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	Charset        string `mapstructure:"charset"`
	Debug          bool   `mapstructure:"debug"`
}

type TencentYun struct {
	AppId     string `mapstructure:"app_id"`
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Path      string `mapstructure:"path"`
	Region    string `mapstructure:"region"`
	Endpoint  string `mapstructure:"endpoint"`
}

// Config for logging
type Logger struct {
	// log level
	LogLevel int8 `mapstructure:"log_level"`

	// Enable console logging
	ConsoleLoggingEnabled bool `mapstructure:"console_logging_enabled"`

	// EncodeLogsAsJSON makes the log framework log JSON
	EncodeLogsAsJSON bool `mapstructure:"encode_logs_as_json"`

	// FileLoggingEnabled makes the framework log to a file, the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `mapstructure:"file_logging_enabled"`

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `mapstructure:"filename"`

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `mapstructure:"max_size"`

	// MaxBackups the max number of rolled files to keep
	MaxBackups int `mapstructure:"max_backups"`

	// MaxAge the max age in days to keep a logfile
	MaxAge int `mapstructure:"max_age"`

	RollingWrite io.Writer
}

// Config is the global configuration variable
var Config *Configuration

func LoadConfig() {
	// Load configuration from file or environment variables
	// This is a placeholder function. Actual implementation will depend on the configuration management strategy.
	var configFile string
	switch env := os.Getenv(CONFIG_ENV_NAME); env {
	case "production":
		// Load production configuration
		configFile = "config-prod"
	case "testing":
		// Load testing configuration
		configFile = "config-test"
	default:
		// Load default configuration
		configFile = "config-dev"
	}
	// Initialize viper to read the configuration file
	viper.SetConfigName(configFile)      // 文件名（不带后缀）
	viper.SetConfigType("yaml")          // 文件类型
	viper.AddConfigPath("../../configs") // 配置文件所在路径
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var cfg Configuration

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	Config = &cfg
	fmt.Printf("Loaded configuration: %+v\n", Config)
	// Set up the logger
	SetLoggerConfig()
}

// SetConfig set logger config
func SetLoggerConfig() {

	if Config.Logger == nil {
		Config.Logger = &Logger{
			LogLevel:              0, // Default log level
			ConsoleLoggingEnabled: true,
			EncodeLogsAsJSON:      false,
			FileLoggingEnabled:    false,
			Filename:              "",
		}
	}

	if Config.Logger.FileLoggingEnabled {
		if Config.Logger.Filename == "" {
			name := filepath.Base(os.Args[0]) + "-fox.log"
			Config.Logger.Filename = filepath.Join(os.TempDir(), name)
		}

		Config.Logger.RollingWrite = &lumberjack.Logger{
			Filename:   Config.Logger.Filename,
			MaxSize:    Config.Logger.MaxSize,
			MaxBackups: Config.Logger.MaxBackups,
			MaxAge:     Config.Logger.MaxAge,
		}
	}
}

func IsApplicationCodeAllowed(code string) bool {
	for _, v := range Config.ApplicationCodes {
		if v == code {
			return true
		}
	}
	return false
}
