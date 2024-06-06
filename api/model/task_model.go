package model

const (
	TaskStatusWaiting = "waiting"
	TaskStatusRunning = "running"
	TaskStatusSuccess = "success"
	TaskStatusError   = "error"
)

type Task interface {
	GetID() string
	GetName() string
	GetError() *string
	GetPercentage() *int
	GetIsIndeterminate() bool
	GetUserID() string
	GetStatus() string
	GetPayload() map[string]string
	HasError() bool
	SetName(string)
	SetError(*string)
	SetPercentage(*int)
	SetIsIndeterminate(bool)
	SetUserID(string)
	SetStatus(string)
	SetPayload(map[string]string)
	GetCreateTime() string
	GetUpdateTime() *string
}
