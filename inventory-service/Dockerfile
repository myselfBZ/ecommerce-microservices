FROM golang:lates 
WORKDIR . ./app/
COPY . . 
RUN apt-get update && \
apt-get install -y unzip wget && \
wget https://github.com/protocolbuffers/protobuf/releases/download/v24.3/protoc-24.3-linux-x86_64.zip && \
unzip protoc-24.3-linux-x86_64.zip -d /usr/local && \
rm -f protoc-24.3-linux-x86_64.zip
RUN go mod tidy 
ENV PATH="/usr/local/bin:/go/bin:${PATH}"
#enviroment variables
#...
CMD ["make", "run"]