#########################################################################################
## docker will invoke this file from ../ dir in order to access the code
#########################################################################################
FROM golang:1.17-alpine as builder

ARG BUILD_VERSION=0.1

WORKDIR /go/src/
ENV CGO_ENABLED=0

# precache go.mod dependencies
COPY go.mod /go/src/
COPY go.sum /go/src/
RUN go mod download

# copy src
COPY . /go/src

RUN CGO_ENABLED=0 go build -o /go/bin/audio-convert -ldflags "-X main.version=$BUILD_VERSION" cmd/audio-convert/main.go
#####################################################################################
FROM alpine:3.13.8 as runner

ARG BUILD_VERSION=0.1

LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="audio-convert-service" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is used to host the Audio convert service" \
      description="This image is used to host the Audio convert service" 

RUN apk add ffmpeg

WORKDIR /app

ENV UID=999
EXPOSE 8000

RUN addgroup app && adduser -S -G app app --uid $UID app

RUN mkdir -p /app && chown -R app:app /app   

COPY --from=builder /go/bin/audio-convert /app/
COPY build/config.yaml /app/

RUN chown app:app /app/* /app
 
USER app

ENTRYPOINT ["./audio-convert"]
