FROM golang:1.14 AS builder

ENV FFMPEG_VERSION 4.2.2
ENV NASM_VERSION=2.14.02
ENV YASM_VERSION=1.3.0
ENV PREFIX=/usr/local

WORKDIR /buildenv

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
          autoconf \
          automake \
          build-essential \
          git-core \
          libass-dev \
          libfreetype6-dev \
          libsdl2-dev \
          libtool \
          libva-dev \
          libvdpau-dev \
          libvorbis-dev \
          libxcb1-dev \
          libxcb-shm0-dev \
          libxcb-xfixes0-dev \
          pkg-config \
          texinfo \
          wget \
          zlib1g-dev \
          libx264-dev \
          libmp3lame-dev \
          libopus-dev \
          libvpx-dev \
          libssl-dev \
          yasm \
          nasm

RUN mkdir -p ffmpeg_sources \
    && cd /buildenv/ffmpeg_sources \
    && wget -O ffmpeg.tar.gz "https://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz" \
    && tar -xzvf ffmpeg.tar.gz \
    && cd ffmpeg-${FFMPEG_VERSION} \
    && ./configure --prefix="${PREFIX}" \
    --pkg-config-flags="--static" \
    --prefix=$PREFIX \
    --extra-cflags="-I$PREFIX/include" \
    --extra-ldflags="-L$PREFIX/lib" \
    --extra-libs="-lpthread -lm" \
    --enable-gpl \
    --enable-libass \
    --enable-libfreetype \
    --enable-libmp3lame \
    --enable-libopus \
    --enable-libvorbis \
    --enable-libvpx \
    --enable-libx264 \
    --enable-nonfree \
    --enable-libsrt \
    && make \
    && make install
