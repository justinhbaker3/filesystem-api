FROM golang:1.19-alpine

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o filesystem-api ./cmd/filesystem-api

EXPOSE 8080

ENV FILESYSTEM_API_DIRECTORY="/app"

CMD [ "./filesystem-api" ]