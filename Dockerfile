FROM golang:1.23
WORKDIR /app
COPY . .
## for sqlite3
ENV CGO_ENABLED=1

RUN go mod download
## fts5 is needed for sqlite full text search
RUN go build -tags "fts5" -o panda
CMD ["./panda"]
