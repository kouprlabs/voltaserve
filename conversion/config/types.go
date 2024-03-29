package config

type Config struct {
	Port     int
	APIURL   string
	Security SecurityConfig
	Limits   LimitsConfig
	S3       S3Config
}

type SecurityConfig struct {
	APIKey string `json:"api_key"`
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds int
	FileProcessingMaxSizeMB       int
	ImagePreviewMaxWidth          int
	ImagePreviewMaxHeight         int
	LanguageScoreThreshold        float64
}

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
}
