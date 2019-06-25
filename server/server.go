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
	adminMux.HandleFunc("/admin/users", server.handlers.GetUsersList)
	adminMux.HandleFunc("/admin/schools/", server.handlers.Schools)
	adminMux.HandleFunc("/admin/groups/", server.handlers.Groups)
	//adminMux.HandleFunc("/admin/price/", server.handlers.Price)	 //prices will be in the next version

	//adminHandler := server.authMiddleware.CheckRole(adminMux)
	//adminHandler = server.authMiddleware.CheckAuth(adminHandler)

	userMux := http.NewServeMux()
	userMux.HandleFunc("/user/", server.handlers.MainUserHandler)

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

//
// child			get, post (getting list, posting new child)
// child/{name}		get

// tabel			get	 (getting list by dates), post ( post by date)
// tabel/date		get (get by specified date)


/////////////
//childMux := http.NewServeMux()
//childMux.HandleFunc("/child", server.handlers.MainChildHandler)		//req like /child
//childMux.HandleFunc("/child/", server.handlers.GetUserChild)		//req like /child/{name}
//
//tabelMux := http.NewServeMux()
//tabelMux.HandleFunc("/tabel", server.handlers.MainTabelHandler)		//req like /tabel
//tabelMux.HandleFunc("/tabel/", server.handlers.GetTabelOnDate) 		//req like /tabel/{date}


