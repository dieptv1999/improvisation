package repository

import (
	"context"
	"fmt"
	pb "git.local/go-app/models"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

func (db *PostgreSQLDB) ReadSong(id int) (*pb.Song, error) {
	var title, link, pic, author, lrc, url, thumbnail string
	var createdAt time.Time
	err := db.conn.QueryRow(context.Background(), `SELECT link, title, author, lrc, url, pic, thumbnail, created_at FROM song WHERE id=$1`, id).Scan(&link, &title, &author, &lrc, &url, &pic, &thumbnail, &createdAt)
	if err != nil && err != pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pb.Song{Id: id, Link: link, Title: title, Author: author, Lrc: lrc, Url: url, Pic: pic, Thumbnail: thumbnail, CreatedAt: createdAt}, nil
}

func (db *PostgreSQLDB) InsertSong(Song *pb.Song) error {
	_, err := db.conn.Exec(context.Background(), `
			INSERT INTO song(id, link, title, author, lrc, url, pic, thumbnail)
			VALUES(?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE created_at=VALUES(created_at)
		`, Song.GetId(), Song.GetLink(), Song.GetTitle(), Song.GetAuthor(), Song.GetLrc(), Song.GetUrl(), Song.GetPic(), Song.GetThumbnail())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update task: %v\n", err)
		return err
	}
	return nil
}

func (db *PostgreSQLDB) ListSong(query string) ([]*pb.Song, error) {
	rows, _ := db.conn.Query(context.Background(), "select id, link, title, author, lrc, url, pic, thumbnail, created_at from song where LOWER(title) like LOWER('%"+query+"%')")

	var Songs []*pb.Song

	for rows.Next() {
		var title, link, pic, author, lrc, url, thumbnail string
		var createdAt time.Time
		var id int
		err := rows.Scan(&id, &link, &title, &author, &lrc, &url, &pic, &thumbnail, &createdAt)
		if err != nil {
			return nil, err
		}
		Songs = append(Songs, &pb.Song{Id: id, Link: link, Title: title, Author: author, Lrc: lrc, Url: url, Pic: pic, Thumbnail: thumbnail, CreatedAt: createdAt})
	}

	return Songs, nil
}

func (db *PostgreSQLDB) ListPopular() ([]*pb.Song, error) {
	rows, _ := db.conn.Query(context.Background(), "select link, title, author, lrc, url, pic, thumbnail, created_at from song inner join popular p on song.id = p.song_id")

	var Songs []*pb.Song

	for rows.Next() {
		var title, link, pic, author, lrc, url, thumbnail string
		var createdAt time.Time
		err := rows.Scan(&link, &title, &author, &lrc, &url, &pic, &thumbnail, &createdAt)
		if err != nil {
			return nil, err
		}
		Songs = append(Songs, &pb.Song{Link: link, Title: title, Author: author, Lrc: lrc, Url: url, Pic: pic, Thumbnail: thumbnail, CreatedAt: createdAt})
	}

	return Songs, nil
}
