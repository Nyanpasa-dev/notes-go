package services

import (
	"errors"
	"simple-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type noteService struct {
	db *gorm.DB
}

type NoteService interface {
	CreateNote(c *gin.Context) (*models.Note, error)
	GetNotes(c *gin.Context) *[]models.Note
	GetNote(c *gin.Context) (*models.Note, error)
	UpdateNote(c *gin.Context) (*models.Note, error)
	DeleteNote(c *gin.Context) error
}

func NewNoteService(db *gorm.DB) NoteService {
	return &noteService{db: db}
}

func (s *noteService) CreateNote(c *gin.Context) (*models.Note, error) {
	// Parse JSON
	var json struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}

	if c.Bind(&json) == nil {
		// Create note
		note := &models.Note{Title: json.Title, Body: json.Body}
		s.db.Create(&note)

		return note, nil
	}

	return nil, errors.New("invalid data")
}

func (s *noteService) GetNotes(c *gin.Context) *[]models.Note {
	var notes []models.Note
	s.db.Find(&notes)
	return &notes
}

func (s *noteService) GetNote(c *gin.Context) (*models.Note, error) {
	var note models.Note
	if err := s.db.First(&note, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note does not exist")
		}
	}

	return &note, nil
}

func (s *noteService) UpdateNote(c *gin.Context) (*models.Note, error) {
	var note models.Note
	if err := s.db.First(&note, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note does not exist")
		}
	} else {
		var json struct {
			Title string `json:"title" binding:"required"`
			Body  string `json:"body"`
		}

		if c.Bind(&json) == nil {
			s.db.Model(&note).Updates(models.Note{Title: json.Title, Body: json.Body})
			return &note, nil
		}
	}
	return nil, errors.New("invalid data")

}

func (s *noteService) DeleteNote(c *gin.Context) error {
	var note models.Note
	if err := s.db.First(&note, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("note does not exist")
		}
	} else {
		s.db.Delete(&note)
	}

	return nil
}
