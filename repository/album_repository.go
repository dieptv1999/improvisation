package repository

import (
	"context"
	pb "git.local/go-app/models"
	"github.com/jackc/pgx/v4"
	"time"
)

func (db *PostgreSQLDB) ReadAlbum(id int) (*pb.Album, []*pb.Song, error) {
	var pic, name, description string
	err := db.conn.QueryRow(context.Background(), `SELECT pic, name, description FROM album WHERE id=$1`, id).Scan(&pic, &name, &description)
	if err != nil && err != pgx.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	//song of album
	rows, _ := db.conn.Query(context.Background(), "select link, title, author, lrc, url, pic, thumbnail, created_at from song inner join popular p on song.id = p.song_id where album_id = $1", id)

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
	return &pb.Album{Id: id, Name: name, Description: description, Pic: pic}, Songs, nil
}

func (db *PostgreSQLDB) ListAlbum() ([]*pb.Album, error) {
	rows, _ := db.conn.Query(context.Background(), "SELECT id, pic, name, description FROM album")

	var Albums []*pb.Album

	for rows.Next() {
		var pic, name, description string
		var id int
		err := rows.Scan(&id, &pic, &name, &description)
		if err != nil {
			return nil, err
		}
		Albums = append(Albums, &pb.Album{Id: id, Name: name, Description: description, Pic: pic})
	}

	return Albums, nil
}
