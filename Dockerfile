FROM golang:alpine

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/api ./main.go
EXPOSE 8080/tcp
CMD ["api"]