package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func Login(ctx *gin.Context) {
	var d struct {
		User string `form:"username" json:"username" binding:"required"`
		Pass string `form:"password" json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&d); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if d.User == "admin" {
		claims := MyCustomClaims{
			"bar",
			jwt.RegisteredClaims{
				// A usual scenario is to set the expiration time relative to the current time
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "test",
				Subject:   "somebody",
				ID:        "1",
				Audience:  []string{"somebody_else"},
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString([]byte("garon"))
		fmt.Println(ss, err)

		ctx.JSON(http.StatusOK, gin.H{"data": ss})
		return
	}

	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found."})
}
