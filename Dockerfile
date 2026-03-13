# Inspired from https://github.com/GoogleContainerTools/distroless/blob/main/examples/go/Dockerfile

# --- Stage 1: Build ---
FROM golang:1.24 AS build

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/app

# --- Stage 2: Normal Mode (Production) ---
FROM gcr.io/distroless/static-debian12 AS normal

COPY --from=build /go/bin/app /

EXPOSE 8000

CMD ["/app", "serve"]

# --- Stage 3: Testing Mode ---
FROM normal AS testing

# Add the specific testing file
COPY --from=build /go/src/service-account.json /
