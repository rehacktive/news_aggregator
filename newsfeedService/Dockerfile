FROM golang:1.13.0-stretch AS builder

ENV GO111MODULE=on

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o service service.go database.go handlers.go utils.go

WORKDIR /dist
RUN cp /build/service ./service


RUN mkdir /data

FROM scratch

COPY --chown=0:0 --from=builder /dist /

COPY --chown=65534:0 --from=builder /data /data
USER 65534
WORKDIR /data

ENTRYPOINT ["/service"]
