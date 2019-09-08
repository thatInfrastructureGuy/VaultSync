FROM golang:alpine as build
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ADD passwd.minimal /etc/passwd
USER nobody
ADD ./vaultsync /vaultsync
# Compile with GOOS=linux and GOARCH=amd64 
CMD ["./vaultsync"]
