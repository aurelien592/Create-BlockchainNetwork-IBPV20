FROM golang:1.8

WORKDIR /go/src/CryptoMotionCoin
WORKDIR /go/src/aurelien592
COPY . .

RUN apt-get update
RUN apt-get install -y vim
RUN go-wrapper download
RUN go-wrapper install
CMD ["go-wrapper", "run"]