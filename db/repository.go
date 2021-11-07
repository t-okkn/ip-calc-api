package db

import (
	"context"

	"github.com/go-gorp/gorp"
	"github.com/t-okkn/ip-calc-practice-api/models"
)

type IpRepository struct {
	*gorp.DbMap
}

func NewIpRepository(dm *gorp.DbMap) *IpRepository {
	return &IpRepository{dm}
}

func (r *IpRepository) GetExpire(ctx context.Context, id string) (models.MstrID, error) {
	var result models.MstrID
	query := GetSQL("get-expire", "")
	val := map[string]interface{}{"id": id}

	if err := r.SelectOne(&result, query, val); err != nil {
		return models.MstrID{}, err
	}

	return result, nil
}

func (r *IpRepository) GetResults(ctx context.Context, id string) ([]models.TranQuestion, error) {
	var result []models.TranQuestion
	query := GetSQL("get-result", "")
	val := map[string]interface{}{"id": id}

	if _, err := r.Select(&result, query, val); err != nil {
		return []models.TranQuestion{}, err
	}

	return result, nil
}

