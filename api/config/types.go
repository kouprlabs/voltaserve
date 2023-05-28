package config

type Config struct {
	APIURL        string         `json:"api_url"`
	UIURL         string         `json:"ui_url"`
	ConversionURL string         `json:"conversion_url"`
	DatabaseURL   string         `json:"database_url"`
	Search        SearchConfig   `json:"search"`
	Redis         RedisConfig    `json:"redis"`
	S3            S3Config       `json:"s3"`
	Limits        LimitsConfig   `json:"limits"`
	Security      SecurityConfig `json:"security"`
	SMTP          SMTPConfig     `json:"smtp"`
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds int `json:"external_command_timeout_seconds"`
	MultipartBodyLengthLimitMB    int `json:"multipart_body_length_limit_mb"`
}

type TokenConfig struct {
	AccessTokenLifetime  int    `json:"access_token_lifetime"`
	RefreshTokenLifetime int    `json:"refresh_token_lifetime"`
	TokenAudience        string `json:"token_audience"`
	TokenIssuer          string `json:"token_issuer"`
}

type SearchConfig struct {
	URL string `json:"url"`
}

type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type S3Config struct {
	URL       string `json:"url"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Secure    bool   `json:"secure"`
}

type SecurityConfig struct {
	JWTSigningKey string   `json:"jwt_signing_key"`
	CORSOrigins   []string `json:"cors_origins"`
	APIKey        string   `json:"api_key"`
}

type SMTPConfig struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Secure        bool   `json:"secure"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	SenderAddress string `json:"sender_address"`
	SenderName    string `json:"sender_name"`
}
