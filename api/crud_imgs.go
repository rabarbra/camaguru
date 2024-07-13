package main

import (
	"database/sql"
	"fmt"
)

func createImg(db *sql.DB, img Img) error {
	res, err := db.Exec(fmt.Sprintf(`
		INSERT INTO imgs(link, user_id)
		VALUES('%s', %d)`,
		img.Link, img.UserId,
	))
	fmt.Println(res, err)
	return err
}

func updateImg(db *sql.DB, img Img) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE imgs SET link='%s WHERE id=%d",
		img.Link, img.Id,
	))
	return err
}

func getImgById(db *sql.DB, id int64) (Img, error) {
	var img Img
	err := db.QueryRow(fmt.Sprintf(
		"SELECT id, link, created_at, user_if FROM imgs WHERE id=%d;",
		id,
	)).Scan(&img.Id, &img.Link, &img.CreatedAt, &img.UserId)
	return img, err
}

func getImgs(db *sql.DB, userId int64, limit int, offset int) ([]Img, error) {
	rows, err := db.Query(fmt.Sprintf(
		"SELECT id, link, created_at, user_id FROM imgs WHERE user_id=%d LIMIT %d OFFSET %d;",
		userId, limit, offset,
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var imgs []Img
	for rows.Next() {
		var img Img
		if err := rows.Scan(&img.Id, &img.Link, &img.CreatedAt, &img.UserId); err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return imgs, nil
}
