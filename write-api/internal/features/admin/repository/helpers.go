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
		ModerationStatus: p.ModerationStatus,
		ModerationScore:  p.ModerationScore,
		ModerationReason: p.ModerationReason,
		ModeratedAt:      p.ModeratedAt,
		ModeratedBy:      p.ModeratedBy,
		LikesCount:       p.LikesCount, ViewsCount: p.ViewsCount,
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

func toAdminLike(l domain.Like) domain.AdminLike {
	al := domain.AdminLike{
		ID: l.ID, UserID: l.UserID, PoemID: l.PoemID, CreatedAt: l.CreatedAt,
	}
	if l.DeletedAt.Valid {
		t := l.DeletedAt.Time
		al.DeletedAt = &t
	}
	return al
}

func toAdminBookmark(b domain.Bookmark) domain.AdminBookmark {
	ab := domain.AdminBookmark{
		ID: b.ID, UserID: b.UserID, PoemID: b.PoemID, CreatedAt: b.CreatedAt,
	}
	if b.DeletedAt.Valid {
		t := b.DeletedAt.Time
		ab.DeletedAt = &t
	}
	return ab
}

func toAdminEmotion(t domain.EmotionTag) domain.AdminEmotion {
	ae := domain.AdminEmotion{
		ID: t.ID, UserID: t.UserID, PoemID: t.PoemID,
		EmotionID: t.EmotionID, CreatedAt: t.CreatedAt,
	}
	if t.DeletedAt.Valid {
		tt := t.DeletedAt.Time
		ae.DeletedAt = &tt
	}
	return ae
}
