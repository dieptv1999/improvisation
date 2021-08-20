package repository

import (
	"context"
	"fmt"
	pb "git.local/go-app/models"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

func (db *PostgreSQLDB) ReadCategory(id int) (*pb.Category, []*pb.Song, error) {
	var name, link, pic, desc, thumbnail string
	var createdAt time.Time
	var _type int
	err := db.conn.QueryRow(context.Background(), `SELECT link, name, c.desc, type, pic, thumbnail, created_at FROM category as c WHERE id=$1`, id).Scan(&link, &name, &desc, &_type, &pic, &thumbnail, &createdAt)
	if err != nil && err != pgx.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	//song of album
	rows, _ := db.conn.Query(context.Background(), "select s.link, s.title, author, lrc, s.url, s.pic, s.thumbnail, s.created_at from song as s inner join category c on s.category_id = $1", id)

	fmt.Println(rows, "rows")

	var Songs []*pb.Song

	for rows.Next() {
		var title, link, pic, author, lrc, url, thumbnail string
		var createdAt time.Time
		err := rows.Scan(&link, &title, &author, &lrc, &url, &pic, &thumbnail, &createdAt)
		if err != nil {
			return nil, nil, err
		}
		Songs = append(Songs, &pb.Song{Link: link, Title: title, Author: author, Lrc: lrc, Url: url, Pic: pic, Thumbnail: thumbnail, CreatedAt: createdAt})
	}
	return &pb.Category{Id: id, Link: link, Name: name, Pic: pic, Thumbnail: thumbnail, Desc: desc, Type: _type, CreatedAt: createdAt}, Songs, nil
}

func (db *PostgreSQLDB) InsertCategory(Category *pb.Category) error {
	_, err := db.conn.Exec(context.Background(), `
			INSERT INTO category(id, link, name, desc, type, pic, thumbnail)
			VALUES(?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE created_at=VALUES(created_at)
		`, Category.GetId(), Category.GetLink(), Category.GetName(), Category.GetDesc(), Category.GetType(), Category.GetPic(), Category.GetThumbnail())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update task: %v\n", err)
		return err
	}
	return nil
}

func (db *PostgreSQLDB) ListCategory() ([]*pb.Category, error) {
	rows, _ := db.conn.Query(context.Background(), "select t.id, t.link, t.name, t.desc, t.type, t.pic, t.thumbnail, t.created_at from public.category t")

	var Categories []*pb.Category

	for rows.Next() {
		var name, link, pic, desc, thumbnail string
		var createdAt time.Time
		var id, _type int
		err := rows.Scan(&id, &link, &name, &desc, &_type, &pic, &thumbnail, &createdAt)
		if err != nil {
			return nil, err
		}
		Categories = append(Categories, &pb.Category{Id: id, Link: link, Name: name, Pic: pic, Thumbnail: thumbnail, Desc: desc, Type: _type, CreatedAt: createdAt})
	}

	return Categories, nil
}
