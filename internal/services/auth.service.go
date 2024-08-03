package services

import (
	"fmt"
	"net/http"
	"simple-api/internal/repositories"
	"simple-api/models"
	"simple-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authService struct {
	db *gorm.DB
}

type AuthService interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

func (s *authService) Login(ctx *gin.Context) {
	hasher := utils.BcryptHasher{}
	clientIp := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	if userAgent == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User-Agent header is required"})
		return
	}

	var json struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if ctx.Bind(&json) == nil {
		newUser := &models.User{Username: json.Username}

		userRepo := repositories.NewUserRepository(s.db)

		var user *models.User
		var err error

		user, err = userRepo.GetUserByUserName(newUser.Username)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		if hasher.ComparePassword(user.Password, json.Password) != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		refreshParams := utils.RefreshParams{
			IpAddress: clientIp,
			UserAgent: userAgent,
			User:      user,
		}

		var refreshUtils = refreshParams
		refreshToken, err := refreshUtils.CreateJWT()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refersh token"})
			return
		}

		accessParams := utils.AccessParams{
			User: user,
		}

		var accessUtils = accessParams
		accessToken, err := accessUtils.CreateJWT()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		}

		if accessToken != "" && refreshToken != "" {
			ectryptedRefreshToken, err := utils.Encrypt(refreshToken)
			fmt.Println("before crypt", refreshToken)
			fmt.Println("after crypt", ectryptedRefreshToken)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error", "details": err})
				return
			}

			ctx.SetCookie(
				"refresh_token",       // cookie name
				ectryptedRefreshToken, // cookie value
				3600*24*7,             // max age in seconds (1 week)
				"/",                   // path
				"",                    // domain (leave empty for default)
				true,                  // secure flag
				true,                  // HttpOnly flag
			)

			ctx.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "user": user})

			s.db.Model(user).Update("ip_address", clientIp)
			s.db.Model(user).Update("user_agent", userAgent)
		}
	}
}

func (s *authService) RefreshToken(c *gin.Context) {
	clientIp := c.ClientIP()
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access error. Go fuck yourself"})
	}

	decryptedRefreshToken, err := utils.Decrypt(refreshToken)

	fmt.Println("refresh", decryptedRefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access error"})
		return
	}

	refreshParams := utils.RefreshParams{
		IpAddress: clientIp,
		UserAgent: c.GetHeader("User-Agent"),
	}

	var refreshUtils = refreshParams
	userData, err := refreshUtils.ExtractUserFromToken(decryptedRefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	refreshClaims := userData.(*utils.RefreshClaims)

	var user models.User

	if err := s.db.Where("id = ? AND ip_address = ? AND user_agent = ?", refreshClaims.ID, refreshClaims.IpAddress, refreshClaims.UserAgent).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if refreshClaims.IpAddress != clientIp || refreshClaims.UserAgent != c.GetHeader("User-Agent") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Why are you stealing tokens?"})
		return
	}

	if _, err := refreshUtils.VerifyToken(decryptedRefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access error"})
	}

	jwtToken, err := utils.AccessParams{
		User: &user,
	}.CreateJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db}
}
