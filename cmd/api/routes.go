package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {
	// setup a new router
	router := httprouter.New()
	// handle 404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
   // handle 405
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	// setup routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	router.HandlerFunc(http.MethodPatch,"/v1/comments/:id", a.updateCommentHandler)
	router.HandlerFunc(http.MethodDelete,"/v1/comments/:id", a.deleteCommentHandler)
	router.HandlerFunc(http.MethodGet,"/v1/comments", a.listCommentsHandler)


	// Set the CORs middleware early in the middleware chain to enable our
	// response to be accepted by the client's browser
	return a.recoverPanic(a.enableCORS(a.rateLimit(a.authenticate(router))))

}
