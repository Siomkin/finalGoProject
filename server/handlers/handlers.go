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
	MainUserHandler(w http.ResponseWriter, r *http.Request)
	GetUsersList(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	//GetUserChild(w http.ResponseWriter, r *http.Request)

	Schools(w http.ResponseWriter, r *http.Request)
	Groups(w http.ResponseWriter, r *http.Request)
	Price(w http.ResponseWriter, r *http.Request)
}

//Struct
type usersHandler struct {
	usersRepository infrastructure.UsersRepository
	schoolsRepository infrastructure.SchoolRepository
	groupsRepository infrastructure.GroupRepository
	childrepository infrastructure.ChildrenRepository
	tabelRepository infrastructure.TabelRepository
}

//Constructor
func NewUsersHandler(usersRepository infrastructure.UsersRepository, schoolsRepository infrastructure.SchoolRepository,
	groupsRepository infrastructure.GroupRepository, childrepository infrastructure.ChildrenRepository,
	tabelRepository infrastructure.TabelRepository) UsersHandler {
	return &usersHandler{
		usersRepository: usersRepository,
		schoolsRepository: schoolsRepository,
		groupsRepository: groupsRepository,
		childrepository: childrepository,
		tabelRepository: tabelRepository,
	}
}

///////////////////////////////////////////////////////////////////////////////////////
//Handlers
///////////////////////////////////////////////////////////////////////////////////////

func (h *usersHandler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	users, err := h.usersRepository.GetUserList(context.Background())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(users)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	http.Redirect(w, r, "/", http.StatusFound)
}

////////////////
//Register
///////////////
func (h *usersHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reg domain.RegistrationData
	decoder.Decode(&reg)

	us, err := h.usersRepository.CreateUser(context.Background(), reg.Login, reg.Password)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
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
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	res, err := json.Marshal(user)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(res))
}

func (h *usersHandler) MainUserHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(string(r.RequestURI), "/")
	userID, err := r.Cookie("session_id")
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusProxyAuthRequired)
		return
	}

	if r.Method == http.MethodGet {
		if len(path) == 3 {                       //three params - address like /user/{name}
			//user info
			user, err := h.usersRepository.GetUserByID(context.Background(), userID.Value)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			mUser, err := json.Marshal(user)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mUser)
			return

		} else if len(path) == 4 {                //four params - address like /user/{name}/child
			//child list
			childList, err := h.childrepository.GetChildrenByUserID(context.Background(), userID.Value)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			mChildList, err := json.Marshal(childList)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mChildList)
			return

		} else if len(path) == 5 {				////five params - address like /user/{name}/child/{childname}
			//child info
			var child *domain.Children
 			childName := path[4]
			childList, err := h.childrepository.GetChildrenByUserID(context.Background(), userID.Value)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for ind, el := range childList {
				if el.Name == childName{
					child = childList[ind]
				}
			}

			if child != nil{				// if we found a child

				mChild, err := json.Marshal(childList)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, mChild)
				return

			}else {

				w.WriteHeader(http.StatusBadRequest)
				return

			}

		}else if len(path) == 6 { 				////six params - address like /user/{name}/child/{childname}/tabel
			//tabel on the default period (current month) or with params as body

			childName := path[4]
			child, err := h.childrepository.GetChildrenByNameAndUserID(context.Background(), childName, userID.Value)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//default dates:
			t := time.Now()
			currentYear, currentMonth, _ := t.Date()
			currentLocation := t.Location()
			datefrom := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
			dateto := datefrom.AddDate(0, 1, -1)

			tabel, err := h.tabelRepository.GetChildTabel(context.Background(), child.ID, datefrom.Unix(), dateto.Unix())
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			mTabel, err :=json.Marshal(tabel)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mTabel)
			return
		}
	} else if r.Method == http.MethodPost {
		if len(path) == 4 {					//four params - address like /user/{name}/child
			//add child
			decoder := json.NewDecoder(r.Body)
			var reg domain.CreateChildData
			err := decoder.Decode(&reg)
			if err != nil{
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ch, err := h.childrepository.AddChild(context.Background(), userID.Value, reg.GroupID, reg.Name)
			if err != nil{
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			mch, err := json.Marshal(ch)
			if err != nil{
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mch)
			return

		} else if len(path) == 7 { ////six params - address like /user/{name}/child/{childname}/tabel/{date}
			//add day state (was child at school or not)

			//1. parse date
			layout := "2006-01-02"
			str := path[6] ///"2014-11-12"
			t, err := time.Parse(layout, str)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//2. find child
			childName := path[4]
			child, err := h.childrepository.GetChildrenByNameAndUserID(context.Background(), childName, userID.Value)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//3. parse body for status
			decoder := json.NewDecoder(r.Body)
			var reg domain.TabelRecord
			err = decoder.Decode(&reg)
			if err != nil{
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			//setting up the value
			tr, err := h.tabelRepository.SetDayValue(context.Background(), t.Unix(), child.ID.Hex(), reg.Value)

			mtr, err := json.Marshal(tr)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, mtr)
			return


		} else{
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//three params - address like /user/{name}/{childname}
	//user, _ := h.usersRepository.GetUserByLogin(context.Background(), path[2])
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintln(w, user)
}

func (h *usersHandler) Schools(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		schools, err := h.schoolsRepository.GetSchools(context.Background())
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(schools)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(newSchool)
		if err != nil{
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(res))

	} else{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func (h *usersHandler) Groups(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		groups, err := h.groupsRepository.GetGroups(context.Background())
		if err != nil {
			http.Redirect(w, r, "/", 301)
		}

		res, err := json.Marshal(groups)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(res))

	} else if r.Method == http.MethodPost {

		decoder := json.NewDecoder(r.Body)
		var reg domain.CreateGroupData
		err := decoder.Decode(&reg)
		if err != nil{
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		gr, err := h.groupsRepository.AddGroup(context.Background(), reg.Name, reg.SchoolID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/register", http.StatusMovedPermanently)
		}

		res, err := json.Marshal(gr)
		if err != nil{
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(res))

	} else{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *usersHandler) Price(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GetUsersList worked")
}
