package services

import (
	"net/http"
	"simple-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type noteService struct {
	db *gorm.DB
}

type NoteService interface {
	CreateNote(c *gin.Context)
	GetNotes(c *gin.Context)
	GetNote(c *gin.Context)
	UpdateNote(c *gin.Context)
	DeleteNote(c *gin.Context)
}

func (s *noteService) CreateNote(c *gin.Context) {
	// Parse JSON
	var json struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}

	if c.Bind(&json) == nil {
		// Create note
		note := models.Note{Title: json.Title, Body: json.Body}
		s.db.Create(&note)
		c.JSON(http.StatusCreated, note)
	}
}

func (s *noteService) GetNotes(c *gin.Context) {
	var notes []models.Note
	s.db.Find(&notes)
	c.JSON(http.StatusOK, notes)
}

func (s *noteService) GetNote(c *gin.Context) {
	var note models.Note
	if err := s.db.First(&note, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		}
	} else {
		c.JSON(http.StatusOK, note)
	}
}

func (s *noteService) UpdateNote(c *gin.Context) {
	var note models.Note
	if err := s.db.First(&note, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		}
	} else {
		var json struct {
			Title string `json:"title" binding:"required"`
			Body  string `json:"body"`
		}

		if c.Bind(&json) == nil {
			s.db.Model(&note).Updates(models.Note{Title: json.Title, Body: json.Body})
			c.JSON(http.StatusOK, note)
		}
	}
}

func (s *noteService) DeleteNote(c *gin.Context) {
	var note models.Note
	if err := s.db.First(&note, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		}
	} else {
		s.db.Delete(&note)
		c.JSON(http.StatusNoContent, nil)
	}
}

func NewNoteService(db *gorm.DB) *noteService {
	return &noteService{db}
}
