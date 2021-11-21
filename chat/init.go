package chat

import (
	"log"
	"os"
	"sync"

	stream "github.com/GetStream/stream-chat-go/v2"
)

type StreamClient struct {
	Stream *stream.Client
}

var (
	setupOnce sync.Once
	appStream *StreamClient
)

func ChatServer() (err error) {
	setupOnce.Do(func() {
		APIKey := os.Getenv("STREAM_API_KEY")
		APISecret := os.Getenv("STREAM_API_SECRET")
		client, err := stream.NewClient(APIKey, []byte(APISecret))
		if err != nil {
			log.Fatal(err)
		}
		appStream = &StreamClient{Stream: client}

		// Create admin user or update if admin already exists
		_, err = client.UpdateUser(&stream.User{
			ID:   "admin",
			Role: "admin",
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create the General channel
		_, err = client.CreateChannel("team", "general", "admin", map[string]interface{}{
			"name": "General",
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return
}

// ChatServerConn returns the global steam connection.
func ChatServerConn() *StreamClient {
	return appStream
}
