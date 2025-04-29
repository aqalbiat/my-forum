package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// merged from 2024-07-13

var sessionList UID
var categoryJar CategoriesList

func ErrorHandler(w http.ResponseWriter, code int) {
	tmpl, err := template.ParseFiles("./web/UI/templates/error.html")
	if err != nil {
		fmt.Fprintf(w, "501, Internal server error: %s\n", err.Error())
		return
	}
	w.WriteHeader(code)
	tmpl.Execute(w, ErrorBar{ErrorCode: code, ErrorMsg: http.StatusText(code)})
}

func StartServer() {
	nerve := http.NewServeMux()
	nerve.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/UI/static/"))))
	nerve.HandleFunc("/", MainHandler)
	nerve.HandleFunc("/login/", LoginHandler)
	nerve.HandleFunc("/register/", RegisterHandler)
	nerve.HandleFunc("/loader/", UserLoadHandler)
	nerve.HandleFunc("/logout/", LogoutHandler)
	nerve.HandleFunc("/menu/", MenuHandler)
	nerve.HandleFunc("/api/credentials/", CheckCredentials)
	nerve.HandleFunc("/api/articles/", ArticlesHandler)
	nerve.HandleFunc("/api/filters/", FilterPostsApi)
	nerve.HandleFunc("/profile/", UserProfileHandler)
	sessionList = UID{
		QueryCount:       0,
		QueryHistory:     []string{"CREATION"},
		CurrentUserQueue: []string{"N/A"},
		InsertionTime:    []int64{time.Now().Unix()},
	}
	InitiateDatabase()
	loadErr := categoryJar.LoadAllData("./web/config/cats.txt")
	if loadErr != nil {
		log.Println(loadErr)
		return
	}
	fmt.Println("Start!")
	fmt.Println("Try to check the address http://localhost:8680/")
	startErr := http.ListenAndServe(":8680", nerve)
	if startErr != nil {
		log.Fatal(startErr.Error())
	}
}
