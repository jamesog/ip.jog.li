FROM golang:1.18 AS build

WORKDIR /go/src/github.com/jamesog/ip.jog.li
COPY . .

RUN go get -d -v ./... && \
	CGO_ENABLED=0 go install -a -tags netgo -ldflags '-w' -v ./cmd/ip.jog.li && \
	echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch
COPY --from=build /go/bin/ip.jog.li /ip.jog.li
COPY --from=build /etc_passwd /etc/passwd

USER nobody

CMD ["/ip.jog.li"]
