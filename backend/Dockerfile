FROM golang:1.27.1-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/calculator .

FROM scratch
COPY --from=build /out/calculator /calculator
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/calculator"]
