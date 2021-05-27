FROM golang:1.11 AS builder
ADD ./main.go main.go
ADD ./service service
# Fetch dependencies
RUN go get -d -v

# Build image as a truly static Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /user_service -a -tags netgo -ldflags '-s -w' .

FROM scratch
COPY --from=builder /user_service /app/user_service
EXPOSE 7070
ENTRYPOINT ["/app/user_service"]
