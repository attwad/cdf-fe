FROM golang:1.8

WORKDIR /go/src/app
COPY server .
COPY . .

RUN go-wrapper download
RUN go-wrapper install

# Provide a sensible default run command.
CMD ["go-wrapper", "run", "--project_id=college-de-france", "--listen_addr=0.0.0.0:80"]
