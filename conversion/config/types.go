package config

type Config struct {
	ConversionURL string         `json:"conversion_url"`
	APIURL        string         `json:"api_url"`
	Security      SecurityConfig `json:"security"`
	Limits        LimitsConfig   `json:"limits"`
	S3            S3Config       `json:"s3"`
}

type SecurityConfig struct {
	APIKey string `json:"api_key"`
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds int `json:"external_command_timeout_seconds"`
	FileProcessingMaxSizeMB       int `json:"file_processing_max_size_mb"`
	ImagePreviewMaxWidth          int `json:"image_preview_max_width"`
	ImagePreviewMaxHeight         int `json:"image_preview_max_height"`
}

type S3Config struct {
	URL       string `json:"url"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Secure    bool   `json:"secure"`
}
