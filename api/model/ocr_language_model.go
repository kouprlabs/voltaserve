package model

type OCRLanguage interface {
	GetID() string
	GetISO639Pt3() string
	GetName() string
}
