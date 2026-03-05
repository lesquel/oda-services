package repository

import "github.com/lesquel/oda-shared/domain"

func toAdminUser(u domain.User) domain.AdminUser {
	au := domain.AdminUser{
		ID: u.ID, Username: u.Username, Email: u.Email,
		Role: u.Role, Bio: u.Bio, AvatarURL: u.AvatarURL,
		IsActive: u.IsActive, CreatedAt: u.CreatedAt,
	}
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		au.DeletedAt = &t
	}
	return au
}

func toAdminPoem(p domain.Poem) domain.AdminPoem {
	ap := domain.AdminPoem{
		ID: p.ID, AuthorID: p.AuthorID, Title: p.Title,
		Content: p.Content, Status: p.Status,
		LikesCount: p.LikesCount, ViewsCount: p.ViewsCount,
		CreatedAt: p.CreatedAt,
	}
	if p.DeletedAt.Valid {
		t := p.DeletedAt.Time
		ap.DeletedAt = &t
	}
	if p.Author != nil {
		ap.Author = &domain.AdminUser{
			ID:        p.Author.ID,
			Username:  p.Author.Username,
			AvatarURL: p.Author.AvatarURL,
		}
	}
	return ap
}
