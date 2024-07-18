package main

import (
	"net/http"
	"strings"
)

type BaseModel interface {
	NewItem(r *http.Request, userId int64) error
}

type User struct {
	Id            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Pass          string `json:"pass"`
	Avatar        string `json:"avatar"`
	EmailVerified bool   `json:"email_verified"`
	LikeNotify    bool   `json:"like_notify"`
	CommNotify    bool   `json:"comm_notify"`
}

type Img struct {
	Id        int64  `json:"id"`
	Link      string `json:"link"`
	UserId    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (i *Img) NewItem(r *http.Request, userId int64) error {
	filePath, err := uploadImg(r)
	if err != nil {
		return err
	}
	i.Link = strings.TrimPrefix(filePath, "assets")
	i.UserId = userId
	return nil
}

type Comment struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UserId    int64  `json:"user_id"`
	ImgId     int64  `json:"img_id"`
}

type Like struct {
	Id     int64 `json:"id"`
	UserId int64 `json:"user_id"`
	ImgId  int64 `json:"img_id"`
}
