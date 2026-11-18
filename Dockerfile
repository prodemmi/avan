FROM golang:1.24.0-alpine as GOBUILDER

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod tidy && go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN apk add --no-cache make bash

RUN make sqlc
RUN make templ

RUN go build -o avan cmd/server/main.go

FROM node:alpine as CSSBUILDER

WORKDIR /app

COPY package*.json .
COPY ./public ./public

RUN npm install
RUN npm run build

FROM golang:1.24.0-alpine as FINAL

WORKDIR /avan

COPY --from=GOBUILDER /app/avan .
COPY --from=GOBUILDER /app/internal/db/migrations/postgresql ./internal/db/migrations/postgresql/
COPY --from=CSSBUILDER /app/public/ ./public/

RUN chmod +x ./avan

ENTRYPOINT ["sh", "-c", "./avan --running-mode=${RUNNING_MODE}"]
