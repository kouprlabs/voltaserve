package model

type Task interface {
	GetID() string
	GetDescription() string
	GetError() *string
	GetPercentage() *int
	GetIsComplete() bool
	GetIsIndeterminate() bool
	SetDescription(string)
	SetError(*string)
	SetPercentage(*int)
	SetIsComplete(bool)
	SetIsIndeterminate(bool)
}
