FROM golang:1.21 as build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main

FROM golang:1.21 as production-stage
WORKDIR /
COPY --from=build-stage /main /main
EXPOSE 8080
CMD [ "/main" ]