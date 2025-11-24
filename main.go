package main

import (
	"BotHSGQ/pkg/http_handler"
	"BotHSGQ/pkg/ngrok"
	"BotHSGQ/pkg/telegram"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listener := ngrok.Init()
	webhook_url := fmt.Sprintf("%v", listener.URL()) + "/webhook"

	hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("TELEGRAM_TOKEN")), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hash_string := string(hash)
	re := regexp.MustCompile(`[^A-Za-z0-9\-_]`)
	telegram.SecretToken = re.ReplaceAllString(hash_string, "")

	init_webhook := telegram.Init(webhook_url)
	if !init_webhook {
		log.Fatal("Can't set webhook")
	}

	router := http.NewServeMux()
	router.HandleFunc("/webhook", http_handler.Webhook)

	fmt.Println("Bot is ready...")

	http.Serve(listener, router)
}
