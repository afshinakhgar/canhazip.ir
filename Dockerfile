FROM scratch
COPY --chown=0:0 api_linux /api
COPY --chown=0:0 certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chown=0:0 certs/zoneinfo /usr/share/zoneinfo
COPY php-legacy/public/GeoLite2-City.mmdb /app/data/
COPY php-legacy/public/GeoLite2-ASN.mmdb /app/data/
COPY data/firehol_level1.netset /app/data/
COPY data/firehol_level2.netset /app/data/
EXPOSE 8080
CMD ["/api"]
