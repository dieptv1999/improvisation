FROM golang

ENV GO111MODULE on

RUN go version

COPY . /src

RUN ls

WORKDIR /src

RUN ls

RUN go mod vendor

RUN go build

EXPOSE 8443

ENTRYPOINT [ "./go-app", "daemon" ]