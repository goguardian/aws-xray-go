package main

import "github.com/goguardian/aws-xray-go/xray"

func main() {
	ctx := xray.NewContext("http-client", nil)
	defer xray.Close(ctx)

	http, err := xray.GetHTTPClient(ctx)
	if err != nil {
		panic(err)
	}

	http.Get("http://127.0.0.1:3000")
}
