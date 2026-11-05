FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY internal ./internal
COPY backend ./backend
RUN CGO_ENABLED=0 go build -o /metrics-traffic ./backend

FROM scratch
COPY --from=build /metrics-traffic /metrics-traffic
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/metrics-traffic"]
