package main

import (
	"firstServer/service"
	"firstServer/storage"
	mux2 "github.com/gorilla/mux"
	"net/http"
)

func main() {

	mux := mux2.NewRouter()
	srv := service.Service{Store: make(map[string]*storage.User)}
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/make_friends", srv.MakeFriends)
	mux.HandleFunc("/user", srv.Delete)
	mux.HandleFunc("/friends/{id:[0-9]+}", srv.GetFriends)
	mux.HandleFunc("/{user_id:[0-9]+}", srv.Put)

	//my handler for debug
	mux.HandleFunc("/get_all", srv.GetAll)
	mux.HandleFunc("/get/{user_id:[0-9]+}", srv.GetUserInfo)

	http.ListenAndServe("localhost:8080", mux)
}
