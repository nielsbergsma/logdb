# Log based stream server / storage

Features:
- Stores every message ever received
- Can replay every message ever received
- Simple HTTP API (see below)
- Streams messages with http chunk encoding
- Content agnostic (handles json, xml, protobuf, ...)
- Based on Biglog (https://github.com/ninibe/netlog/tree/master/biglog)

Usage:
- Start: go run logdb.go --data-directory=data/ --http=:9999
- Insert message: `curl -X POST -H "Content-Type: application/json" -d '{ "message": 1 }' "http://localhost:9999/streams/stream1"`
- Stream messages (from offset): `curl "http://localhost:9999/streams/stream1?offset=0"`

//https://www.youtube.com/watch?v=ysjcEN548yc&t=1267s