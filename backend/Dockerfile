FROM golang:1.26.1-bookworm

RUN apt-get update && apt-get install -y \
    unzip \
    curl \
    libaio1 \
    gcc \
    libc6-dev

WORKDIR /opt/oracle

RUN curl -o ic.zip https://download.oracle.com/otn_software/linux/instantclient/2360000/instantclient-basic-linux.x64-23.6.0.24.10.zip \
    && unzip ic.zip \
    && rm ic.zip

ENV LD_LIBRARY_PATH="/opt/oracle/instantclient_23_6"
ENV CGO_ENABLED=1

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o dbmt .

WORKDIR /app

ENV LD_LIBRARY_PATH="/opt/oracle/instantclient_23_6"

EXPOSE 5461

CMD ["./dbmt"]
