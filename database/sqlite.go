package sqlite

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func openConnection() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "file:./database/sqlite.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Init() (err error) {

	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	var sql string

	// Search :: Table
	sql = `
    CREATE TABLE IF NOT EXISTS search (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        query TEXT NOT NULL DEFAULT '',
        page INT NOT NULL DEFAULT '0',
        results TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Search :: Index idx_query_page
	sql = `
    CREATE UNIQUE INDEX IF NOT EXISTS idx_query_page
	ON search (query,page)
    `
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Games :: Table
	sql = `
    CREATE TABLE IF NOT EXISTS games (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        guid TEXT,
        results TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Games : Index idx_guid
	sql = `
    CREATE UNIQUE INDEX IF NOT EXISTS idx_guid
	ON games (guid)
    `
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func AddSearchResults(q string, page int, data string) (id int, err error) {
	db, err := openConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sql := `INSERT INTO search (query, page, results) VALUES (?, ?, ?)`

	result, err := db.Exec(sql, q, page, data)
	if err != nil {
		return 0, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id64), nil
}

func DeleteSearchResults(id int) (err error) {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	sql := `DELETE FROM search WHERE id = ?`

	_, err = db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}

type SearchResult struct {
	Id      int
	Query   string
	Page    int
	Results string
}

func GetSearchResults(query string, page int) (searchResult SearchResult, err error) {
	db, err := openConnection()
	if err != nil {
		return SearchResult{}, err
	}
	defer db.Close()

	sql := `SELECT id, query, page, results FROM search WHERE query = ? AND page = ?`
	row := db.QueryRow(sql, query, page)

	err = row.Scan(&searchResult.Id, &searchResult.Query, &searchResult.Page, &searchResult.Results)
	if err != nil {
		return SearchResult{}, err
	}

	return searchResult, nil
}

func AddGame(guid string, data string) (id int, err error) {
	db, err := openConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sql := `INSERT INTO games (guid, results) VALUES (?, ?)`

	result, err := db.Exec(sql, guid, data)
	if err != nil {
		return 0, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id64), nil
}

type GameResult struct {
	Id      int
	Guid    string
	Results string
}

func GetGameResults(guid string) (gameResult GameResult, err error) {
	db, err := openConnection()
	if err != nil {
		return GameResult{}, err
	}
	defer db.Close()

	sql := `SELECT id, guid, results FROM games WHERE guid = ?`
	row := db.QueryRow(sql, guid)

	err = row.Scan(&gameResult.Id, &gameResult.Guid, &gameResult.Results)
	if err != nil {
		return GameResult{}, err
	}

	return gameResult, nil
}
