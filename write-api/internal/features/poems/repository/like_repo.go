package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

type likeRepo struct{ db *gorm.DB }

func NewLikeRepository(db *gorm.DB) domain.LikeRepository { return &likeRepo{db: db} }

func (r *likeRepo) Toggle(userID, poemID string) (bool, error) {
	var existing domain.Like
	err := r.db.Where("user_id = ? AND poem_id = ?", userID, poemID).First(&existing).Error
	if err == nil {
		// Unlike
		if delErr := r.db.Delete(&existing).Error; delErr != nil {
			return false, delErr
		}
		r.db.Model(&domain.Poem{}).Where("id = ?", poemID).
			UpdateColumn("likes_count", gorm.Expr("GREATEST(likes_count - 1, 0)"))
		return false, nil
	}
	// Like
	like := &domain.Like{ID: uuid.NewString(), UserID: userID, PoemID: poemID}
	if createErr := r.db.Create(like).Error; createErr != nil {
		return false, createErr
	}
	r.db.Model(&domain.Poem{}).Where("id = ?", poemID).
		UpdateColumn("likes_count", gorm.Expr("likes_count + 1"))
	return true, nil
}

func (r *likeRepo) IsLiked(userID, poemID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Like{}).Where("user_id = ? AND poem_id = ?", userID, poemID).Count(&count).Error
	return count > 0, err
}

func (r *likeRepo) GetUserLikes(userID string, limit, offset int) ([]*domain.Poem, error) {
	var likes []*domain.Like
	if err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&likes).Error; err != nil {
		return nil, err
	}
	poemIDs := make([]string, len(likes))
	for i, l := range likes {
		poemIDs[i] = l.PoemID
	}
	if len(poemIDs) == 0 {
		return []*domain.Poem{}, nil
	}
	var poems []*domain.Poem
	err := r.db.Preload("Author").Where("id IN ?", poemIDs).Find(&poems).Error
	return poems, err
}
