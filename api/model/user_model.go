package model

type UserModel interface {
	GetId() string
	GetFullName() string
	GetUsername() string
	GetEmail() string
	GetPicture() *string
	GetIsEmailConfirmed() bool
	GetCreateTime() string
	GetUpdateTime() *string
}
