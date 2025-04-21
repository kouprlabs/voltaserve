package model

type StorageQuota interface {
	GetID() string
	GetUserID() string
	GetStorageCapacity() int64
	GetCreateTime() string
	GetUpdateTime() *string
	SetID(string)
	SetUserID(string)
	SetStorageCapacity(int64)
	SetCreateTime(string)
	SetUpdateTime(*string)
}
