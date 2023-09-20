package DbConnectionUtils

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/Database"
)

func EstablishDbConnection(dbFactory Database.DatabaseFactory) (*sql.DB, error) {
	connection, err := dbFactory.NewDbConnection()
	if err != nil {
		return nil, err
	}

	err = connection.Ping()
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func CheckIfEmailIsTaken(email string, connection *sql.DB) (int, error) {
	query := fmt.Sprintf("select id from user_table where email = '%s'", email)

	var id int

	err := connection.QueryRow(query).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	if id > 0 {
		return id, nil
	} else {
		return 0, nil
	}
}

func GetUserPassword(id int, connection *sql.DB) (string, error) {
	query := fmt.Sprintf("select password from user_table where id = '%d'", id)

	var pass string

	err := connection.QueryRow(query).Scan(&pass)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	if len(pass) > 0 {
		return pass, nil
	} else {
		return "", errors.New("empty password in db")
	}
}
