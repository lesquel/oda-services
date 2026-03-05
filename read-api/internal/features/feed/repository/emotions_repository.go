package repository

import "github.com/lesquel/oda-shared/domain"

func (r *ReadRepository) GetEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	var items []*domain.EmotionCatalog
	return items, r.db.Find(&items).Error
}

func (r *ReadRepository) GetEmotionDistribution(poemID string) (map[string]int, error) {
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
	for _, row := range rows {
		result[row.EmotionID] = row.Count
	}
	return result, nil
}
