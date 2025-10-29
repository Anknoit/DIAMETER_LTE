FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN apk add --no-cache git && go env -w GOPATH=/go && go build -o /bin/diameter ./cmd/server

FROM alpine:3.18
COPY --from=build /bin/diameter /usr/local/bin/diameter
COPY config /etc/diameter
EXPOSE 3868 8080
CMD ["/usr/local/bin/diameter"]
