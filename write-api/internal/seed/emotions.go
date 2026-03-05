package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedEmotions upserts the default emotion catalog. Safe to call on every startup.
func SeedEmotions(db *gorm.DB) {
	emotions := []*domain.EmotionCatalog{
		{Name: "Melancólico", Emoji: "😔", Description: "Tristeza reflexiva, contemplación profunda"},
		{Name: "Esperanzador", Emoji: "🌟", Description: "Optimista, mirando hacia adelante"},
		{Name: "Sereno", Emoji: "☮️", Description: "Paz, calma, tranquilidad"},
		{Name: "Apasionado", Emoji: "🔥", Description: "Emoción intensa, fervoroso"},
		{Name: "Nostálgico", Emoji: "🍂", Description: "Anhelo del pasado, melancolía dulce"},
		{Name: "Inspirador", Emoji: "✨", Description: "Edificante, motivador, llena de energía"},
		{Name: "Romántico", Emoji: "🌹", Description: "Amor, ternura, conexión emocional profunda"},
		{Name: "Misterioso", Emoji: "🌙", Description: "Enigmático, oscuro, lleno de profundidad"},
		{Name: "Alegre", Emoji: "☀️", Description: "Luminoso, festivo, desbordante de vida"},
		{Name: "Doliente", Emoji: "💧", Description: "Dolor profundo, pérdida, luto silencioso"},
		{Name: "Rebelde", Emoji: "⚡", Description: "Desafiante, inconformista, voz propia"},
		{Name: "Tierno", Emoji: "🌺", Description: "Delicado, suave, llena de cariño"},
	}
	for _, e := range emotions {
		if e.ID == "" {
			e.ID = uuid.NewString()
		}
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"emoji", "description"}),
		}).Create(e).Error; err != nil {
			log.Printf("⚠️  emotion seed failed for %s: %v", e.Name, err)
		}
	}
	log.Println("✅ Emotion catalog seeded (12 emotions)")
}
