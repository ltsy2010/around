package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"
    "regexp"
    "time"

    jwt "github.com/form3tech-oss/jwt-go"
    "github.com/pborman/uuid"
)

var (
    mediaTypes = map[string]string{
        ".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
    }
)

// Parse from body of request to get a json object.
//for debug
//construct a post object //if *post, if(p), not &p
//convert body to p object, &p: reference, change p directly. 
//try catch block, if error != null, panic..->throw exception
//print to w:responsebody. 

var mySigningKey = []byte("secret") //byte array 

func uploadHandler(w http.ResponseWriter, r *http.Request) {
   
    fmt.Println("Received one post request") 

	w.Header().Set("Access-Control-Allow-Origin", "*") //allow cross region 
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

    if r.Method == "OPTIONS" {
        return
    }
    
    user := r.Context().Value("user") //user: rawTokenstring, get token from user, request header里user key stores token.
    claims := user.(*jwt.Token).Claims //claims: payload
    username := claims.(jwt.MapClaims)["username"]
    //user.(*jwt.Toekn): cast/type assertion

    p := Post{
        Id: uuid.New(),
        User: username.(string),
        Message: r.FormValue("message"),
    }

    file, header, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    suffix := filepath.Ext(header.Filename)
    if t, ok := mediaTypes[suffix]; ok {
        p.Type = t
    } else {
        p.Type = "unknown"
    }

    err = savePost(&p, file)
    if err != nil {
        http.Error(w, "Failed to save post to GCS or Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save post to GCS or Elasticsearch %v\n", err)
        return
    }

    fmt.Println("Post is saved successfully.")

}
//responsewriter: interface, request: struct-pointer
func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for search")

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
    w.Header().Set("Content-Type", "application/json")

    if r.Method == "OPTIONS" {
        return
    }

    user := r.URL.Query().Get("user") //info after ?
    keywords := r.URL.Query().Get("keywords")

    var posts []Post
    var err error
    if user != "" {
        posts, err = searchPostsByUser(user)
    } else {
        posts, err = searchPostsByKeywords(keywords)
    }

    if err != nil {
        http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read post from Elasticsearch %v.\n", err)
        return
    }

    //return json to frontend(go/post->json: marshal)
    js, err := json.Marshal(posts)
    if err != nil { //add error to response w
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js) //wrtie js to reponse body
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signin request")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        return
    }
    //post
    //  Get User information from client
    decoder := json.NewDecoder(r.Body)
    var user User
    //&user, change on user directly
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

    exists, err := checkUser(user.Username, user.Password)
    if err != nil {
        http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
        return
    }

    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }
    //Claim => payload on the website
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }) //expire in 24 hours, convert to unix time. 

    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }

    w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        return
    }

    decoder := json.NewDecoder(r.Body)
    var user User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }
    //^: start, $; end
    if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

    success, err := addUser(&user)
    if err != nil {
        http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
        return
    }

    if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
    }
    fmt.Printf("User added successfully: %s.\n", user.Username)
}
//signouthandler: no need to do anythign when receiving log out from frontend, session needs to be destroyed. 
