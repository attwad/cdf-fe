FROM golang:1.9

WORKDIR /go/src/app
COPY server .
COPY . .

RUN go-wrapper download
RUN go-wrapper install

# Provide a sensible default run command.
CMD ["go-wrapper", "run", "--project_id=college-de-france", "--listen_addr=127.0.0.1:80", "--elastic_address=http://127.0.0.1:9200"]
