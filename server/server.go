package server

import (
	"fmt"
	"log"
	"main/server/handlers"
	"net/http"
)

type Server struct {
	port     string
	handlers handlers.UsersHandler
	//adminsHandler	handlers.AdminsHandler
	//authMiddleware	middleware.AuthMiddleware
}

func NewServer(port string, usersHandler handlers.UsersHandler) *Server {
	return &Server{
		port:     port,
		handlers: usersHandler,
		//adminsHandler:  adminsHandler,
		//authMiddleware: authMiddleware,
	}
}

func (server *Server) ConfigureAndRun() {

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/users/", server.handlers.GetUsersList)
	adminMux.HandleFunc("/admin/schools/", server.handlers.Schools)
	adminMux.HandleFunc("/admin/groups/", server.handlers.Groups)
	adminMux.HandleFunc("/admin/price/", server.handlers.Price)

	//adminHandler := server.authMiddleware.CheckRole(adminMux)
	//adminHandler = server.authMiddleware.CheckAuth(adminHandler)

	userMux := http.NewServeMux()
	//userMux.HandleFunc("/", server.handlers.GetUser)
	userMux.HandleFunc("/user/", server.handlers.GetUser)
	//userMux.HandleFunc("/user/{name}/{child}/", server.handlers.GetUserChild)
	userMux.HandleFunc("/user/children/", server.handlers.GetUserChild)


	//userHandler := server.authMiddleware.CheckAuth(userMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminMux) //adminHandler)
	siteMux.Handle("/user/", userMux)   //userHandler)

	siteMux.HandleFunc("/", server.handlers.MainPage)
	siteMux.HandleFunc("/register", server.handlers.AddUser)
	siteMux.HandleFunc("/login", server.handlers.Login)
	siteMux.HandleFunc("/logout", server.handlers.Logout)

	fmt.Printf("listening at %s", server.port)
	log.Fatal(http.ListenAndServe(server.port, siteMux))
}

// login
// register
// admin/
//		 users		  get
//		 schools	  get, post
//		 groups		  get, post
//		 groups/{id}  get
//		 price		  get, post
//
// user/{name}			get
//		 users			get
//		 tabel			get, post
//		 children		get, post
//		 children/name	get
