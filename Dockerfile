#

FROM ffmpeg:builder as builder

WORKDIR /app

#COPY --from=ffmpeg /usr/local/lib/libav* /usr/local/lib/
#COPY --from=ffmpeg /usr/local/lib/libsw** /usr/local/lib/
#COPY --from=ffmpeg /usr/local/include/libav** /usr/local/include/
#COPY --from=ffmpeg /usr/local/include/libsw** /usr/local/include/

COPY . .

ENV CGO_LDFLAGS="-lstdc++ -lm -lcrypto -lssl -ldl -lavformat -lavcodec -lswscale -lavutil -lavfilter -lswresample -lavdevice -lx264"
ENV CGO_ENABLED=1

RUN go build -o tutorial example/tutorial.go

ENTRYPOINT ["/app/tutorial"]

# docker build -t avtest .
# docker run -it --rm -v $HOME:/mnt/home avtest /mnt/home/Desktop/vlc-output.ts
