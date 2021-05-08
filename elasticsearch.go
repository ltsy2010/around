package main

import (
    "context"

    "github.com/olivere/elastic/v7"
)

const (
        ES_URL = "http://10.128.0.2:9200"
)
func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    client, err := elastic.NewClient(
        elastic.SetURL(ES_URL),
        elastic.SetBasicAuth("elastic", "123456"))
    if err != nil {
        return nil, err
    }

    searchResult, err := client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

//interface{}: object
func saveToES(i interface{}, index string, id string) error{
    client, err := elastic.NewClient(
        elastic.SetURL(ES_URL),
        elastic.SetBasicAuth("elastic", "123456"))
    if err != nil {
        return err
    }

    _, err = client.Index(). //client.insert/update -> upsert
        Index(index).//into post
        Id(id). //user_id = xxx
        BodyJson(i). //user = xxx, url = xxx
        Do(context.Background())
    return err
}

