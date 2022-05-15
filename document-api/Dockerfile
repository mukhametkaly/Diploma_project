FROM golang:1.15-alpine

RUN mkdir $HOME/goproject

WORKDIR $HOME/goproject

COPY . $HOME/goproject


EXPOSE 8080


CMD ["go", "run", "main.go"]