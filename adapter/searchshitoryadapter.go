package adapter

import "github.com/tsongpon/ginraidee/model"

type SearchHistoryAdapter interface {
	Save(history model.SearchHistory) (model.SearchHistory, error)
}
