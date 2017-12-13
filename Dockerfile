FROM golang:latest
COPY . /usr/src/app
WORKDIR /usr/src/app
RUN go get -u github.com/gorilla/mux
RUN go get -u gopkg.in/mgo.v2
RUN go get -u github.com/opentracing/opentracing-go
RUN go get -u sourcegraph.com/sourcegraph/appdash
RUN go get -u github.com/opentracing/basictracer-go
RUN go get -u sourcegraph.com/sourcegraph/appdash-data
CMD go run payment.go util.go
