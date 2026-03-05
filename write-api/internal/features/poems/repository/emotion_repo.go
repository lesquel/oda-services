package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type emotionRepo struct{ db *gorm.DB }

func NewEmotionRepository(db *gorm.DB) domain.EmotionRepository { return &emotionRepo{db: db} }

func (r *emotionRepo) Tag(userID, poemID, emotionID string) error {
	tag := &domain.EmotionTag{
		ID:        uuid.NewString(),
		PoemID:    poemID,
		UserID:    userID,
		EmotionID: emotionID,
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(tag).Error
}

func (r *emotionRepo) Remove(userID, poemID, emotionTagID string) error {
	return r.db.Where("id = ? AND user_id = ? AND poem_id = ?", emotionTagID, userID, poemID).
		Delete(&domain.EmotionTag{}).Error
}

func (r *emotionRepo) GetByPoem(poemID string) ([]*domain.EmotionTag, error) {
	var tags []*domain.EmotionTag
	err := r.db.Where("poem_id = ?", poemID).Find(&tags).Error
	return tags, err
}

func (r *emotionRepo) GetDistribution(poemID string) (map[string]int, error) {
	type row struct {
		EmotionID string
		Count     int
	}
	var rows []row
	err := r.db.Model(&domain.EmotionTag{}).
		Select("emotion_id, COUNT(*) as count").
		Where("poem_id = ?", poemID).
		Group("emotion_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]int, len(rows))
	for _, r := range rows {
		result[r.EmotionID] = r.Count
	}
	return result, nil
}
