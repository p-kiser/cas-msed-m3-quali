curl -v -X PUT http://172.19.73.21:17468/tx  -H 'diva-token-api:TOKEN' -H 'Content-Type: application/json' --data-binary '[{"command":"data","ns":"test.i2p","d":"TESTDATA-HERE"}]'
