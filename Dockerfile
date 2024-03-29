FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN export GOPROXY=https://proxy.golang.com.cn,direct
RUN go mod download

COPY *.go ./
COPY certificate/ certificate/

#RUN go build -gcflags="-l" -o /usr/bin/go-server-sample
RUN go build -o /usr/bin/go-server-sample

CMD ["/usr/bin/go-server-sample"]