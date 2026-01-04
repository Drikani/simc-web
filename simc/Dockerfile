FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y \
      git \
      cmake \
      build-essential \
      libcurl4-openssl-dev \
      libssl-dev \
      zlib1g-dev \
      && rm -rf /var/lib/apt/lists/*

WORKDIR /opt

RUN git clone https://github.com/simulationcraft/simc.git
WORKDIR /opt/simc

RUN make -j$(nproc)

WORKDIR /data

ENTRYPOINT ["/opt/simc/engine/simc"]