FROM golang:latest AS build

WORKDIR /random_work_dir

# first download dependencies
COPY go.mod /random_work_dir
COPY go.sum /random_work_dir
RUN go mod download

# then copy source code
COPY / /random_work_dir


RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /ocrmymail ./cmd/ocrmymail



FROM jbarlow83/ocrmypdf:latest


# Copy build results (https://github.com/jbarlow83/OCRmyPDF)
WORKDIR /

COPY --from=build /ocrmymail /bin/ocrmymail

WORKDIR /ocrmymail
RUN mkdir ./tmp

RUN chmod +x /bin/ocrmymail

EXPOSE 1234

ENTRYPOINT []
CMD ["/bin/ocrmymail"]