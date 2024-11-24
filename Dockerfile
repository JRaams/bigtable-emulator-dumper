FROM golang:1.22.4 AS build
    WORKDIR /go/src
    
    COPY ./src/go.mod ./src/go.sum ./
    RUN --mount=type=cache,target=/go/pkg/mod/ go mod download -x

    COPY ./src ./
    RUN --mount=type=cache,target=/go/pkg/mod/ CGO_ENABLED=0 go build \
        -installsuffix 'static' \
        -ldflags '-s -w' \
        -o /go/bin/app
    
FROM scratch AS final
    COPY --from=build /go/bin/app /go/bin/app
    ENTRYPOINT ["/go/bin/app"]