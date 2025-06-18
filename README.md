# protoc-gen-fasthttp

**`protoc-gen-fasthttp`** is a `protoc` plugin that generates [FastHTTP](https://github.com/valyala/fasthttp)-based HTTP server stubs from your `.proto` files. It uses `google.api.http` annotations to map RPCs to RESTful HTTP endpoints with high performance and minimal dependencies.

---

## ✨ Features

- Generate FastHTTP-compatible HTTP handlers from `.proto`
- Declarative routing via `google.api.http` annotations
- Ultra-low overhead, ideal for high-performance microservices
- Works with `protoc` and `buf`

---

## 🚀 Installation

```bash
go install github.com/streamdealer/protoc-gen-fasthttp@latest
```

Ensure `$GOPATH/bin` or `$GOBIN` is in your `PATH`.

---

## 📄 Example `.proto`

```proto
syntax = "proto3";

package helloworld;

import "google/api/annotations.proto";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {
    option (google.api.http) = {
      post: "/v1/hello"
      body: "*"
    };
  }
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

---

## ⚙️ Generate Code

```bash
protoc   --proto_path=.   --proto_path=third_party   --go_out=.   --go_opt=paths=source_relative   --fasthttp_out=.   --fasthttp_opt=paths=source_relative   helloworld.proto
```

---

## 🛠️ Implement Server

```go
type greeterServer struct{}

func (s *greeterServer) SayHello(ctx *fasthttp.RequestCtx, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
  return &helloworld.HelloReply{
    Message: "Hello, " + req.Name,
  }, nil
}
```

---

## 🚦 Setup and Serve

```go
package main

import (
  "log"

  "github.com/fasthttp/router"
  "github.com/valyala/fasthttp"
  "your_project/helloworld"
)

func main() {
  r := router.New()
  helloworld.RegisterGreeterHandler(r, &greeterServer{})

  log.Println("Listening on :8080")
  log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
```

---

## 📦 Buf Integration

Add this to `buf.gen.yaml`:

```yaml
version: v1
plugins:
  - name: go
    out: .
    opt: paths=source_relative
  - name: fasthttp
    out: .
    opt: paths=source_relative
```

Run:

```bash
buf generate
```

---

## 🧪 Examples

### ✅ cURL Request

```bash
curl -X POST http://localhost:8080/v1/hello \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'
```

### 🔁 Response

```json
{
  "message": "Hello, Alice"
}
```

---

## 📜 License

[MIT](LICENSE)

---

## 🤝 Contributing

Contributions and issues are welcome. Please open a PR or file a bug to get started.
