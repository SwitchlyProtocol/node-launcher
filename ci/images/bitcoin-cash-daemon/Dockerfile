FROM debian:11-slim AS installer

WORKDIR /work

RUN apt-get update \
  && apt-get -y install --no-install-recommends \
  ca-certificates=20210119 \
  gnupg=2.2.27-2+deb11u2 \
  wget=1.21-1+deb11u1 \
  jq=1.6-2.1

ENV BITCOIN_SHA256=d69ee632147f886ca540cecdff5b1b85512612b4c005e86b09083a63c35b64fa
COPY version /version

# trunk-ignore(hadolint/DL4006): sh does not support -o pipefail
RUN set -ex && \
  VERSION="$(sed 's/-.*//' /version)" && \
  BITCOIN_URL=https://github.com/bitcoin-cash-node/bitcoin-cash-node/releases/download/v$VERSION/bitcoin-cash-node-$VERSION-x86_64-linux-gnu.tar.gz && \
  wget -qO bitcoin.tar.gz "$BITCOIN_URL" && \
  sha256sum bitcoin.tar.gz && \
  echo "$BITCOIN_SHA256  bitcoin.tar.gz" | sha256sum -c - && \
  mkdir -p /work/bitcoin && \
  tar -xzvf bitcoin.tar.gz -C /work/bitcoin --strip-components=1 --exclude=*-qt


FROM debian:11-slim

COPY --from=installer /work/bitcoin/bin/* /usr/local/bin/
COPY --from=installer /work/bitcoin/lib/* /usr/local/lib/
COPY --from=installer /work/bitcoin/share/* /usr/local/share/

COPY --from=installer /usr/bin/wget /usr/local/bin/
COPY --from=installer /usr/lib/x86_64-linux-gnu/libpsl.so.5  /lib/x86_64-linux-gnu/

COPY --from=installer /usr/bin/jq /usr/local/bin/
COPY --from=installer /usr/lib/x86_64-linux-gnu/libjq.so.1 /usr/lib/x86_64-linux-gnu/
COPY --from=installer /usr/lib/x86_64-linux-gnu/libonig.so.5 /usr/lib/x86_64-linux-gnu

COPY ./scripts /scripts

RUN useradd --create-home bitcoin
USER bitcoin

EXPOSE 8332 8333 18332 18333 18443 18444
VOLUME ["/home/bitcoin/.bitcoin"]

ENTRYPOINT ["/scripts/entrypoint.sh"]
CMD ["bitcoind"]
