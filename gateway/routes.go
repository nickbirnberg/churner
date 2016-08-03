package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getRouter() http.Handler {
	router := httprouter.New()

	router.GET("/api/v1/actions/:action", GetAction)
	router.POST("/api/v1/actions", PostAction)

	return router
}
