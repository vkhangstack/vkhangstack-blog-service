package services

import (
	"encoding/json"

	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/meili"
)

type SearchEngineService struct {
	engine *meili.MeilisearchAdapter
}

func NewSearchEngineService(engine *meili.MeilisearchAdapter) *SearchEngineService {
	return &SearchEngineService{
		engine: engine,
	}
}

func (s *SearchEngineService) IndexDocument(indexName string, document interface{}) error {
	_, err := s.engine.AddDocument(indexName, document)
	if err != nil {
		return err
	}
	return nil
}

func (s *SearchEngineService) Search(indexName string, query string, limit int) ([]map[string]interface{}, error) {
	task, err := s.engine.Search(indexName, query, limit)
	if err != nil {
		return nil, err
	}
	res := make([]map[string]interface{}, 0, len(task.Hits))
	for _, hit := range task.Hits {
		m := make(map[string]interface{})
		for key, raw := range hit {
			var v interface{}
			if err := json.Unmarshal(raw, &v); err != nil {
				m[key] = string(raw)
				continue
			}
			m[key] = v
		}
		res = append(res, m)
	}

	return res, nil
}

func (s *SearchEngineService) DeleteDocument(indexName string, documentID string) error {
	_, err := s.engine.DeleteDocument(indexName, documentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SearchEngineService) UpdateDocument(indexName string, document interface{}) error {
	_, err := s.engine.UpdateDocument(indexName, document)
	if err != nil {
		return err
	}
	return nil
}
