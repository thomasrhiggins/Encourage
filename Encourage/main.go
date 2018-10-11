package main

import (
	"net/http"
	"trh/EncourageDeps/models"
	"trh/EncourageDeps/routes"
	"trh/EncourageDeps/utils"

	"google.golang.org/appengine" // Required external App Engine library
)

/*
type Info struct {
	Title       string
	Updates     []*models.Update
	DisplayForm bool
	Username    string
}
UserData := Info{
	Title:       username,
	Updates:     updates,
	DisplayForm: currentUserId == userId,
}
*/

func main() {

	models.Init()
	utils.LoadTemplates("templates/*.*")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

	appengine.Main()
}
