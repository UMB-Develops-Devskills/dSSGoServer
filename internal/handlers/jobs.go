package handlers

import (
	"encoding/json"
	"context"
	"strings"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

const bucketName = "devskills-499815.firebasestorage.app"

// this function downloads the job postings from firebase as a jobfile and stores them into []job
func DownloadJobs(ctx context.Context, client *storage.Client, bucketName string, filepath string) ([]Job, error) {
	bkt := client.Bucket(bucketName) 
	object := bkt.Object(filepath)

	reader, err := object.NewReader(ctx) 
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var jf JobFile
	err = json.NewDecoder(reader).Decode(&jf)
	if err != nil {
		return nil, err
	}
	return jf.Jobs, nil
}

// this function uploads file to the given filepath
func UploadFile(ctx context.Context, client *storage.Client, bucketName, filepath string, data []byte) error {
	bkt := client.Bucket(bucketName) 
	object := bkt.Object(filepath)
	w := object.NewWriter(ctx)
  w.ContentType = "application/json"
  _, err := w.Write(data)
	if err != nil {
    return err
  }
  return w.Close()
}

// this function retrieves job files 
func GetJobFiles(ctx context.Context, client *storage.Client, bucketName string, prefix string) ([]string, error) {
	var files []string
	bkt := client.Bucket(bucketName)
	query := &storage.Query{Prefix: prefix}
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		// ignore folders
    if strings.HasSuffix(attrs.Name, "/") {
        continue
    }
    // only files
    if !strings.HasSuffix(attrs.Name, "-jobs.json") {
        continue
    }
		files = append(files, attrs.Name)
	}
	return files, nil
}
