package main

//type: image or video for frontend
//I/i: public/private

import (
    "mime/multipart"
    "reflect"

    "github.com/olivere/elastic/v7"
)

const (
    POST_INDEX  = "post"
)

type Post struct {
    Id      string `json:"id"`
    User    string `json:"user"`
    Message string `json:"message"`
    Url     string `json:"url"`
    Type    string `json:"type"`  
	
}
//success: return[] post, or return error
func searchPostsByUser(user string) ([]Post, error) {
    query := elastic.NewTermQuery("user", user) //select *from post where user =?..
    searchResult, err := readFromES(query, POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) {
    query := elastic.NewMatchQuery("message", keywords) //matchquery: message
    query.Operator("AND") //取交集if multiple keywords
    if keywords == "" {
        query.ZeroTermsQuery("all")
    }//return all posts if keywords is empty, default limit: 20
    searchResult, err := readFromES(query, POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
    var ptype Post
    var posts []Post

    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
        p := item.(Post) //cast Post to item 
        posts = append(posts, p)
    }
    return posts
}

//mulipart.file from http
func savePost(post *Post, file multipart.File) error {
  mediaLink, err := saveToGCS(file, post.Id,)
  if err != nil {
      return err
  }
  post.Url = mediaLink
  return saveToES(post, POST_INDEX, post.Id)
}