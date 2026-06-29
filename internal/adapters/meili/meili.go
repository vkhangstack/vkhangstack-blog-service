package meili

import (
	"context"

	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type MeilisearchConfig struct {
	Host   string
	APIKey string
}

type MeilisearchAdapter struct {
	client meilisearch.ServiceManager
}

func NewMeilisearchAdapter(ctx context.Context, cfg MeilisearchConfig) (*MeilisearchAdapter, error) {
	client := meilisearch.New(cfg.Host, meilisearch.WithAPIKey(cfg.APIKey))
	return &MeilisearchAdapter{
		client: client,
	}, nil
}

func (m *MeilisearchAdapter) GetClient() meilisearch.ServiceManager {
	return m.client
}

func (m *MeilisearchAdapter) AddDocument(index string, document interface{}) (*meilisearch.TaskInfo, error) {
	task, err := m.client.Index(index).AddDocuments(document, &meilisearch.DocumentOptions{
		PrimaryKey: utils.StringPtr("id"),
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *MeilisearchAdapter) Search(index string, query string, limit int) (*meilisearch.SearchResponse, error) {
	searchResponse, err := m.client.Index(index).Search(query, &meilisearch.SearchRequest{
		Limit: int64(limit),
	})
	if err != nil {
		return nil, err
	}
	return searchResponse, nil
}

func (m *MeilisearchAdapter) DeleteDocument(index string, documentID string) (*meilisearch.TaskInfo, error) {
	task, err := m.client.Index(index).DeleteDocument(documentID, &meilisearch.DocumentOptions{})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *MeilisearchAdapter) UpdateDocument(index string, document interface{}) (*meilisearch.TaskInfo, error) {
	task, err := m.client.Index(index).UpdateDocuments(document, &meilisearch.DocumentOptions{
		PrimaryKey: utils.StringPtr("id"),
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}
