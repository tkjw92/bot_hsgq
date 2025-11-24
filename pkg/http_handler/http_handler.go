package http_handler

import (
	"BotHSGQ/pkg/snmp"
	"BotHSGQ/pkg/telegram"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type TelegramUpdate struct {
	UpdateId int `json:"update_id"`
	Message  struct {
		MessageID int    `json:"message_id"`
		Text      string `json:"text"`
		From      struct {
			Username string `json:"username"`
		} `json:"from"`
		Chat struct {
			Id int `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

const (
	helpCmd = `
List Commands:

Cari dan tampilkan data onu berdasarkan nama.
/name \[name]

Cari dan tampilkan data onu berdasarkan mac address.
/mac \[mac]

List semua onu, menampilkan (nama, status, mac address)
/list
	`
)

func commandHandler(c TelegramUpdate) {
	if string(c.Message.Text[0]) != "/" {
		telegram.SendMessage("Perintah tidak diketahui...", c.Message.Chat.Id)
		telegram.SendMessage(helpCmd, c.Message.Chat.Id)
		return
	}

	cmd := strings.Split(c.Message.Text, " ")
	if len(cmd) < 2 && cmd[0] != "/start" && cmd[0] != "/list" {
		telegram.SendMessage("Perintah tidak diketahui...", c.Message.Chat.Id)
		telegram.SendMessage(helpCmd, c.Message.Chat.Id)
		return
	}

	snmp_client := snmp.NewSNMPClient(os.Getenv("SNMP_ADDRESS"), os.Getenv("SNMP_COMMUNITY"))
	switch cmd[0] {
	case "/start":
		telegram.SetMenuButton(c.Message.Chat.Id)
		telegram.SendMessage(helpCmd, c.Message.Chat.Id)

	case "/name":
		res := snmp.GetByName(snmp_client, cmd[1])
		if len(res.Name) < 1 {
			telegram.SendMessage("Data tidak ditemukan...", c.Message.Chat.Id)
			return
		}

		telegram.OnuMessageComposer(res, c.Message.Chat.Id)

	case "/mac":
		res := snmp.GetByMac(snmp_client, cmd[1])
		if len(res.Mac) < 1 {
			telegram.SendMessage("Data tidak ditemukan...", c.Message.Chat.Id)
			return
		}

		telegram.OnuMessageComposer(res, c.Message.Chat.Id)

	case "/list":
		res := snmp.GetOnuList(snmp_client)
		fmt.Println(res)
		if len(res.Mac) < 1 {
			telegram.SendMessage("Data tidak ditemukan...", c.Message.Chat.Id)
			return
		}

		telegram.OnuMessageList(res, c.Message.Chat.Id)

	default:
		telegram.SendMessage("Perintah tidak diketahui...", c.Message.Chat.Id)
		telegram.SendMessage(helpCmd, c.Message.Chat.Id)
	}
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, "Only allowed post method")
		return
	}

	pre_response, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	recv_secret := r.Header["X-Telegram-Bot-Api-Secret-Token"]
	if len(recv_secret) < 1 {
		return
	}

	arbitary_secret := recv_secret[0]

	if arbitary_secret != telegram.SecretToken {
		return
	}

	var response TelegramUpdate
	json.Unmarshal(pre_response, &response)

	if response.Message.Text == "" {
		return
	}

	file, err := os.Open(os.Getenv("ALLOW_LIST_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	allowed := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == response.Message.From.Username {
			allowed = true
		}
	}

	if !allowed {
		return
	}

	commandHandler(response)
}
