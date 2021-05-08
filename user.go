package main

import (
    "fmt"
    "reflect"

    "github.com/olivere/elastic/v7"
)

const (
    USER_INDEX = "user"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int64  `json:"age"`
    Gender   string `json:"gender"`
}

func checkUser(username, password string) (bool, error) { //check whether user exists
    query := elastic.NewBoolQuery() //select * from users where username = ? AND pw=?
	//boolquery => AND
    query.Must(elastic.NewTermQuery("username", username))
    query.Must(elastic.NewTermQuery("password", password))
    searchResult, err := readFromES(query, USER_INDEX)
    if err != nil {
        return false, err
    }

    var utype User
    for _, item := range searchResult.Each(reflect.TypeOf(utype)) { //TypeOf() cannot use classname User in (), only object
        u:= item.(User)
        if u.Password == password {
                fmt.Printf("Login as %s\n", username)
                return true, nil
        }
        
    }//return searchResult.TotalHit()>0
    return false, nil
}

func addUser(user *User) (bool, error) { //*: reference, no need to copy user object
    query := elastic.NewTermQuery("username", user.Username)
    searchResult, err := readFromES(query, USER_INDEX)
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 {
        return false, nil
    }

    err = saveToES(user, USER_INDEX, user.Username)
    if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}

