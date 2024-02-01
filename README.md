# MSED Qualifiktationsschritt M3

Qualifikationsschritt für das Modul M3 CAS MSED HSLU.

## Spezifikation

Das Programm stellt via einem lokalen HTTP Server auf einem spezifischen Port zwei Endpunkte zur Verfügung:

- GET: /[domain-name]
- PUT: /[domain-name]/[b32-string]

Validation:

- domain-name: `^[a-z0-9-_]{3,64}\\.i2p$`
- b32-string: `^[a-z0-9]{52}$`

Beispiele:

- `GET /diva.i2p`
- `PUT /diva.i2p/auoqibfnyujhcht4v3nzahpqztwlyomesfywltuls5bqqi3nd3ka`

Detailspezifikation: https://gist.github.com/diva-exchange/aa6b1adbfefe909cd3ea07ac3cdfc322

## Run locally

Run locally with default values:

    go run main.go

Run locally different settings:

    PORT=8080 BASE_URL=http://127.19.73.21:17468 go run main.go

Environment variables:

- `PORT`: port that the application listens to, default `8080`
- `BASE_URL`: where to send the requests to, default is `http://127.19.73.21:17468`




