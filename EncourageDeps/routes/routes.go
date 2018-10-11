package routes

import (
	"fmt"
	"log"
	"net/http"
	"trh/EncourageDeps/middleware"
	"trh/EncourageDeps/models"
	"trh/EncourageDeps/sessions"
	"trh/EncourageDeps/utils"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/logout", logoutGetHandler).Methods("GET")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	r.HandleFunc("/getuser", GetCurrentUserHandler).Methods("GET")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/{username}",
		middleware.AuthRequired(userGetHandler)).Methods("GET")
	return r
}

// My first closure to return a username

func usrFunc() func(*http.Request) string {

	return func(r *http.Request) string {
		session, _ := sessions.Store.Get(r, "session")
		untypedUserId := session.Values["user_id"]
		CurrentUserId, _ := untypedUserId.(int64)
		//	log.Printf(" Closure User id %v ", CurrentUserId)
		ThisUser := fmt.Sprintf("user:%d", CurrentUserId)
		//fmt.Println("thins clousrre user : ", ThisUser)
		Username, _ := models.GetUsername(ThisUser)
		//log.Println(" Closure User name %v ", Username) //v
		return Username
	}
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	updates, err := models.GetAllUpdates()
	u := usrFunc()
	fmt.Println("indexGetHandler run ")
	log.Println("THIS u as closure USERNAMe  ", u(r))

	type info struct {
		Title       string
		Updates     []*models.Update
		DisplayForm bool
		Username    string
	}
	usr := info{
		Title:       "All updates",
		Updates:     updates,
		DisplayForm: true,
		Username:    "ThisUsername",
	}

	usr.Username = u(r)

	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.gohtml", usr)
} /*  Saved Origional
utils.ExecuteTemplate(w, "index.gohtml", struct {
	Title       string
	Updates     []*models.Update
	DisplayForm bool
	Username    string
}{
	Title:       "All updates",
	Updates:     updates,
	DisplayForm: true,
	Username:    ThisUsername,
})
} */

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	currentUserId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUserByUsername(username)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	updates, err := models.GetUpdates(userId)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.gohtml", struct {
		Title       string
		Updates     []*models.Update
		DisplayForm bool
	}{
		Title:       username,
		Updates:     updates,
		DisplayForm: currentUserId == userId,
	})
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.gohtml", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.gohtml", "unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.gohtml", "invalid login")
		default:
			utils.InternalServerError(w)
		}
		return
	}
	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["user_id"] = userId
	session.Save(r, w)
	// Save Session State

	// for usernames
	http.Redirect(w, r, "/", 302)
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.gohtml", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err == models.ErrUsernameTaken {
		utils.ExecuteTemplate(w, "register.gohtml", "username taken")
		return
	} else if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	currentUserId, ok := untypedUserId.(int64)
	log.Printf(" User id %v ", currentUserId)
	ThisUser := fmt.Sprintf("user:%d", currentUserId)
	Username, _ := models.GetUsername(ThisUser)
	log.Printf(" User name %v ", Username) //var user1 models.User
	//user1 = models.GetUserById(2)
	//log.Printf((models.GetUserById(currentUserId).GetUsername()))

	if !ok {
		utils.InternalServerError(w)
		return
	}

}
