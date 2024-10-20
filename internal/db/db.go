// Пакет для работы с базой данных
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/iliamikado/UrlShortener/internal/logger"
)

// URLShortenerDB - структура работающая с БД
type URLShortenerDB struct {
	db *sql.DB
}

// URLDB - глобальная переменная для обращения к базе данных.
var URLDB URLShortenerDB

// Функция для подключения к базе данных.
func Initialize(host string) {
	logger.Log.Info("Opening DB with host=" + host)
	db, err := sql.Open("pgx", host)
	if err != nil {
		logger.Log.Error("Failed to open DB")
		panic(err)
	}
	URLDB = URLShortenerDB{db}
}

// Функция для создания таблицы сокращенных ссылок в БД.
func (urlDB *URLShortenerDB) CreateURLTable() {
	urlDB.db.Exec(`create table if not exists urls (
		uuid text PRIMARY KEY NOT NULL,
		long_url text UNIQUE NOT NULL,
		user_id text,
		is_deleted boolean default false
	);`)
}

// Функция для добавления ссылки в БД. Возвращает id.
func (urlDB *URLShortenerDB) AddURL(longURL string, userID string, getID func() string) (string, error) {
	var err error
	var id string
	for {
		id = getID()
		_, err = urlDB.db.Exec(`insert into urls values ($1, $2, $3)`, id, longURL, userID)
		if err != nil && isErrorWithID(err) {
			continue
		}
		break
	}
	return id, err
}

// Функция для добавления нескольких ссылок в БД. Возвращает массив из id.
func (urlDB *URLShortenerDB) AddManyURLs(longURLs []string, userID string, getID func() string) ([]string, error) {
	tx, _ := urlDB.db.Begin()
	stmt, _ := tx.Prepare(`insert into urls values ($1, $2, $3)`)
	defer stmt.Close()
	ids := make([]string, len(longURLs))
	for i := 0; i < len(longURLs); i++ {
		var err error
		for {
			ids[i] = getID()
			_, err = stmt.Exec(ids[i], longURLs[i], userID)
			if err != nil && isErrorWithID(err) {
				continue
			}
			break
		}
		if err != nil {
			return nil, tx.Rollback()
		}
	}
	return ids, tx.Commit()
}

// GetURL - функция возвращающая URL по id
func (urlDB *URLShortenerDB) GetURL(id string) (string, bool, error) {
	row := urlDB.db.QueryRow(`select long_url, is_deleted from urls where uuid = $1`, id)
	var longURL string
	var isDeleted bool
	err := row.Scan(&longURL, &isDeleted)
	return longURL, isDeleted, err
}

// GetIDByURL - функция возвращающая id по URL
func (urlDB *URLShortenerDB) GetIDByURL(longURL string) (string, error) {
	row := urlDB.db.QueryRow(`select uuid from urls where long_url = $1`, longURL)
	var id string
	err := row.Scan(&id)
	return id, err
}

// GetUserURLs - функция возвращающая все URL пользователя
func (urlDB *URLShortenerDB) GetUserURLs(userID string) [][2]string {
	rows, _ := urlDB.db.Query("select uuid, long_url from urls where user_id = $1 and not is_deleted", userID)
	defer rows.Close()
	var ans [][2]string
	for rows.Next() {
		var uuid, longURL string
		rows.Scan(&uuid, &longURL)
		ans = append(ans, [2]string{uuid, longURL})
	}
	if rows.Err() != nil {
		panic("error in reading db")
	}
	return ans
}

// DeleteURLs - функция удаляющая URL по id у пользователя
func (urlDB *URLShortenerDB) DeleteURLs(ids []string, userID string) {
	tx, _ := urlDB.db.Begin()
	valueStrings := make([]string, 0, len(ids))
	variables := make([]any, 0, len(ids)+1)
	variables = append(variables, userID)
	for i, id := range ids {
		valueStrings = append(valueStrings, fmt.Sprintf("$%d", i+2))
		variables = append(variables, id)
	}
	query := fmt.Sprintf("update urls set is_deleted = true where user_id = $1 and uuid in (%s)", strings.Join(valueStrings, ","))
	_, err := urlDB.db.Exec(query, variables...)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

// Close - функция закрывающая связь с БД
func (urlDB *URLShortenerDB) Close() {
	urlDB.db.Close()
}

// Ping - проверка подключения к бд
func (urlDB *URLShortenerDB) Ping() error {
	return urlDB.db.PingContext(context.TODO())
}

func isErrorWithID(e error) bool {
	var err *pgconn.PgError
	if errors.As(e, &err) {
		if err.Code == pgerrcode.UniqueViolation && err.ConstraintName == "urls_uuid_key" {
			return true
		}
	}
	return false
}
