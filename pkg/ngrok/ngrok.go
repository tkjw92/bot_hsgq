package ngrok

import (
	"context"
	"log"
	"os"

	"golang.ngrok.com/ngrok/v2"
)

func Init() ngrok.EndpointListener {
	ngrok_agent, err := ngrok.NewAgent(ngrok.WithAuthtoken(os.Getenv("NGROK_AUTHTOKEN")))
	if err != nil {
		log.Fatal(err)
	}

	listener, err := ngrok_agent.Listen(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return listener
}
