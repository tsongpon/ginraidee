FROM golang:1.12

COPY . /ginraidee
WORKDIR /ginraidee

ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build -o ginraidee

EXPOSE 5000
CMD ["./ginraidee"]