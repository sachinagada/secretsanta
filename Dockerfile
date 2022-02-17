# syntax=docker/dockerfile:1

FROM golang:1.17.7-alpine as build
RUN apk add --update --no-cache ca-certificates shadow
# create a group and a user
# -d specifies the home directory and -m flag creates the home directory if it
# doesn't exist
RUN groupadd --gid 1000 build && useradd --uid 1000 --gid 1000 -d /home/app -m build
USER build
# set the working directory to be used for copying and and building the app
WORKDIR /home/app
ENV GOMODULE=on

# copy everything and change owner of files to the build user
COPY --chown=build:build . .
RUN go mod download
RUN go build -mod=readonly -o secret-santa ./cmd

FROM alpine
RUN apk --no-cache add ca-certificates
RUN adduser -D santa
USER santa
WORKDIR /app/
COPY --from=build /home/app/secret-santa ./
EXPOSE 3000 8080
ENTRYPOINT [ "./secret-santa" ]
