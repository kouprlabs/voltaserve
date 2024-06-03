package model

type Task interface {
	GetID() string
	GetName() string
	GetError() *string
	GetPercentage() *int
	GetIsIndeterminate() bool
	GetUserID() string
	HasError() bool
	SetName(string)
	SetError(*string)
	SetPercentage(*int)
	SetIsIndeterminate(bool)
	SetUserID(string)
	GetCreateTime() string
	GetUpdateTime() *string
}
