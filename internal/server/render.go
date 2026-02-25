package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// Renders templ components or error.
func render(ctx *gin.Context, status int, template templ.Component) {
	ctx.Status(status)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err := template.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
