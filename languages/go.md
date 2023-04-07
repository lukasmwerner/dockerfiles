# Go Programming Language

```dockerfile
FROM golang:alpine AS build

RUN apk add git

WORKDIR /src/

COPY go.* /src/

RUN go mod download -x

COPY . /src

#Compiler Settings
ENV CGO_ENABLED=0

# for full parings check out https://go.dev/doc/install/source#environment
ENV GOOS=linux
# this will be the cpu arch
# Can be amd64 arm64 386 ppc64
ENV GOARCH=amd64

RUN go build -o /out/app .

FROM scratch AS run

COPY --from=build /out/app /

# if needed
EXPOSE 80 

ENTRYPOINT [ "/app" ]
```
