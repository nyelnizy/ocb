FROM golang:1.17-alpine3.15 as build

WORKDIR /app
COPY ./go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o authms .

FROM alpine:3.15
WORKDIR /app
COPY --from=build /app .
EXPOSE 9200
CMD ["/app/authms"]