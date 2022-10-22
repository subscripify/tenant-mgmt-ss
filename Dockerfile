## this is a custom container that has authorization to access private repos on subscripify's azure devops
FROM subscripifycontreg.azurecr.io/builders/subscripifygolang:1.19.2

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
RUN go build -v -o /usr/local/bin/tenant-mgmt-ss .

EXPOSE 8080


ENTRYPOINT [ "tenant-mgmt-ss" ]
CMD ["serve"]
