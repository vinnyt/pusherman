package main


import (
	"crypto/tls"
	"encoding/json"
	"flag"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	. "github.com/sideshow/apns2/payload"
	"io"
	"log"
	"net/http"
)

// internal push representation
type notification struct {
	token   string
	message string
	sound   string
	badge   int
}

// POST request to / with JSON body
// {
//   "tokens": [
//     "70455fc162e0577d9ff5f05737f5aaf091c64d864573f1db5a139e52e3a2b8ac"
//   ],
//   "message": "hello from remote",
//   "badge": 0,
//   "sound": "",
//   "extra": ""
// }
type PushRequest struct {
	Tokens  []string `json:"tokens"`
	Message string   `json:"message"`
	Badge   int      `json:"badge"`
	Sound   string   `json:"sound"`
	Extra   string   `json:"extra"`
}

func worker(input <-chan notification, cert tls.Certificate, topic string, production bool) {
	// TODO pass in environment
	client := apns.NewClient(cert)
	if production {
		client = client.Production()
	} else {
		client = client.Development()
	}

	log.Println(topic)

	// loop on the channel and send when something
	// is ready to send
	for messageParts := range input {
		log.Println("Sending message to this user: ", messageParts.token)

		payload := NewPayload()

		// construct outbound payload
		if len(messageParts.message) > 0 {
			payload.Alert(messageParts.message)
		}

		if messageParts.badge >= 0 {
			payload.Badge(messageParts.badge)
		}

		// prepare notification
		notification := &apns.Notification{}
		notification.DeviceToken = messageParts.token
		notification.Topic = topic
		notification.Payload = payload

		// send it
		res, err := client.Push(notification)

		log.Println("res:", res)

		if err != nil {
			log.Println("Error:", err)
			return
		}
	}
}

// handle web request this function returns a handler function
func requestHandler(messages chan notification) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// we parse the body of the request
		decoder := json.NewDecoder(r.Body)
		var t PushRequest
		err := decoder.Decode(&t)

		// it wasn't good
		if err != nil {
			io.WriteString(w, "bad")
		} else {
			// loop on tokens and write to channel for each token
			for _, token := range t.Tokens {
				log.Println("Token:", token)
				messages <- notification{token: token, message: t.Message, badge: t.Badge}
			}

			io.WriteString(w, "Hello world!")
		}
	}
}

func main() {
	topicPtr := flag.String("topic", "com.foo.bar", "topic on which to send messages typically your bundle identifier")
	environmentPtr := flag.Bool("production", false, "run in production mode")
	flag.Parse()
	// setup our buffered channel, we'll queue 1024 before we block
	messages := make(chan notification, 1024)

	// open the cert
	// should probably exit here
	cert, pemErr := certificate.FromPemFile("../cert.pem", "")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	// kick off some workers
	go worker(messages, cert, *topicPtr, *environmentPtr)
	go worker(messages, cert, *topicPtr, *environmentPtr)

	// http worker
	http.HandleFunc("/", requestHandler(messages))
	http.ListenAndServe(":8000", nil)
}
