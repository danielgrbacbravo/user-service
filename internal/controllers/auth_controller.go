package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Temp Register handler
func Register(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Register handler not implemented yet")
}

// Temp Login handler
func Login(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Login handler not implemented yet")
}

// Temp RefreshToken handler
func RefreshToken(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "RefreshToken handler not implemented yet")
}

// Temp AppleLogin handler
func AppleLogin(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "AppleLogin handler not implemented yet")
}
