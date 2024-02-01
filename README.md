# MSED Qualifiktationsschritt M3

*Qualifikationsschritt für das Modul M3 CAS MSED HSLU.*

- Gallus Bühlmann <gallus.buehlmann@stud.hslu.ch>
- Pascal Kiser <pascal.kiser@stud.hslu.ch>

Detailspezifikation: https://gist.github.com/diva-exchange/aa6b1adbfefe909cd3ea07ac3cdfc322

## Run locally

Run locally with default values:

    go run main.go

Run locally different settings:

    PORT=8080 BASE_URL=http://127.19.73.21:17468 go run main.go

Environment variables:

- `PORT`: port that the application listens to, default `8080`
- `BASE_URL`: where to send the requests to, default is `http://127.19.73.21:17468`




