FROM alpine:latest
RUN apk add -U ca-certificates tzdata
COPY bin/iex-exporter-linux /bin/iex-exporter
EXPOSE 9099
ENTRYPOINT ["/bin/iex-exporter"]
CMD [""]
