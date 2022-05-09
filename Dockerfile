FROM golang:1.17

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN go build -o dolan main.go
RUN cp dolan /usr/local/bin

COPY views/ ./views/
ENTRYPOINT [ "/usr/local/bin/dolan" ]
