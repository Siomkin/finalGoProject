package server

import (
	"fmt"
	"main/server/handlers"
	"net/http"
)

type Server struct {
	port        	string
	handlers		handlers.UsersHandler
	//adminsHandler	handlers.AdminsHandler
	//authMiddleware	middleware.AuthMiddleware
}

func NewServer(port string, usersHandler handlers.UsersHandler) *Server {
	return &Server{
		port:       port,
		handlers:   usersHandler,
		//adminsHandler:  adminsHandler,
		//authMiddleware: authMiddleware,
	}
}

func (server *Server) ConfigureAndRun() {

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/userslist", server.handlers.GetUsersList)

	//adminHandler := server.authMiddleware.CheckRole(adminMux)
	//adminHandler = server.authMiddleware.CheckAuth(adminHandler)

	userMux := http.NewServeMux()
	userMux.HandleFunc("/user/profile", server.handlers.GetUser)

	//userHandler := server.authMiddleware.CheckAuth(userMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminMux); //adminHandler)
	siteMux.Handle("/user/", userMux); //userHandler)
	siteMux.HandleFunc("/register", server.handlers.AddUser)

	fmt.Printf("listening at %s", server.port)
	http.ListenAndServe(server.port, siteMux)
}