package redis

import (
	"catalog_service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// ПОка в разработке!!!!

const categoriesKey = "categories:all"

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (s *Cache) hasCategories(ctx context.Context) (bool, error) {
	exists, err := s.client.Exists(ctx, categoriesKey).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists failed: %w", err)
	}
	return exists > 0, nil
}

func (s *Cache) getCategoriesFromCache(ctx context.Context) ([]models.Categories, error) {
	data, err := s.client.Get(ctx, categoriesKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis get failed: %w", err)
	}
	var categories []models.Categories
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	return categories, nil
}

func (s *Cache) fetchCategoriesFromDBAndCache(ctx context.Context) ([]models.Categories, error) {
	categories, err := s.db.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories from db: %w", err)
	}
	data, err := json.Marshal(categories)
	if err != nil {
		// логируем, но не возвращаем ошибку
		fmt.Printf("warning: marshal categories failed: %v\n", err)
	} else {
		if err := s.client.Set(ctx, categoriesKey, data, 0).Err(); err != nil {
			fmt.Printf("warning: failed to set categories cache: %v\n", err)
		}
	}
	return categories, nil
}

func (s *Cache) GetCategories(ctx context.Context) ([]models.Categories, error) {
	has, err := s.hasCategories(ctx)
	if err != nil {
		return nil, err
	}
	if has {
		return s.getCategoriesFromCache(ctx)
	}
	return s.fetchCategoriesFromDBAndCache(ctx)
}
