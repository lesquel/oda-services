package repository

import "github.com/lesquel/oda-shared/domain"

func (r *adminRepo) GetDashboardStats() (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}
	r.db.Model(&domain.User{}).Count(&stats.TotalUsers)
	r.db.Model(&domain.Poem{}).Count(&stats.TotalPoems)
	r.db.Model(&domain.Like{}).Count(&stats.TotalLikes)
	r.db.Model(&domain.Bookmark{}).Count(&stats.TotalBookmarks)
	r.db.Model(&domain.EmotionTag{}).Count(&stats.TotalEmotions)
	return stats, nil
}
