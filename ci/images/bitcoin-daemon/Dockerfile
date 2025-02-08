FROM debian:11-slim AS installer

# install dependencies
RUN apt-get update && \
  apt-get -y install --no-install-recommends \
  ca-certificates=20210119 \
  git=1:2.30.2-1+deb11u2 \
  gnupg=2.2.27-2+deb11u2 \
  wget=1.21-1+deb11u1 \
  jq=1.6-2.1


# import builder keys
WORKDIR /keys
RUN git init && \
  git remote add origin https://github.com/bitcoin-core/guix.sigs.git && \
  git fetch --depth 1 origin 32acdd8af28fae2980625c304bcd70fe7636de28 && \
  git checkout FETCH_HEAD && \
  for file in builder-keys/*; do gpg --import "$file"; done

# pull and verify binaries
COPY version /version
WORKDIR /work
RUN VERSION="$(sed 's/-.*//' /version)" && \
  wget --progress=dot:giga "https://bitcoincore.org/bin/bitcoin-core-${VERSION}/bitcoin-${VERSION}-x86_64-linux-gnu.tar.gz" && \
  wget --progress=dot:giga "https://bitcoincore.org/bin/bitcoin-core-${VERSION}/SHA256SUMS" && \
  wget --progress=dot:giga "https://bitcoincore.org/bin/bitcoin-core-${VERSION}/SHA256SUMS.asc"
RUN VERSION="$(sed 's/-.*//' /version)" && \
  gpg --verify SHA256SUMS.asc SHA256SUMS && \
  grep "bitcoin-${VERSION}-x86_64-linux-gnu.tar.gz" SHA256SUMS | sha256sum -c && \
  mkdir bitcoin && \
  tar xzvf "bitcoin-${VERSION}-x86_64-linux-gnu.tar.gz" --strip-components=1 -C bitcoin

FROM debian:11-slim

COPY --from=installer /work/bitcoin/bin/* /usr/local/bin/
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
