package config

type Config struct {
	Port     int
	Security SecurityConfig
	Limits   LimitsConfig
}

type SecurityConfig struct {
	APIKey string `json:"api_key"`
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds int
	MultipartBodyLengthLimitMB    int
}
