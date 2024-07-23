package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckAdmin(c *gin.Context) (*AccessClaims, error) {
	claims, exists := c.Get("user")

	if !exists {
		return nil, errors.New("user claim does not exist")
	}

	myClaims, ok := claims.(*AccessClaims)

	if !ok {
		return nil, errors.New("invalid user claim")
	}

	if !myClaims.IsAdmin {
		return nil, errors.New("user is not an admin")
	}

	return myClaims, nil
}
