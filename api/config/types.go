package config

type Config struct {
	Port          int
	PublicUIURL   string
	ConversionURL string
	DatabaseURL   string
	Search        SearchConfig
	Redis         RedisConfig
	S3            S3Config
	Limits        LimitsConfig
	Security      SecurityConfig
	SMTP          SMTPConfig
	Defaults      DefaultsConfig
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds int
	MultipartBodyLengthLimitMB    int
	FileProcessingMaxSizeMB       int
}

type DefaultsConfig struct {
	WorkspaceStorageCapacityBytes int64
}

type TokenConfig struct {
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	TokenAudience        string
	TokenIssuer          string
}

type SearchConfig struct {
	URL string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
}

type SecurityConfig struct {
	JWTSigningKey string
	CORSOrigins   []string
	APIKey        string
}

type SMTPConfig struct {
	Host          string
	Port          int
	Secure        bool
	Username      string
	Password      string
	SenderAddress string
	SenderName    string
}
