package firebase

import (
  "context"
	"cloud.google.com/go/storage"
)

// this function returns the connection client to firebase storage
// BTW Clients should be reused instead of created as needed
func ConnectFireBaseStorage(ctx context.Context) (*storage.Client, error) {
	/*client, err := storage.NewClient(ctx) 
	if err != nil {
		return nil, nil, fmt.Errorf("could not create google cloud storage client %w", err)
	}*/
	return storage.NewClient(ctx)
}