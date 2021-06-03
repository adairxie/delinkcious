FROM golang:1.13 AS builder
ADD ./main.go main.go
ADD ./service service
# Fetch dependencies
RUN export GO111MODULE=on
RUN go get -d -v

# Build image as a truly static Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /news_service -a -tags netgo -ldflags '-s -w' .

FROM scratch
COPY --from=builder /news_service /app/news_service
EXPOSE 7070
ENTRYPOINT ["/app/news_service"]