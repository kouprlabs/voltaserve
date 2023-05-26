package model

type User interface {
	GetID() string
	GetFullName() string
	GetUsername() string
	GetEmail() string
	GetPicture() *string
	GetIsEmailConfirmed() bool
	GetCreateTime() string
	GetUpdateTime() *string
}
