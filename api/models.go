package main

import (
	"net/http"
	"orm"
	"strings"
	"time"
)

type User struct {
	orm.BaseModel
	Username      string `json:"username"`
	Email         string `json:"email"`
	Pass          string `json:"pass"`
	Avatar        string `json:"avatar"`
	EmailVerified bool   `json:"email_verified"`
	LikeNotify    bool   `json:"like_notify"`
	CommNotify    bool   `json:"comm_notify"`
}

func (*User) TableName() string {
	return "users"
}

func (u *User) NewItem(r *http.Request, userId int64) error {
	return nil
}

type Img struct {
	orm.BaseModel
	Link      string `json:"link"`
	UserId    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (*Img) TableName() string {
	return "imgs"
}

func (i *Img) NewItem(r *http.Request, userId int64) error {
	filePath, err := uploadImg(r)
	if err != nil {
		return err
	}
	i.Link = strings.TrimPrefix(filePath, "assets")
	i.UserId = userId
	i.CreatedAt = time.Now().Format(time.RFC3339)
	return nil
}

type Comment struct {
	orm.BaseModel
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UserId    int64  `json:"user_id"`
	ImgId     int64  `json:"img_id"`
}

func (*Comment) TableName() string {
	return "comments"
}

type Like struct {
	orm.BaseModel
	UserId int64 `json:"user_id"`
	ImgId  int64 `json:"img_id"`
}

func (*Like) TableName() string {
	return "likes"
}
