FROM golang:1.13 AS build
WORKDIR /src
# Copy go.mod and go.sum to download all dependencies (this is cached if go.mod and go.sum don't change)
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy everything else to actually build the project
COPY . .
# Build project to build folder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64\
 go build -ldflags="-w -s"\
 -o build/test-server

FROM scratch
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/build/ /app/
WORKDIR /app
ENTRYPOINT ["/app/test-server"]
