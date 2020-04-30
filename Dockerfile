FROM golang:1.14.1

ADD . /go/src/brain-debt

WORKDIR /go/src/brain-debt

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/brain-debt ./cmd/main.go

ENTRYPOINT ./out/brain-debt

EXPOSE 8080

CMD ["./out/brain-debt"]
