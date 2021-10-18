FROM golang:1.17 as builder
COPY go.mod go.sum /go/src/github.com/aemiralfath/IH-Userland-Onboard/
WORKDIR /go/src/github.com/aemiralfath/IH-Userland-Onboard
RUN go mod download
COPY . /go/src/github.com/aemiralfath/IH-Userland-Onboard
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/IH-Userland-Onboard github.com/aemiralfath/IH-Userland-Onboard

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/aemiralfath/IH-Userland-Onboard/build/IH-Userland-Onboard /usr/bin/IH-Userland-Onboard
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/IH-Userland-Onboard"]