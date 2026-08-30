FROM golang:1.26-alpine as stc

RUN apk add --no-cache make cmake gcc musl-dev linux-headers git bash build-base libc-dev

ADD . /stc
WORKDIR /stc
ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__" 
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"
RUN make geth
RUN mv /stc/build/bin/geth /usr/local/bin/geth

EXPOSE 8545 8547 30303 30303/udp
ENTRYPOINT [ "/usr/local/bin/geth" ]

FROM ethereum/solc:0.6.4-alpine as stc-genesis

RUN apk add --no-cache ca-certificates npm nodejs bash alpine-sdk expect

RUN git clone https://github.com/secblockchain/stc-genesis-contract.git /root/genesis \
    && cd /root/genesis && npm install
#RUN curl -L https://foundry.paradigm.xyz | bash

COPY --from=stc /usr/local/bin/geth /usr/local/bin/geth

ENTRYPOINT [ "/bin/bash" ]
