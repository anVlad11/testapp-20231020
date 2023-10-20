package config

type App struct {
	HTTPServer HTTPServer `mapstructure:"http_server"`
	Tracing    Tracing    `mapstructure:"tracing"`
	Metrics    HTTPServer `mapstructure:"metrics"`
}

type HTTPServer struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
	Listen  bool   `mapstructure:"listen"`
}

type Tracing struct {
	Enabled     bool   `mapstructure:"enabled"`
	ProviderURL string `mapstructure:"provider_url"`
}
