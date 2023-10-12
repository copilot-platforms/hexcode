FROM golang:1.19
WORKDIR /app
COPY . /app
RUN apt update
RUN apt install sqlite3
RUN go mod download
RUN go build -o hexcode

CMD ["./hexcode"]