package repository

import (
	"context"
	"fmt"
	"os"

	pb "git.local/go-app/models"
	"github.com/jackc/pgx/v4"
)

func (db *PostgreSQLDB) ReadAccount(id int) (*pb.Account, error) {
	var photoURL, first_name, last_name, user_email, user_display_name string
	var createdAt int64
	var type_login int
	err := db.conn.QueryRow(context.Background(), `SELECT photoURL, first_name, last_name, user_email, username, user_display_name, type_login, created_at FROM account WHERE id=$1`,
		id).Scan(&photoURL, &first_name, &last_name, &user_email, &user_display_name, &type_login, &createdAt)
	if err != nil && err != pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pb.Account{Id: id, PhotoURL: photoURL, FirstName: first_name, LastName: last_name, UserEmail: user_email, UserDisplayName: user_display_name, TypeLogin: type_login, CreatedAt: createdAt}, nil
}

func (db *PostgreSQLDB) InsertAccount(account *pb.Account) error {
	_, err := db.conn.Exec(context.Background(), `
			INSERT INTO account(photoURL, first_name, last_name, user_email, username, user_display_name, type_login)
			VALUES(?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE created_at=VALUES(created_at)
		`, account.PhotoURL, account.FirstName, account.LastName, account.UserEmail, account.Username, account.UserDisplayName, account.TypeLogin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update task: %v\n", err)
		return err
	}
	return nil
}

func (db *PostgreSQLDB) ListAccount() ([]*pb.Account, error) {
	rows, _ := db.conn.Query(context.Background(), "select * from account")

	var accounts []*pb.Account

	for rows.Next() {
		var photoURL, first_name, last_name, user_email, user_display_name string
		var createdAt int64
		var type_login, id int
		err := rows.Scan(&id, &photoURL, &first_name, &last_name, &user_email, &user_display_name, &type_login, &createdAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &pb.Account{Id: id, PhotoURL: photoURL, FirstName: first_name, LastName: last_name, UserEmail: user_email, UserDisplayName: user_display_name, TypeLogin: type_login, CreatedAt: createdAt})
	}

	return accounts, nil
}
