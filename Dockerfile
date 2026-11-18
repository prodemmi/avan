# syntax=docker/dockerfile:1
FROM golang:1.26.0-alpine as GOBUILDER

ENV GOTOOLCHAIN local
WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn

RUN go tool templ generate

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o avan cmd/server/main.go

FROM node:alpine as CSSBUILDER

WORKDIR /app

COPY package*.json .
COPY ./public ./public
COPY ./views ./views

RUN --mount=type=cache,target=/root/.npm npm install
RUN npm run build

FROM golang:1.26.0-alpine as FINAL

WORKDIR /avan

COPY --from=GOBUILDER /app/avan .
COPY --from=GOBUILDER /app/internal/db/migrations/postgresql ./internal/db/migrations/postgresql/
COPY --from=CSSBUILDER /app/public/ ./public/

RUN chmod +x ./avan

ENTRYPOINT ["sh", "-c", "./avan --running-mode=${RUNNING_MODE}"]
