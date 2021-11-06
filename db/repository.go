package db

import (
	"context"

	"github.com/go-gorp/gorp"
	"github.com/t-okkn/ip-calc-practice-api/models"
)

type IpRepository gorp.SqlExecutor

func NewIpRepository(exec gorp.SqlExecutor) *IpRepository {
	return &IpRepository(exec)
}

func (r *IpRepository) GetExpire(ctx context.Context, id string) (models.MstrID, error) {
	var result models.MstrID
	query := GetSQL("get-expire", "")
	val := map[string]interface{}{"id": id}

	if err := r.exec.SelectOne(&result, query, val); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *IpRepository) GetResults(ctx context.Context, id string) ([]models.TranQuestion, error) {
	var result []models.TranQuestion
	query := GetSQL("get-result", "")
	val := map[string]interface{}{"id": id}

	if _, err := r.exec.Select(&result, query, val); err != nil {
		return nil, err
	}

	return result, nil
}

