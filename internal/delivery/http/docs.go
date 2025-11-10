//go:build !debug
// +build !debug

package http

import "github.com/gin-gonic/gin"

func AddDocsForDebugVersion(router *gin.Engine) {}
