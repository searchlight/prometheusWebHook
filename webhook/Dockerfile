FROM golang:alpine
WORKDIR /src
COPY . .
RUN go build -o ./webhook .

CMD ["./webhook"]