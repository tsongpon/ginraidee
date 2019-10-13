package adapter

import (
	"database/sql"
	"github.com/tsongpon/ginraidee/model"
	"log"
)

type SearchHistoryDBAdapter struct {
	db *sql.DB
}

func NewSearchHistoryDBAdapter(db *sql.DB) *SearchHistoryDBAdapter {
	repo := new(SearchHistoryDBAdapter)
	repo.db = db
	return repo
}

func (a *SearchHistoryDBAdapter) Save(history model.SearchHistory) (model.SearchHistory, error) {
	log.Printf("data %v", history)
	insert := "INSERT INTO searchhistory (id, userid, keyword, time) VALUES ($1, $2, $3, $4)"
	stmt, err := a.db.Prepare(insert)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.Exec(history.ID, history.UserID, history.Keyword, history.Time)

	return history, nil
}