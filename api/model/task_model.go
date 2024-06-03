package model

type Task interface {
	GetID() string
	GetName() string
	GetError() *string
	GetPercentage() *int
	GetIsComplete() bool
	GetIsIndeterminate() bool
	GetUserID() string
	SetName(string)
	SetError(*string)
	SetPercentage(*int)
	SetIsComplete(bool)
	SetIsIndeterminate(bool)
	SetUserID(string)
	GetCreateTime() string
	GetUpdateTime() *string
}
