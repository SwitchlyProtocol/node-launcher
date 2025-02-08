FROM debian:11-slim AS installer

WORKDIR /work

RUN apt-get update \
  && apt-get -y install --no-install-recommends \
  ca-certificates=20210119 \
  gnupg=2.2.27-2+deb11u2 \
  wget=1.21-1+deb11u1 \
  jq=1.6-2.1

COPY version /version

RUN VERSION="$(sed 's/-.*//' /version)" && \
  wget  --progress=dot:giga "https://github.com/dogecoin/dogecoin/releases/download/v${VERSION}/dogecoin-${VERSION}-x86_64-linux-gnu.tar.gz" && \
  mkdir dogecoin && \
  tar xzvf "dogecoin-${VERSION}-x86_64-linux-gnu.tar.gz" --strip-components=1 -C dogecoin

FROM debian:11-slim

COPY --from=installer /work/dogecoin/bin/* /usr/local/bin/
COPY --from=installer /work/dogecoin/lib/* /usr/local/lib/
COPY --from=installer /work/dogecoin/share/* /usr/local/share/

COPY --from=installer /usr/bin/wget /usr/local/bin/
COPY --from=installer /usr/lib/x86_64-linux-gnu/libpsl.so.5  /lib/x86_64-linux-gnu/

COPY --from=installer /usr/bin/jq /usr/local/bin/
COPY --from=installer /usr/lib/x86_64-linux-gnu/libjq.so.1 /usr/lib/x86_64-linux-gnu/
COPY --from=installer /usr/lib/x86_64-linux-gnu/libonig.so.5 /usr/lib/x86_64-linux-gnu

COPY ./scripts /scripts

RUN useradd --create-home dogecoin
USER dogecoin

EXPOSE 22555 22556 44555 44556 18332 18444
VOLUME ["/home/dogecoin/.dogecoin"]

ENTRYPOINT ["/scripts/entrypoint.sh"]
CMD ["dogecoind"]
