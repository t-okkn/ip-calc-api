package db

import (
	"errors"

	"github.com/go-gorp/gorp"
	"ip-calc-practice-api/models"
)

type IpRepository struct {
	*gorp.DbMap
}

func NewIpRepository(dm *gorp.DbMap) *IpRepository {
	return &IpRepository{dm}
}

func (r *IpRepository) GetID(id string) (models.TranID, error) {
	var result models.TranID
	query := GetSQL("get-id", "")
	val := map[string]interface{}{"id": id}

	if err := r.SelectOne(&result, query, val); err != nil {
		return models.TranID{}, err
	}

	return result, nil
}

func (r *IpRepository) GetExpired() ([]models.TranID, error) {
	var result []models.TranID
	query := GetSQL("get-expired", "")

	if _, err := r.Select(&result, query); err != nil {
		return []models.TranID{}, err
	}

	return result, nil
}

func (r *IpRepository) GetQuestion(id string, num int) (models.TranQuestion, error) {
	var result models.TranQuestion
	query := GetSQL("get-question", "")
	val := map[string]interface{}{"id": id, "number": num}

	if err := r.SelectOne(&result, query, val); err != nil {
		return models.TranQuestion{}, err
	}

	return result, nil
}

func (r *IpRepository) GetResults(id string) ([]models.TranQuestion, error) {
	var result []models.TranQuestion
	query := GetSQL("get-results", "")
	val := map[string]interface{}{"id": id}

	if _, err := r.Select(&result, query, val); err != nil {
		return []models.TranQuestion{}, err
	}

	return result, nil
}

func (r *IpRepository) CheckNow(id string) (models.NowNumber, error) {
	var result models.NowNumber
	query := GetSQL("check-now", "")
	val := map[string]interface{}{"id": id}

	if err := r.SelectOne(&result, query, val); err != nil {
		return models.NowNumber{}, err
	}

	return result, nil
}

func (r *IpRepository) InsertFirstData(mid models.TranID, tq models.TranQuestion) error {
	tx, err := r.Begin()

	if err != nil {
		return err
	}

	if err := tx.Insert(&mid); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Insert(&tq); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *IpRepository) InsertQuestion(tq models.TranQuestion) error {
	tx, err := r.Begin()

	if err != nil {
		return err
	}

	if err := tx.Insert(&tq); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *IpRepository) UpdateQuestion(tq models.TranQuestion) error {
	tx, err := r.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.Update(&tq); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *IpRepository) DeleteExpiredData(ids []models.TranID) error {
	str_ids := make([]string, len(ids), len(ids))
	for i, v := range ids {
		str_ids[i] = v.Id
	}

	var result []models.TranQuestion
	query := GetSQL("delete-expired", "")
	val := map[string]interface{}{"ids": str_ids}

	if _, err := r.Select(&result, query, val); err != nil {
		return err
	}

	if result == nil || len(result) == 0 {
		e := errors.New("No expired data in T_QUESTION")
		return e
	}

	tx, err := r.Begin()

	if err != nil {
		return err
	}

	errs := make([]error, 0, 2 * len(ids))

	for _, v := range ids {
		if _, err := tx.Delete(&v); err != nil {
			errs = append(errs, err)
		}
	}

	for _, v := range result {
		if _, err := tx.Delete(&v); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		tx.Rollback()

		return errs[0]
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

