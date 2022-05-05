package logger

type Config struct {
	Verbose bool `mapstructure:"verbose"`
	Pretty  bool `mapstructure:"pretty"`
}
