FROM golang:1.10.2-alpine3.7 AS build

WORKDIR /go/src/github.com/bitfield/edger/
COPY *.go go.* /go/src/github.com/bitfield/edger/
RUN CGO_ENABLED=0 go test
RUN CGO_ENABLED=0 go build -o /bin/edger

FROM scratch
COPY --from=build /bin/edger /bin/edger
ENTRYPOINT ["/bin/edger"]
