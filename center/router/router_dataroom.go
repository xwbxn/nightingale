package router

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"
)

// api函数
func (rt *Router) drProxy(c *gin.Context) {
	var target = "127.0.0.1:8081"
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.Host = target

		// fe request e.g. /api/n9e/proxy/:id/*
		arr := strings.Split(c.Request.URL.Path, "/")
		if len(arr) < 3 {
			c.String(http.StatusBadRequest, "invalid url path")
			return
		}

		req.URL.Path = "/" + strings.Join(arr[2:], "/")
		req.URL.RawQuery = c.Request.URL.RawQuery
		logger.Debug("-------------------------------------------------------------")
		logger.Debug(req.URL.Path)
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}
