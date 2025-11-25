# Bot Telegram HSGQ

## List Commands:

- Cari dan tampilkan data onu berdasarkan nama.\
  /name [name]

- Cari dan tampilkan data onu berdasarkan mac address.\
/mac [mac]

- List semua onu, menampilkan (nama, status, mac address)\
/list

# Installasi
Copy file <code>.env.example</code> dengan nama <code>.env</code>
Edit file .env
```
SNMP_ADDRESS= # IP Address OLT
SNMP_COMMUNITY= # SNMP Community OLT

NGROK_AUTHTOKEN= # Ngrok Auth Token
TELEGRAM_TOKEN= # Telegram Token

ALLOW_LIST_FILE=/app/operators
```

Edit file <code>operators</code>
File ini berisi username telegram untuk member yang di izinkan mengirimkan perintah pada bot.\
Jika username tidak terdaftar di file ini dan mencoba untuk mengirimkan perintah pada bot, maka bot tidak akan merespon apapun.
```
username1
username2
username3
```

```
docker compose up -d
```
