# CAS MSED M3 Qualifikationsschritt

Qualifikationsschritt für das Modul M3 CAS MSED HSLU.

Detailspezifikation: https://gist.github.com/diva-exchange/aa6b1adbfefe909cd3ea07ac3cdfc322

## API Endpunkte

    GET: /[domain-name]
    PUT: /[domain-name]/[b32-string]

## Validierung


    GET /[a-z0-9-_]{3-64}\.i2p$
    PUT /[a-z0-9-_]{3-64}\.i2p$/[a-z0-9]{52}$

## Beispiel

    GET /diva.i2p
    PUT /diva.i2p/auoqibfnyujhcht4v3nzahpqztwlyomesfywltuls5bqqi3nd3ka
