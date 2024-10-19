
FROM golang:latest AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux 

RUN go build -o server ./cmd/server/main.go


RUN ls -la /app


FROM alpine:latest


WORKDIR /root/


COPY --from=builder /app/server .


RUN chmod +x ./server


VOLUME /root/logs


CMD ["./server"]
