package telegram

import (
	"BotHSGQ/pkg/snmp"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var SecretToken string

func cmd(endpoint string, data map[string]any) (response map[string]any) {
	telegram_token := os.Getenv("TELEGRAM_TOKEN")
	if telegram_token == "" {
		log.Fatal("Can't load telegram auth token")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", telegram_token, endpoint)
	r, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	return response
}

func Init(webhook_url string) bool {
	payload := map[string]any{
		"url":                  webhook_url,
		"drop_pending_updates": true,
		"allowed_updates":      []string{"message"},
		"secret_token":         SecretToken,
	}

	response := cmd("setWebhook", payload)

	if ok, exists := response["ok"].(bool); exists && ok {
		return true
	}

	return false
}

func SendMessage(text string, chat_id int) {
	cmd("sendMessage", map[string]any{
		"text":       text,
		"chat_id":    chat_id,
		"parse_mode": "Markdown",
	})
}

func OnuMessageComposer(data snmp.Onu, chat_id int) {
	base_message := `
*Nama*: *%s*
*Status*: %s
*Mac*: %s
*Distance*: %s m
*Rx*: %s dBm
*Tx*: %s dBm
-------------------------------

`

	result := []string{}

	for i := range data.Name {
		status := data.Status[i]
		if status == "up" {
			status = "✅ Up"
		} else {
			status = "❌ Down"
		}

		result = append(result, fmt.Sprintf(
			base_message,
			data.Name[i],
			status,
			strings.ToUpper(data.Mac[i]),
			data.Distance[i],
			data.Rx[i],
			data.Tx[i],
		))
	}

	SendMessage(strings.Join(result, ""), chat_id)
}

func OnuMessageList(data snmp.Onu, chat_id int) {
	base_message := `
*Nama*: *%s*
*Status*: %s
*Mac*: %s

`

	overview_message := `
-------------------------------
👤 Total: %d
✅ Online: %d
❌ Offline: %d
`

	online := 0
	offline := 0

	result := []string{}

	for i := range data.Name {
		status := data.Status[i]
		if status == "up" {
			status = "✅ Up"
			online++
		} else {
			status = "❌ Down"
			offline++
		}

		result = append(result, fmt.Sprintf(
			base_message,
			data.Name[i],
			status,
			strings.ToUpper(data.Mac[i]),
		))
	}

	result = append(result, fmt.Sprintf(overview_message, (online+offline), online, offline))

	SendMessage(strings.Join(result, ""), chat_id)
}

func setMyCommands() {
	cmd("setMyCommands", map[string]any{
		"commands": []map[string]string{
			{
				"command":     "start",
				"description": "Start bot",
			},
			{
				"command":     "name",
				"description": "/name [name], Cari dan tampilkan data onu berdasarkan nama.",
			},
			{
				"command":     "mac",
				"description": "/mac [mac], Cari dan tampilkan data onu berdasarkan mac address.",
			},
			{
				"command":     "list",
				"description": "List All Onu, menampilkan (nama, status, mac address)",
			},
		},
	})
}

func SetMenuButton(chat_id int) {
	setMyCommands()

	cmd("setChatMenuButton", map[string]any{
		"chat_id": chat_id,
	})
}
