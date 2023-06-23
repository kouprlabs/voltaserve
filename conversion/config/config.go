package config

import (
	"os"
	"strconv"
)

var config *Config

func GetConfig() Config {
	if config == nil {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		config = &Config{
			Port: port,
		}
		readURLs(config)
		readSecurity(config)
		readS3(config)
		readLimits(config)
	}
	return *config
}

func readURLs(config *Config) {
	config.APIURL = os.Getenv("API_URL")
	config.LanguageURL = os.Getenv("LANGUAGE_URL")
	config.FFMPEGURL = os.Getenv("VOLTASERVE_FFMPEG_URL")
	config.ExiftoolURL = os.Getenv("VOLTASERVE_EXIFTOOL_URL")
	config.ImageMagickURL = os.Getenv("VOLTASERVE_IMAGEMAGICK_URL")
	config.LibreOfficeURL = os.Getenv("VOLTASERVE_LIBREOFFICE_URL")
	config.OCRMyPDFURL = os.Getenv("VOLTASERVE_OCRMYPDF_URL")
	config.PopplerURL = os.Getenv("VOLTASERVE_POPPLER_URL")
	config.TesseractURL = os.Getenv("VOLTASERVE_TESSERACT_URL")
}

func readSecurity(config *Config) {
	config.Security.APIKey = os.Getenv("SECURITY_API_KEY")
}

func readS3(config *Config) {
	config.S3.URL = os.Getenv("S3_URL")
	config.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
	config.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
	config.S3.Region = os.Getenv("S3_REGION")
	if len(os.Getenv("S3_SECURE")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("S3_SECURE"))
		if err != nil {
			panic(err)
		}
		config.S3.Secure = v
	}
}

func readLimits(config *Config) {
	if len(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ExternalCommandTimeoutSeconds = int(v)
	}
	if len(os.Getenv("LIMITS_FILE_PROCESSING_MAX_SIZE_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_FILE_PROCESSING_MAX_SIZE_MB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.FileProcessingMaxSizeMB = int(v)
	}
	if len(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_WIDTH")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_WIDTH"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ImagePreviewMaxWidth = int(v)
	}
	if len(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_HEIGHT")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_HEIGHT"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ImagePreviewMaxHeight = int(v)
	}
	if len(os.Getenv("LIMITS_LANGUAGE_SCORE_THRESHOLD_PERCENT")) > 0 {
		v, err := strconv.ParseFloat(os.Getenv("LIMITS_LANGUAGE_SCORE_THRESHOLD_PERCENT"), 64)
		if err != nil {
			panic(err)
		}
		config.Limits.LanguageScoreThreshold = v
	}
}
