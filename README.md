# Go SDK for AWX X-Ray

> SDK for instrumenting Go applications for AWS X-Ray

## AWS X-Ray
More information about AWS X-Ray can be found at the [product website](https://aws.amazon.com/xray/).

## Getting Started

### Client Example
The outbound request is traced
```go
import "github.com/goguardian/aws-xray-go/xray"

func main() {
	ctx := xray.NewContext("service-name", nil)
	defer xray.Close(ctx)

	http := xray.GetHTTPClient(ctx)
	http.Get("http://127.0.0.1:3000")
}
```

### Server Example
All inbound requests will be traced
```go
import (
	"log"
	"net/http"
	"github.com/goguardian/aws-xray-go/xray"
)

func main() {
	http.HandleFunc("/", xray.Middleware("service-name", handler))
	http.ListenAndServe(":3000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hiya"))
}
```

### Sampling
The sampling configuration can be modified using `xray.SetSampler`.  The default settings are to sample all of the first ten requests of any second and five percent of all other requests that second.
```go
import "github.com/goguardian/aws-xray-go/xray"

func example() {
	fixedTarget := 10 // First ten requests per second
	fallbackRate := 0.05 // Five percent of all other requests
	xray.SetSampler(fixedTarget, fallbackRate)
}
