## this is a custom container that has authorization to access private repos on subscripify's azure devops
FROM subscripifycontreg.azurecr.io/builders/subscripifygolang:1.19.2 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
# only copy the needed folders
COPY ./cmd ./cmd
COPY ./configs ./configs
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./main.go ./main.go
## CGO_ENABLED is set to 0 so that the binary can run on Alpine
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/tenant-mgmt-ss .






FROM alpine:3.14
# try this if CGO_ENABLED suddenly stops working
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/tenant-mgmt-ss .


EXPOSE 8080
ENTRYPOINT [ "tenant-mgmt-ss" ]
CMD ["serve"]
