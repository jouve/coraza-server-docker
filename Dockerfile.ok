FROM golang:1.17.6-alpine3.15

ARG LIBINJECTION_VERSION=3.10.0
ARG LIBINJECTION_SHA512SUM=15d144439d99dafd04703c7d6db329120bc24abb5d954d9006d1c5aba6acd79cde1e1af3faad6f263782c042d458c567ce12d4c4c6e44baf1d0a44c94d6564ee
RUN set -e; \
    apk add --no-cache gcc make musl-dev pcre-dev pkgconfig python2; \
    wget "https://github.com/client9/libinjection/archive/refs/tags/v${LIBINJECTION_VERSION}.tar.gz"; \
    echo "${LIBINJECTION_SHA512SUM}  v${LIBINJECTION_VERSION}.tar.gz" | sha512sum -c; \
    tar xf "v${LIBINJECTION_VERSION}.tar.gz"; \
    cd "libinjection-${LIBINJECTION_VERSION}"; \
    sed -i -e '/${INSTALL} -c ${STATICLIB} ${LIB_DIR}/a\\t${INSTALL} -c ${SHAREDLIB} ${LIB_DIR}' src/Makefile; \
    ./configure-gcc-hardened.sh; \
    make -C src install PREFIX="$pkgdir/usr";

COPY go.mod go.sum main.go src/

RUN cd src; go build 

FROM alpine:3.15.0

RUN apk add --no-cache pcre

COPY --from=0 /usr/lib/libinjection.so /usr/lib

COPY --from=0 /go/src/coraza-server /usr/bin

COPY config.yml /etc/coraza-server/

ARG CORAZA_WAF_VERSION=2.0.0-rc.1
RUN wget "https://raw.githubusercontent.com/jptosso/coraza-waf/v${CORAZA_WAF_VERSION}/coraza.conf-recommended" -O /etc/coraza-server/coraza.conf

ARG CRS_VERSION=3.3.2
RUN set -e; \
    wget "https://github.com/coreruleset/coreruleset/archive/refs/tags/v${CRS_VERSION}.tar.gz" -O- | tar -xz; \
    mv "coreruleset-${CRS_VERSION}/crs-setup.conf.example" /etc/coraza-server/crs-setup.conf; \
    mv "coreruleset-${CRS_VERSION}/rules" /etc/coraza-server; \
    rm -rf "coreruleset-${CRS_VERSION}"

EXPOSE 12345

CMD ["coraza-server"]
