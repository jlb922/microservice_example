#FROM golang:1.9.2 as builder
#ARG SOURCE_LOCATION=/
#WORKDIR ${SOURCE_LOCATION}
#RUN go get -d -v gopkg.in/mgo.v2/bson \
#    && go get -d -v gopkg.in/mgo.v2
#COPY main.go .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
#
#FROM alpine:latest  
#ARG SOURCE_LOCATION=/
#RUN apk --no-cache add curl
#EXPOSE 9090
#WORKDIR /root/
#COPY --from=builder ${SOURCE_LOCATION} .
#CMD ["./app"]  

FROM golang:alpine as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
#RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest  
RUN apk --no-cache add curl

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Export necessary port
EXPOSE 9090

# Command to run when starting the container
CMD ["/dist/main"]