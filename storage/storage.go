package storage

import "fmt"

type User struct {
	ID      int    //Users ID's
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

type FriendsMaker struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type NewAge struct {
	NewAge int `json:"new_age"`
}

func (u *User) ToString() string {
	return fmt.Sprintf("id %d, name %s, age %d, friends %v\n", u.ID, u.Name, u.Age, u.Friends)
}
