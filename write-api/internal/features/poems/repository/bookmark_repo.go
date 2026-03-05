package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

type bookmarkRepo struct{ db *gorm.DB }

func NewBookmarkRepository(db *gorm.DB) domain.BookmarkRepository { return &bookmarkRepo{db: db} }

func (r *bookmarkRepo) Toggle(userID, poemID string) (bool, error) {
	var existing domain.Bookmark
	err := r.db.Where("user_id = ? AND poem_id = ?", userID, poemID).First(&existing).Error
	if err == nil {
		if delErr := r.db.Delete(&existing).Error; delErr != nil {
			return false, delErr
		}
		return false, nil
	}
	bm := &domain.Bookmark{ID: uuid.NewString(), UserID: userID, PoemID: poemID}
	return true, r.db.Create(bm).Error
}

func (r *bookmarkRepo) IsBookmarked(userID, poemID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Bookmark{}).
		Where("user_id = ? AND poem_id = ?", userID, poemID).Count(&count).Error
	return count > 0, err
}

func (r *bookmarkRepo) GetUserBookmarks(userID string, limit, offset int) ([]*domain.Poem, error) {
	var bms []*domain.Bookmark
	if err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&bms).Error; err != nil {
		return nil, err
	}
	poemIDs := make([]string, len(bms))
	for i, b := range bms {
		poemIDs[i] = b.PoemID
	}
	if len(poemIDs) == 0 {
		return []*domain.Poem{}, nil
	}
	var poems []*domain.Poem
	err := r.db.Preload("Author").Where("id IN ?", poemIDs).Find(&poems).Error
	return poems, err
}
