FROM centos:7

ENV GOVERSION 1.10.1

## install git
RUN yum update -y && yum install wget git make gcc gcc-c++ kernel-devel redis -y
RUN git config --global user.name "IOST" && git config --global user.email "chain@iost.io"

# Install Redis.
RUN \
  cd /tmp && \
  wget http://download.redis.io/redis-stable.tar.gz && \
  tar xvzf redis-stable.tar.gz && \
  cd redis-stable && \
  make && \
  make install && \
  cp -f src/redis-sentinel /usr/local/bin && \
  mkdir -p /etc/redis && \
  cp -f *.conf /etc/redis && \
  rm -rf /tmp/redis-stable* && \
  sed -i 's/^\(bind .*\)$/# \1/' /etc/redis/redis.conf && \
  sed -i 's/^\(daemonize .*\)$/# \1/' /etc/redis/redis.conf && \
  sed -i 's/^\(logfile .*\)$/# \1/' /etc/redis/redis.conf

EXPOSE 6379

## install go
RUN mkdir /goroot && \
    mkdir /gopath && \
    curl https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz | \
         tar xzf - -C /goroot --strip-components=1

ENV CGO_ENABLED 1
ENV GOPATH /gopath
ENV GOROOT /goroot
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

# Install Python
RUN yum install -y epel-release
RUN yum install -y python python-devel python-pip

# Install project
RUN mkdir -p $GOPATH/src/github.com/iost-official && cd $GOPATH/src/github.com/iost-official && \
git clone https://445789ea93ff81d814c78fccae8e25000f96e539@github.com/iost-official/prototype && \
cd prototype && git checkout testnet && go get github.com/kardianos/govendor && govendor sync -v && \
pip install -r scripts/backup/requirements.txt && \
cd iserver && go install && cd ../imonitor && go install && mkdir /workdir

EXPOSE 30301
EXPOSE 30302
EXPOSE 30303
EXPOSE 30304
EXPOSE 30305
EXPOSE 30306
EXPOSE 30307
EXPOSE 30308
EXPOSE 30309
EXPOSE 30310

WORKDIR /workdir

## docker deploy
## docker build -t iost .
## docker run --name iost_container -p 30301:30301 -p 30303:30303 -p 30310:30310 -p 8080:8080 -v /home/ec2-user/workdir:/workdir  -d iost imonitor
## sudo docker run --name iost_container1 -p 30302:30302 -p 30304:30304 -p 30309:30310 -v /home/ec2-user/workdir1:/workdir  -d iost imonitor
## sudo docker run --name iost_container2 -p 30305:30305 -p 30307:30307 -p 30308:30310 -v /home/ec2-user/workdir2:/workdir  -d iost imonitor
