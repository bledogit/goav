FROM ffmpeg:builder

WORKDIR /app

RUN apt-get update -y \
    && apt-get install -y \
          build-essential \
          ca-certificates \
          g++ \
          gcc \
          libc-dev \
          make \
          cmake \
          libssl-dev \
          autoconf

COPY . .

ENV CGO_LDFLAGS="-lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter"
ENV CGO_LDFLAGS="-lstdc++ -lm -lcrypto -lssl -ldl -lavformat -lavcodec -lswscale -lavutil -lavfilter -lswresample -lavdevice -lx264"
ENV CGO_ENABLED=1

RUN go build -o tutorial example/tutorial.go

ENTRYPOINT ["/app/tutorial"]

# docker build -t avtest .
# docker run -it --rm -v $HOME:/mnt/home avtest /mnt/home/Desktop/vlc-output.ts
