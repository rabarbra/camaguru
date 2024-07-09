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

func getImgs(db *sql.DB, limit int, offset int) ([]Img, error) {
	rows, err := db.Query(fmt.Sprintf(
		"SELECT id, link, created_at, user_if FROM imgs LIMIT=%d OFFSET=%d;",
		limit, offset,
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return nil, nil
	// cols, _ := rows.Columns()
	// row := make([]Img, len(cols))
	// imgs := make([]*Img, len(cols))
	// for i := range row {
	// 	imgs[i] = &row[i]
	// }
	// for rows.Next() {
	// 	err = rows.Scan < Img > (imgs)
	// 	if err != nil {
	// 		fmt.Println("cannot scan row:", err)
	// 	}
	// 	fmt.Println(row...)
	// }
	// return rows.Err()
}
