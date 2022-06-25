FROM golang:1.18.3-alpine as build
RUN apk add --update --no-cache build-base
WORKDIR /go/src/coffee-cup-counter
COPY *.go go.mod go.sum /go/src/coffee-cup-counter/

RUN go test -race

RUN CGO_ENABLED=0 go build -o main .

FROM scratch
WORKDIR /go/src/coffee-cup-counter/
COPY --from=build /go/src/coffee-cup-counter/main .
EXPOSE 8080
CMD ["./main"]