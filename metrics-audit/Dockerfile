FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY internal ./internal
COPY backend ./backend
RUN CGO_ENABLED=0 go build -o /metrics-audit ./backend

FROM scratch
COPY --from=build /metrics-audit /metrics-audit
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/metrics-audit"]
