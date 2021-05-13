# OCRmyMail

OCRmyMail is an SMTP server relay that adds an OCR text layer to PDF mail attachments and sends them to the original recipient.  
[OCRmyPDF](https://github.com/jbarlow83/OCRmyPDF) is used to do the heavy lifting a.k.a perform the OCR.

OCRmyMail doesn't have any SMTP authentication mechanism or TLS. It is aimed to be used within a local or dockerized environment where authentication doesn't matter. **It is not meant to be put somewhere publicly!**

[![Build Status](https://travis-ci.com/DenBeke/ocrmymail.svg?branch=master)](https://travis-ci.com/DenBeke/ocrmymail)
[![Go Report Card](https://goreportcard.com/badge/github.com/DenBeke/ocrmymail)](https://goreportcard.com/report/github.com/DenBeke/ocrmymail)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/denbeke/ocrmymail?sort=date)](https://hub.docker.com/r/denbeke/ocrmymail)


## Usage (Docker)

### Docker-compose

The easiest way to run OCRmyMail is with docker-compose.
Edit the `.env` file with your settings,  download the [docker-compose.yml](./docker-compose.yml) file and run it with:

```bash
docker-compose up -d
```


### Docker run

If you don't want to use Docker compose, you can always run the command manually:

```bash
docker run -it\
    -e REMOTE_SMTP_HOST=${REMOTE_SMTP_HOST} \
    -e REMOTE_SMTP_PORT=${REMOTE_SMTP_PORT} \
    -e REMOTE_SMTP_DISABLE_TLS=${REMOTE_SMTP_DISABLE_TLS} \
    -e SMTP_FROM_EMAIL=${SMTP_FROM_EMAIL} \
    -e ADMIN_MAIL=${ADMIN_MAIL} \
    -p 25:25 \
    denbeke/ocrmymail
```



## Usage (binary)

[Install OCRmPDF](https://github.com/jbarlow83/OCRmyPDF#installation) first.  
Then download the latest OCRmyMail from the [releases page](https://github.com/DenBeke/ocrmymail/releases).

Configure your settings in the `.env` and run OCRmyMail with:

```bash
./ocrmypdf
```


## Development

[Install OCRmPDF](https://github.com/jbarlow83/OCRmyPDF#installation) first. 
Then run it manually with Go (requires Go 1.15 or newer):

```bash
go run cmd/ocrmymail/*.go
```

To test the email functionality, you can send the `test.txt` SMTP mail with a tool like netcat:

```bash
nc localhost 26 < mail.txt
```

Because netcat and GoPistolet don't run well together, I added a small SMTP gateway between them on port `26`. (Hence the port `26` in the above command):

```bash
docker run -it -e 'ACCEPTED_NETWORKS=192.168.0.0/16 172.16.0.0/12 10.0.0.0/8' -e RELAY_HOST_NAME=test -e 'EXT_RELAY_HOST=10.0.0.49' -e 'EXT_RELAY_PORT=25' -e 'SMTP_LOGIN= ' -e 'SMTP_PASSWORD= ' -p 26:25 alterrebe/postfix-relay
```


## Acknowledgments

- [gopistolet/smtp](https://github.com/gopistolet/smtp)
- [DusanKasan/parsemail](https://github.com/DusanKasan/parsemail)
- [sirupsen/logrus](https://github.com/sirupsen/logrus)
- [gosimple/slug](https://github.com/gosimple/slug)
- [go-mail/mail](https://github.com/go-mail/mail)
- [google/uuid](https://github.com/google/uuid)


## TODO

- [ ] Handle multiple recipients
- [x] Handle collissions when multiple attachments with same filename are to be handled
- [ ] Metrics
- [ ] Error handling with remote tool
- [x] Delete files after sending out mail
- [ ] Tests
- [x] Async mail handling (use `HANDLE_ASYNC=1` env variable)


## Author

[Mathias Beke](https://denbeke.be)

