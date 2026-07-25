package handlers

import (
	"encoding/json"
	"context"
	"cloud.google.com/go/storage"
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