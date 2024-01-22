# Stage: install dependencies and build code
# install and generate swagger docs
FROM golang:1.20-alpine3.19 as go-build
WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /rest.out ./cmd/rest/*

# Stage: serve
FROM alpine:3.16

RUN apk add tzdata
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /apps
COPY --from=go-build /rest.out ./rest

ENV GIN_MODE=release
EXPOSE 8080
CMD ["/apps/rest"]