package main

type User struct {
	Id            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Pass          string `json:"pass"`
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
