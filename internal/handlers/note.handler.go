package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"simple-api/internal/services"
)

type NoteHandler struct {
	service services.NoteService
}

func NewNoteHandler(db *gorm.DB) *NoteHandler {
	return &NoteHandler{
		service: services.NewNoteService(db),
	}
}

func (h *NoteHandler) CreateNote(ctx *gin.Context) {
	var noteRequest struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&noteRequest); err == nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request"})
		return
	}

	note, err := h.service.CreateNote(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, note)
	return
}

func (h *NoteHandler) GetNote(ctx *gin.Context) {
	note, err := h.service.GetNote(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, note)
	return
}

func (h *NoteHandler) GetNotes(ctx *gin.Context) {
	notes := h.service.GetNotes(ctx)

	ctx.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) UpdateNote(ctx *gin.Context) {
	var noteRequest struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&noteRequest); err == nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request"})
	}
	note, err := h.service.UpdateNote(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(ctx *gin.Context) {
	err := h.service.DeleteNote(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, nil)
	return
}
