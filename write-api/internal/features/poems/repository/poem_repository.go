package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ── Poem repository ───────────────────────────────────────────────────────────

type poemRepo struct{ db *gorm.DB }

func NewPoemRepository(db *gorm.DB) domain.PoemRepository { return &poemRepo{db: db} }

func (r *poemRepo) Create(poem *domain.Poem) error {
	if poem.ID == "" {
		poem.ID = uuid.NewString()
	}
	return r.db.Create(poem).Error
}

func (r *poemRepo) FindByID(id string) (*domain.Poem, error) {
	var poem domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").First(&poem, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &poem, nil
}

func (r *poemRepo) Update(poem *domain.Poem) error {
	return r.db.Omit("Author", "EmotionTags").Save(poem).Error
}

func (r *poemRepo) Delete(id string) error {
	return r.db.Delete(&domain.Poem{}, "id = ?", id).Error
}

func (r *poemRepo) GetFeed(limit, offset int) ([]*domain.Poem, error) {
	var poems []*domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, err
}

func (r *poemRepo) GetUserPoems(userID string, limit, offset int) ([]*domain.Poem, error) {
	var poems []*domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, err
}

func (r *poemRepo) Search(query string, limit, offset int) ([]*domain.Poem, error) {
	var poems []*domain.Poem
	like := "%" + query + "%"
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("status = ? AND (title ILIKE ? OR content ILIKE ?)", "published", like, like).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, err
}

func (r *poemRepo) IncrementViews(id string) error {
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
}

func (r *poemRepo) GetStats(poemID string) (map[string]interface{}, error) {
	var poem domain.Poem
	if err := r.db.Select("id, likes_count, views_count").First(&poem, "id = ?", poemID).Error; err != nil {
		return nil, err
	}
	var emotionCount int64
	r.db.Model(&domain.EmotionTag{}).Where("poem_id = ?", poemID).Count(&emotionCount)
	return map[string]interface{}{
		"likes_count":   poem.LikesCount,
		"views_count":   poem.ViewsCount,
		"emotion_count": emotionCount,
	}, nil
}

// ── Like repository ───────────────────────────────────────────────────────────

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

// ── Emotion repository ────────────────────────────────────────────────────────

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

// ── Bookmark repository ───────────────────────────────────────────────────────

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

// ── EmotionCatalog repository ─────────────────────────────────────────────────

type emotionCatalogRepo struct{ db *gorm.DB }

func NewEmotionCatalogRepository(db *gorm.DB) domain.EmotionCatalogRepository {
	return &emotionCatalogRepo{db: db}
}

func (r *emotionCatalogRepo) FindAll() ([]*domain.EmotionCatalog, error) {
	var items []*domain.EmotionCatalog
	return items, r.db.Find(&items).Error
}

func (r *emotionCatalogRepo) FindByID(id string) (*domain.EmotionCatalog, error) {
	var item domain.EmotionCatalog
	return &item, r.db.First(&item, "id = ?", id).Error
}

func (r *emotionCatalogRepo) FindByName(name string) (*domain.EmotionCatalog, error) {
	var item domain.EmotionCatalog
	return &item, r.db.First(&item, "name = ?", name).Error
}

func (r *emotionCatalogRepo) Create(e *domain.EmotionCatalog) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	return r.db.Create(e).Error
}

func (r *emotionCatalogRepo) Update(e *domain.EmotionCatalog) error {
	return r.db.Save(e).Error
}

func (r *emotionCatalogRepo) Delete(id string) error {
	return r.db.Delete(&domain.EmotionCatalog{}, "id = ?", id).Error
}
