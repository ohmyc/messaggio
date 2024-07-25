FROM golang
WORKDIR /app
COPY . .
ENV GOPROXY=direct
RUN go mod download
RUN go build ./cmd/app

FROM busybox
COPY --from=0 /app/app ./app
CMD ["./app"]
