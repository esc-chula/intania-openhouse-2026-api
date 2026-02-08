# Inspired from https://github.com/GoogleContainerTools/distroless/blob/main/examples/go/Dockerfile

FROM golang:1.24 AS build

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /

EXPOSE 8000

CMD ["/app", "serve"]
