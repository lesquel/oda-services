package repository

import "github.com/lesquel/oda-shared/domain"

func (r *ReadRepository) GetUserBookmarks(userID string, limit, offset int) ([]*domain.Poem, int64, error) {
	var bookmarks []*domain.Bookmark
	var total int64
	r.db.Model(&domain.Bookmark{}).Where("user_id = ?", userID).Count(&total)
	if err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}
	poemIDs := make([]string, len(bookmarks))
	for i, b := range bookmarks {
		poemIDs[i] = b.PoemID
	}
	if len(poemIDs) == 0 {
		return []*domain.Poem{}, total, nil
	}
	var poems []*domain.Poem
	err := r.db.Preload("Author").Where("id IN ?", poemIDs).Find(&poems).Error
	return poems, total, err
}
