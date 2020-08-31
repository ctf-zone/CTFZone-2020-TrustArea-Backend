FROM golang:1.14.0
RUN mkdir /src
RUN mkdir /app
RUN mkdir /app/config
COPY ./src /src

WORKDIR /src
RUN go mod download
RUN go build -o /app/main .
RUN cp /src/config/prod.json /app/config/prod.json
RUN rm -rf /src

EXPOSE 8080
WORKDIR /app
CMD ["/app/main"]
