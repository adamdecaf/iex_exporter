FROM alpine:latest AS ca
RUN apk add -U ca-certificates

FROM scratch
COPY --from=ca /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY bin/iex-exporter-linux /bin/iex-exporter
EXPOSE 9099
ENTRYPOINT ["/bin/iex-exporter"]
CMD [""]
