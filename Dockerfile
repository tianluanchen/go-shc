FROM golang:alpine3.18
WORKDIR /app
ADD . /app
RUN go build -o app main.go
CMD [ "./app" , "-a",":80"]
EXPOSE 80