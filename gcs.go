package main

import (
    "context"
    "fmt"
    "io"

    "cloud.google.com/go/storage"
)

const (
    BUCKET_NAME = "sheryl-bucket"
)

func saveToGCS(r io.Reader, objectName string) (string, error) {
    //empty context
	ctx := context.Background()
    //set connection
    client, err := storage.NewClient(ctx)
    if err != nil {
        return "", err
    }
    //copy local file r to remote wc
    object := client.Bucket(BUCKET_NAME).Object(objectName)
    wc := object.NewWriter(ctx)
    if _, err := io.Copy(wc, r); err != nil {
        return "", err
    }

    if err := wc.Close(); err != nil {
        return "", err
    }

	//AllUsers: serviceaccount check frontend account info
    if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
        return "", err
    }

    attrs, err := object.Attrs(ctx)
    if err != nil {
        return "", err
    }

    fmt.Printf("Image is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil
}
