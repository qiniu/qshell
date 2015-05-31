FROM ubuntu:14.04
#install software
RUN apt-get install -y wget
RUN apt-get install -y git
RUN apt-get install -y vim
RUN cd ~ && wget http://qdisk.qiniudn.com/go1.4.2.linux-amd64.tar.gz
RUN cd ~ && tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz

#set env variables
RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
RUN echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
RUN cd ~ && mkdir GoProjects
RUN echo 'export GOPATH=~/GoProjects' >> ~/.bashrc

#get remote source
RUN ["/bin/bash", "-c", "export GOPATH=~/GoProjects && /usr/local/go/bin/go get github.com/qiniu/api"]
RUN ["/bin/bash", "-c", "export GOPATH=~/GoProjects && /usr/local/go/bin/go get github.com/qiniu/api"]
RUN ["/bin/bash", "-c", "export GOPATH=~/GoProjects && /usr/local/go/bin/go get github.com/qiniu/rpc"]
RUN ["/bin/bash", "-c", "export GOPATH=~/GoProjects && /usr/local/go/bin/go get github.com/qiniu/log"]
RUN ["/bin/bash", "-c", "export GOPATH=~/GoProjects && /usr/local/go/bin/go get github.com/syndtr/goleveldb/leveldb"]
RUN ["/bin/bash", "-c", "export GOPATH=~/GoProjects && /usr/local/go/bin/go get github.com/yanunon/oss-go-api/oss"]
RUN ["/bin/bash", "-c", "export GOPATH=~/GoProjects && /usr/local/go/bin/go get github.com/golang/text"]
RUN mkdir -p ~/GoProjects/src/golang.org/x/text && cp -R ~/GoProjects/src/github.com/golang/text ~/GoProjects/src/golang.org/x
RUN mkdir ~/Projects
RUN ["/bin/bash", "-c", "cd ~/Projects && /usr/bin/git clone https://github.com/jemygraw/qshell"]