package main

import (
	"database/sql"
	"firstServer/service"
	"firstServer/storage"
	"flag"
	"fmt"
	mux2 "github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	HOST     = "localhost"
	DATABASE = "postgres"
	PORT     = "5432"
	USER     = "postgres"
	PASSWORD = "fs"
)

func main() {

	mux := mux2.NewRouter()
	srv := service.Service{Store: make(map[string]*storage.User)}
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/make_friends", srv.MakeFriends)
	mux.HandleFunc("/user", srv.Delete)
	mux.HandleFunc("/friends/{id:[0-9]+}", srv.GetFriends)
	mux.HandleFunc("/{user_id:[0-9]+}", srv.Put)

	mux.HandleFunc("/get_all", srv.GetAll)
	mux.HandleFunc("/get/{user_id:[0-9]+}", srv.GetUserInfo)

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, DATABASE)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully created connection to database")
	var newPort string
	flag.StringVar(&newPort, "port", "default_port", "The name of a port")

	flag.Parse()
	http.ListenAndServe(fmt.Sprintf("localhost:%s", newPort), mux)

	/*var dat storage.User
	for {
		sql_statement := "INSERT INTO users VALUES ($1,$2,$3);"
		_, err = db.Exec(sql_statement,
			dat.ID,
			dat.Name,
			dat.Age)
		if err != nil {
			log.Fatal(err)
		}

		for _, userFriend := range dat.Friends {
			var response []string
			response = append(response, strconv.Itoa(userFriend))
			sql_statement = " INSERT INTO users VALUES ($4);"
			_, err = db.Exec(sql_statement,
				response)
		}
	}*/
	for _, user := range srv.Store {
		response := "INSERT INTO users VALUES ($1,$2,$3);"
		_, err = db.Exec(response,
			user.ID,
			user.Name,
			user.Age)
		if err != nil {
			log.Fatal(err)
		}
	}

}
