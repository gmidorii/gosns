FROM golang:1.8.3-onbuild

WORKDIR /go/src/app
RUN touch "subscribed.json"

CMD ["go-wrapper", "run", "-p", "8888"]