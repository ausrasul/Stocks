FROM ubuntu:16.04
MAINTAINER ausrasul

# Package installer requirements
RUN apt-get update
RUN apt-get install -yq build-essential
RUN apt-get install -yq tcl8.5
RUN apt-get install -yq wget
# Install Redis server/client
RUN wget http://download.redis.io/releases/redis-stable.tar.gz
RUN tar xzf redis-stable.tar.gz
RUN cd redis-stable; make
RUN cd redis-stable; make test
RUN cd redis-stable; make install
RUN cd redis-stable; cd utils; ./install_server.sh
