package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/domain"
	"main/infrastructure"
	"net/http"
)

type UsersHandler interface {
	AddUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetUsersList(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type usersHandler struct {
	usersRepository infrastructure.UsersRepository
}

func NewUsersHandler(usersRepository infrastructure.UsersRepository) UsersHandler {
	return &usersHandler{
		usersRepository: usersRepository,
	}
}

func (h *usersHandler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *usersHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reg domain.RegistrationData
	decoder.Decode(&reg)

	us, err := h.usersRepository.CreateUser(context.Background(), reg.Login, reg.Password)
	if err != nil{
		http.Redirect(w, r,"/", 301)
	}
	fmt.Printf("New User: \n Login - %s, Password - %s \n", reg.Login, reg.Password)
	fmt.Fprintf(w, "New User: \n Login - %s, Password - %s \n ID - %s \n", reg.Login, reg.Password, us.ID)
}

func (h *usersHandler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reg domain.RegistrationData
	decoder.Decode(&reg)

	ok := h.usersRepository.UserLogin(context.Background(), reg.Login, reg.Password)
	if !ok{
		http.Redirect(w, r,"/register", 301)
	}

	fmt.Printf("New User: \n Login - %s, Password - %s \n", reg.Login, reg.Password)
	fmt.Fprintf(w, "New User: \n Login - %s, Password - %s \n ID - %s \n", reg.Login, reg.Password, us.ID)
}

func (h *usersHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.RequestURI)

	user, _ := h.usersRepository.GetUserByLogin(context.Background(),"ldo")

	fmt.Fprintln(w, user)

}


