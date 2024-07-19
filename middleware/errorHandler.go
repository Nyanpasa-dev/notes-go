package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// Перебор всех ошибок
			for _, e := range c.Errors {
				// Определение типа ошибки и установка соответствующего HTTP-статуса
				switch e.Type {
				case gin.ErrorTypePublic:
					// Если ошибка должна быть показана пользователю
					c.JSON(c.Writer.Status(), gin.H{"error": e.Error()})
				case gin.ErrorTypeBind:
					// Ошибки связанные с привязкой данных
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				default:
					// Внутренние ошибки сервера
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				}
			}
			c.Abort() // Прекращение обработки запроса
		}
	}
}
