version: '2'
services:
  bootnode:
    build: .
    restart: on-failure
    ports:
      - 9669:9669
    volumes:
      - .:/go/src/CryptoMotionCoin
      - .:/go/src/aurelien592
    networks:
      regnet:
        ipv4_address: 10.5.0.99
    entrypoint:
      - go
      - run
      - main.go
  node:
    build: .
    restart: on-failure
    ports:
      - 9668
    volumes:
      - .:/go/src/CryptoMotionCoin
      - .:/go/src/aurelien592
    networks:
      - regnet
    entrypoint:
      - go
      - run
      - main.go
      - -entry=10.5.0.99
networks:
  regnet:
    driver: bridge
    ipam:
      config:
        -
          subnet: 10.5.0.0/16
          gateway: 10.5.0.1