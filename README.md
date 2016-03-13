# pusherman
a golang based microweb service API for sending Push Notifications.  The goal of pusherman is to provide a webservice that can be communicated with to easily send push notifications.  The webservice should be responsible for things like queuing and maintenance of persistent connections as dictated by APNS and eventually the Android equivalent.  Currently only APNS is supported

## Install

1. `go get -u golang.org/x/net/http2`
2. `go get -u golang.org/x/crypto/pkcs12`
3. `go get -u github.com/sideshow/apns2`


# Getting started

1. start the server

```
go run main.go -topic=<BUNDLE ID/ or topic> -production
```

2. Send a notification
```
curl -H "Content-Type: application/json" -X POST -d '{"tokens":["70455fc162e0577d9ff5f05737f5aaf091c64d864573f1db5a139e52e3a2b8ac"],"message":"hello from remote","badge":0,"sound":"","extra":""}' http://localhost:8000/
```

# Feedback and getting Involved
My hope is this can evolve to at least provide support for Android Push.

Please file issues and pull requests

My design philosophy was to create an easily configurable and robustly running service that could be launched with minimal systems requirements and easily communicated with to support push requirements for various applications.
