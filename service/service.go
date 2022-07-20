package service

import (
	"encoding/json"
	"firstServer/storage"
	"fmt"
	mux2 "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

//var GlobalId int

type Service struct {
	Store map[string]*storage.User
}

func (s *Service) GetUser(id int) *storage.User {
	for k, v := range s.Store {
		if v.ID == id {
			return s.Store[k]
		}
	}
	return nil
}

func (s *Service) DeleteByID(id int) {
	for _, v := range s.Store {
		for k, u := range v.Friends {
			if u == id {
				v.Friends = append(v.Friends[:k], v.Friends[k+1:]...)
			}
		}
	}
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer r.Body.Close()

		var u *storage.User

		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		s.Store[u.Name] = u
		for _, v := range s.Store {
			if v.ID == len(s.Store)-1 {
				v.ID++
				u.ID = v.ID
			}
		}
		//GlobalId += 1
		//u.ID = GlobalId //костылёк с ID

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("User id:%d\n", u.ID)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		defer r.Body.Close()

		var fm *storage.FriendsMaker

		if err := json.Unmarshal(content, &fm); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		userFirst := s.GetUser(fm.SourceId)
		userSecond := s.GetUser(fm.TargetId)
		userFirst.Friends = append(userFirst.Friends, userSecond.ID)
		userSecond.Friends = append(userSecond.Friends, userFirst.ID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья\n", userFirst.Name, userSecond.Name)))

		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		defer r.Body.Close()

		var ud *storage.FriendsMaker

		if err := json.Unmarshal(content, &ud); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		user := s.GetUser(ud.TargetId)
		_, ok := s.Store[user.Name] //проверка наличия имени
		if ok {
			delete(s.Store, user.Name)
			s.DeleteByID(user.ID)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s\n", user.Name)))

		return
	}
}

func (s *Service) GetFriends(w http.ResponseWriter, r *http.Request) {

	vars := mux2.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if r.Method == http.MethodGet {
		response := ""

		user := s.GetUser(id)

		//проверка наличия имени
		if user == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с таким id"))
			return
		}

		for _, userFriend := range user.Friends {
			response += fmt.Sprintf("%s\n", s.GetUser(userFriend).Name) //get username by id
		}
		if response == "" {
			response = fmt.Sprintf("У %s нет друзей :(\n", user.Name)
		} else {
			response = fmt.Sprintf("Друзья %s:\n%s", user.Name, response)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) Put(w http.ResponseWriter, r *http.Request) {

	vars := mux2.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if r.Method == http.MethodPut && r.Header.Get("Content-Type") == "application/json" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer r.Body.Close()

		var na *storage.NewAge

		if err = json.Unmarshal(content, &na); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Println(string(content))
		user := s.GetUser(userId)

		if user == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с таким id"))
			return
		}

		user.Age = na.NewAge
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint("Возраст успешно обновлён\n")))

		return
	}
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		response := ""
		for _, user := range s.Store {
			response += user.ToString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetUserInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux2.Vars(r)
	id, err := strconv.Atoi(vars["user_id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if r.Method == http.MethodGet {
		user := s.GetUser(id)
		//проверка наличия имени
		if user == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с таким id"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.ToString()))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
