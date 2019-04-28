package main

import (
	"fmt"
	"html/template"
	"model"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("tmpl/login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]

		user := model.Users{Username: username, Password: password}
		row, err := user.Check()
		fmt.Println(row)
		fmt.Println(err)
	}
}
