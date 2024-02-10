FROM golang:1.20
WORKDIR /app
COPY . .
## sqlite3 needs gcc
ENV CGO_ENABLED=1

RUN go mod download
## fts5 is needed for full text search
RUN go build -tags "fts5" -o panda
CMD ["./panda"]
