package controllers

import "github.com/kryptomind/bidboxapi/bitgetms/middleware"

func (r *Server) initializeRoutes() {
	s := r.Router.PathPrefix("/trades").Subrouter()
	s.HandleFunc("/", middleware.MiddlewareJSON(r.Home)).Methods("GET")
}
