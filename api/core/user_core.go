package core

import "voltaserve/model"

type User struct {
	Id         string  `json:"id"`
	FullName   string  `json:"fullName"`
	Picture    *string `json:"picture,omitempty"`
	Email      string  `json:"email"`
	Username   string  `json:"username"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime"`
}

type userMapper struct {
}

func newUserMapper() *userMapper {
	return &userMapper{}
}

func (mp *userMapper) mapUser(user model.UserModel) *User {
	return &User{
		Id:         user.GetID(),
		FullName:   user.GetFullName(),
		Picture:    user.GetPicture(),
		Email:      user.GetEmail(),
		Username:   user.GetUsername(),
		CreateTime: user.GetCreateTime(),
		UpdateTime: user.GetUpdateTime(),
	}
}

func (mp *userMapper) mapUsers(users []model.UserModel) ([]*User, error) {
	res := []*User{}
	for _, u := range users {
		res = append(res, mp.mapUser(u))
	}
	return res, nil
}
