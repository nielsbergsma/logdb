# Log based stream server

Features:
- Stores every message ever received
- Streams messages with http chunk encoding
- No content restriction (handles json, xml, protobuf, ...)
- Simple HTTP API (see below)
- Based on Biglog (https://github.com/ninibe/netlog/tree/master/biglog)

Usage:
- Start: go run logdb.go --data-directory=data/ --http=:9999
- Insert message: `curl -X POST -H "Content-Type: application/json" -d '{ "message": 5 }' "http://localhost:9999/streams/stream1"`
- Stream messages (from offset): `curl "http://localhost:9999/streams/stream1?offset=0`
