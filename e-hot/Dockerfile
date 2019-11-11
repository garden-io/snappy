FROM golang:alpine
RUN apk add entr ca-certificates
EXPOSE 8080
RUN mkdir app
WORKDIR /app
COPY src/ /app/
CMD ls main.go | entr -r ./buildrun.sh