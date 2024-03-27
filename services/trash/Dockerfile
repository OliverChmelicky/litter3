FROM golang:1.22 as builder

# Set destination for COPY
WORKDIR /app

COPY . .
# Download Go modules
RUN go mod download

WORKDIR /app/cmd
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /service

# Run
FROM scratch
COPY --from=0 /service /service
ENTRYPOINT [ "/service" ]