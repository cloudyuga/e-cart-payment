FROM golang:latest
COPY . /usr/src/app
WORKDIR /usr/src/app
RUN go get -u github.com/gorilla/mux
RUN go get -u gopkg.in/mgo.v2
CMD go run payment.go
