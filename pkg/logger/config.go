package logger

type Config struct {
	TimeFormat string
	Verbose    bool `mapstructure:"verbose"`
	Pretty     bool `mapstructure:"pretty"`
}
