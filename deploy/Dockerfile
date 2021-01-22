FROM alpine:3.12.1 as runner
RUN apk add ffmpeg

FROM runner
LABEL Description="This image is used to host the Audio len service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY data/audio-convert /app/
COPY config.yaml /app/

ENTRYPOINT ["./audio-convert"]
