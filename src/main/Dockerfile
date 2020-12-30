FROM golang:alpine


ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    MONGODB_HOST=mongodb://root:dummypw@mongodb.default.svc.cluster.local:27017/?authSource=admin \
    MONGODB_DATABASE=k8sportal\
    MONGODB_COLLECTION=portal-services


WORKDIR /build


COPY go.mod .
COPY go.sum .
RUN go mod download


COPY . .


RUN go build -o main .


WORKDIR /dist


RUN cp /build/main .


EXPOSE 80


CMD ["/dist/main"]