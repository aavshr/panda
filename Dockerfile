FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod download
## sqlite3 needs gcc
RUN CGO_ENABLED=1 go build -o panda
CMD ["./panda"]
