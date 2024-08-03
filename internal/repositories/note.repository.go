package repositories

import (
	"simple-api/models"

	"gorm.io/gorm"
)

type NoteRepository interface {
	CreateNote(note *models.Note) error
	GetNotes() ([]models.Note, error)
	GetNoteByID(id uint) (*models.Note, error)
	UpdateNote(note *models.Note) error
	DeleteNote(id uint) error
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db}
}

func (r *noteRepository) CreateNote(note *models.Note) error {
	return r.db.Create(note).Error
}

func (r *noteRepository) GetNotes() ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Find(&notes).Error
	return notes, err
}

func (r *noteRepository) GetNoteByID(id uint) (*models.Note, error) {
	var note models.Note
	err := r.db.First(&note, id).Error
	return &note, err
}

func (r *noteRepository) UpdateNote(note *models.Note) error {
	return r.db.Save(note).Error
}

func (r *noteRepository) DeleteNote(id uint) error {
	return r.db.Delete(&models.Note{}, id).Error
}
