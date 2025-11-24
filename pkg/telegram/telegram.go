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
	const maxLength = 4096

	// Jika text tidak melebihi batas, kirim langsung
	if len(text) <= maxLength {
		cmd("sendMessage", map[string]any{
			"text":       text,
			"chat_id":    chat_id,
			"parse_mode": "Markdown",
		})
	}

	// Split text menjadi beberapa bagian
	chunks := splitMessage(text, maxLength)

	// Kirim setiap chunk
	for _, chunk := range chunks {
		cmd("sendMessage", map[string]any{
			"text":       chunk,
			"chat_id":    chat_id,
			"parse_mode": "Markdown",
		})
	}
}

func splitMessage(text string, maxLength int) []string {
	if len(text) <= maxLength {
		return []string{text}
	}

	var chunks []string
	runes := []rune(text)

	for len(runes) > 0 {
		// Tentukan panjang chunk yang akan diambil
		chunkSize := maxLength
		if len(runes) < chunkSize {
			chunkSize = len(runes)
		}

		// Cari posisi pemisah yang baik (newline atau space) jika memungkinkan
		if len(runes) > chunkSize {
			// Cari newline terdekat dari belakang
			breakPoint := chunkSize
			for i := chunkSize - 1; i >= chunkSize-200 && i >= 0; i-- {
				if runes[i] == '\n' {
					breakPoint = i + 1
					break
				}
			}

			// Jika tidak ada newline, cari space
			if breakPoint == chunkSize {
				for i := chunkSize - 1; i >= chunkSize-100 && i >= 0; i-- {
					if runes[i] == ' ' {
						breakPoint = i + 1
						break
					}
				}
			}

			chunkSize = breakPoint
		}

		// Ambil chunk dan tambahkan ke hasil
		chunks = append(chunks, string(runes[:chunkSize]))
		runes = runes[chunkSize:]
	}

	return chunks
}

func OnuMessageComposer(data snmp.Onu, chat_id int) {
	// Validasi panjang semua slice
	lengths := map[string]int{
		"Name":     len(data.Name),
		"Status":   len(data.Status),
		"Mac":      len(data.Mac),
		"Distance": len(data.Distance),
		"Rx":       len(data.Rx),
		"Tx":       len(data.Tx),
		"Vendor":   len(data.Vendor),
	}

	// Cek apakah semua panjang sama
	expectedLen := len(data.Name)
	for field, length := range lengths {
		if length != expectedLen {
			log.Fatalf("Array length mismatch! Name=%d, Status=%d, Mac=%d, Distance=%d, Rx=%d, Tx=%d. Field '%s' has different length",
				lengths["Name"], lengths["Status"], lengths["Mac"],
				lengths["Distance"], lengths["Rx"], lengths["Tx"], lengths["Vendor"], field)
		}
	}

	base_message := `
*Nama*: *%s*
*Status*: %s
*Mac*: %s
*ONU Vendor*: %s
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
			data.Vendor[i],
			data.Distance[i],
			data.Rx[i],
			data.Tx[i],
		))
	}
	SendMessage(strings.Join(result, ""), chat_id)
}

func OnuMessageList(data snmp.Onu, chat_id int) {
	// Validasi panjang semua slice
	lengths := map[string]int{
		"Name":   len(data.Name),
		"Status": len(data.Status),
		"Mac":    len(data.Mac),
	}

	// Cek apakah semua panjang sama
	expectedLen := len(data.Name)
	for field, length := range lengths {
		if length != expectedLen {
			log.Fatalf("Array length mismatch! Name=%d, Status=%d, Mac=%d. Field '%s' has different length",
				lengths["Name"], lengths["Status"], lengths["Mac"], field)
		}
	}

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
