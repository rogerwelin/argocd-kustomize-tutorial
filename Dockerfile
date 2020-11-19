FROM golang:1.14 as build
ARG version
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=$version"

FROM scratch
WORKDIR /app
EXPOSE 8080
COPY --from=build /build/hello-argocd-app .
CMD ["./hello-argocd-app"]
