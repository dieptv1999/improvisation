package repository

import (
	"context"
	"fmt"
	pb "git.local/go-app/models"
	"github.com/jackc/pgx/v4"
	"os"
)

func (db *PostgreSQLDB) ReadAccount(id int) (*pb.Account, error) {
	var name, email, pic string
	var dateOfBirth, createdAt int64
	err := db.conn.QueryRow(context.Background(), `SELECT name, email, date_of_birth, pic, created_at FROM account WHERE id=$1`, id).Scan(&name, &createdAt, &email, &pic, &dateOfBirth)
	if err != nil && err != pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pb.Account{Id: id, Name: name, Email: email, CreatedAt: createdAt, DateOfBirth: dateOfBirth}, nil
}

func (db *PostgreSQLDB) InsertAccount(Account *pb.Account) error {
	_,err := db.conn.Exec(context.Background(),`
			INSERT INTO account(id, name, email, date_of_birth, pic)
			VALUES(?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE created_at=VALUES(created_at)
		`, Account.GetId(), Account.GetName(), Account.GetEmail(), Account.GetDateOfBirth(), Account.GetPic())

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
		var name, email, pic string
		var dateOfBirth, createdAt int64
		err := rows.Scan(&name, &createdAt, &email, &pic, &dateOfBirth)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &pb.Account{Name: name, Email: email, CreatedAt: createdAt, DateOfBirth: dateOfBirth})
	}

	return accounts, nil
}
