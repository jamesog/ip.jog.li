FROM golang:1.9 AS build

WORKDIR /go/src/ip.jog.li
COPY main.go .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -v ./... .

FROM scratch
COPY --from=build /go/src/ip.jog.li/ip.jog.li /ip.jog.li

EXPOSE 8000

CMD ["/ip.jog.li"]
