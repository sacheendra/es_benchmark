FROM golang:1.5.3
MAINTAINER Sacheendra Talluri <sacheendra.t@gmail.com>

# Set GOPATH
ENV GOPATH /go

# Make directories for api_frontend
RUN mkdir -p /go/src/github.com/sacheendra/es_benchmark

# Add api_frontend files
ADD . /go/src/github.com/sacheendra/es_benchmark

# Define working directory
WORKDIR /go/src/github.com/sacheendra/es_benchmark

# Restore Dependencies and Install Application
RUN \
	cd /go/src/github.com/sacheendra/es_benchmark && \
	go install

# Define default command
CMD ["/go/src/github.com/sacheendra/es_benchmark/start.sh"]