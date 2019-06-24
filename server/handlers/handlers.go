package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/domain"
	"main/infrastructure"
	"net/http"
	"strings"
	"time"
)

//Interface
type UsersHandler interface {
	MainPage(w http.ResponseWriter, r *http.Request)
	AddUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetUsersList(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	GetUserChild(w http.ResponseWriter, r *http.Request)

	Schools(w http.ResponseWriter, r *http.Request)
	Groups(w http.ResponseWriter, r *http.Request)
	Price(w http.ResponseWriter, r *http.Request)
}

//Struct
type usersHandler struct {
	usersRepository infrastructure.UsersRepository
	schoolsRepository infrastructure.SchoolRepository
	groupsRepository infrastructure.GroupRepository
}

//Constructor
func NewUsersHandler(usersRepository infrastructure.UsersRepository, schoolsRepository infrastructure.SchoolRepository,
	groupsRepository infrastructure.GroupRepository) UsersHandler {
	return &usersHandler{
		usersRepository: usersRepository,
		schoolsRepository: schoolsRepository,
		groupsRepository: groupsRepository,
	}
}

///////////////////////////////////////////////////////////////////////////////////////
//Handlers
///////////////////////////////////////////////////////////////////////////////////////

func (h *usersHandler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	users, err := h.usersRepository.GetUserList(context.Background())
	if err != nil {
		http.Redirect(w, r, "/", 301)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := json.Marshal(users)
	if err != nil{
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(res))
}

/////////////////////////////////
//MainPage
/////////////////////////////////
func (h *usersHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	// Fake auth check
	loggedIn := (err != http.ErrNoCookie)

	user, err := h.usersRepository.GetUserByID(context.Background(), session.Value)
	if err != nil{
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	if user != nil {
		if loggedIn {
			fmt.Fprintln(w, `<a href="/logout">logout</a>`)
			fmt.Fprintln(w, "Welcome, " + user.Name)
		} else {
			fmt.Fprintln(w, `<a href="/login">login</a>`)
			fmt.Fprintln(w, "You need to login")
		}
	}
}

/////////////////////////////////
//LogOut
/////////////////////////////////
func (h *usersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)}


////////////////
//Register
///////////////
func (h *usersHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reg domain.RegistrationData
	decoder.Decode(&reg)

	us, err := h.usersRepository.CreateUser(context.Background(), reg.Login, reg.Password)
	if err != nil {
		http.Redirect(w, r, "/", 301)
	}

	fmt.Printf("New User: \n Login - %s, Password - %s \n", reg.Login, reg.Password)
	fmt.Fprintf(w, "New User: \n Login - %s, Password - %s \n ID - %s \n", reg.Login, reg.Password, us.ID)
}

////////////////
//Login
///////////////
func (h *usersHandler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reg domain.RegistrationData
	err := decoder.Decode(&reg)
	if err != nil{
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	user, _ := h.usersRepository.UserLogin(context.Background(), reg.Login, reg.Password)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/register", http.StatusMovedPermanently)
	}
	w.Header().Set("Content-Type", "application/json")

	expiration := time.Now().Add(10 * time.Hour)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   user.ID.Hex(),
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)

	res, err := json.Marshal(user)
	if err != nil{
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(res))
}

func (h *usersHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(string(r.RequestURI), "/")

	if r.Method == http.MethodGet {
		if len(path) == 2 {                       //two params - address like /user/{name}



		} else if len(path) == 3 {                //three params - address like /user/{name}/{childname}



		}

	} else if r.Method == http.MethodPost {
		if len(path) == 2 {



		} else if len(path) == 3 {



		}

	} else {

	}
	//three params - address like /user/{name}/{childname}

	user, _ := h.usersRepository.GetUserByLogin(context.Background(), path[2])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, user)
}

func (h *usersHandler) GetUserChild(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GetUserChild worked")
}

func (h *usersHandler) Schools(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		schools, err := h.schoolsRepository.GetSchools(context.Background())
		if err != nil {
			http.Redirect(w, r, "/", 301)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res, err := json.Marshal(schools)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(res))
	} else if r.Method == http.MethodPost {

		decoder := json.NewDecoder(r.Body)
		var reg domain.School
		err := decoder.Decode(&reg)
		if err != nil{
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}

		newSchool, err := h.schoolsRepository.AddSchool(context.Background(), reg.Name)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/register", http.StatusMovedPermanently)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res, err := json.Marshal(newSchool)
		if err != nil{
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(res))

	} else{
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
	}

}

func (h *usersHandler) Groups(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		groups, err := h.groupsRepository.GetGroups(context.Background())
		if err != nil {
			http.Redirect(w, r, "/", 301)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res, err := json.Marshal(groups)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(res))

	} else if r.Method == http.MethodPost {

		decoder := json.NewDecoder(r.Body)
		var reg domain.CreateGroupData
		err := decoder.Decode(&reg)
		if err != nil{
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}

		gr, err := h.groupsRepository.AddGroup(context.Background(), reg.Name, reg.SchoolID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/register", http.StatusMovedPermanently)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res, err := json.Marshal(gr)
		if err != nil{
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(res))

	} else{
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
	}
}

func (h *usersHandler) Price(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GetUsersList worked")
}
