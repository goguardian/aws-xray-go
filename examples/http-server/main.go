package main

import (
	"github.com/goguardian/aws-xray-go/xray"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	listenAddr = ":3000"
	serverName = "http-server"
)

func main() {
	http.HandleFunc("/", xray.Middleware(serverName, handler))

	log.Println("Listening on:", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	http, err := xray.GetHTTPClient(r.Context())
	if err != nil {
		panic(err)
	}

	http.Get("http://www.example.com")

	session, _ := session.NewSession()
	s3Client := s3.New(session, &aws.Config{
		Region:     aws.String("us-west-2"),
		HTTPClient: http,
	})
	s3Client.ListBuckets(&s3.ListBucketsInput{})

	w.Write([]byte("hiya"))
}
