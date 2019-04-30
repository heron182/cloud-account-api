# build stage
FROM golang:1.12.4-stretch AS build-env
ADD . $GOPATH/src/github.com/heron182/cloud-account-api
WORKDIR $GOPATH/src/github.com/heron182/cloud-account-api
# install go dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go build -o api

# final stage
FROM golang:1.12.4-stretch
WORKDIR /app
COPY --from=build-env $GOPATH/src/github.com/heron182/cloud-account-api /app/
ENTRYPOINT ./api