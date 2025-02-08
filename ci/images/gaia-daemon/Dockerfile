########################################################################################
# Build
########################################################################################

# Source: https://github.com/cosmos/gaia/blob/main/Dockerfile

FROM golang:1.22-alpine AS build

WORKDIR /app
COPY version /version

ENV PACKAGES="curl make git libc-dev bash file gcc linux-headers eudev-dev python3"
RUN apk add --no-cache $PACKAGES

ARG SRC_COMMIT=cc8fa2bbc24d55f0cacf51f8c06b3fd1c2762d17
RUN VERSION="$(sed 's/-.*//' /version)" && \
  git init && \
  git remote add origin https://github.com/cosmos/gaia && \
  git fetch --depth 1 origin "v${VERSION}" && \
  git checkout $SRC_COMMIT

# See https://github.com/CosmWasm/wasmvm/releases
ARG WASMVM_VERSION=v2.1.4
ADD https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 090b97641157fae1ae45e7ed368a1a8c091f3fef67958d3bc7c2fa7e7c54b6b4
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep a4a3d09b36fabb65b119d5ba23442c23694401fcbee4451fe6b7e22e325a4bac
RUN cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

# build
RUN go mod download
RUN LEDGER_ENABLED=false LINK_STATICALLY=true BUILD_TAGS=muslc make build
RUN echo "Ensuring binary is statically linked ..."  \
  && file /app/build/gaiad | grep "statically linked"

########################################################################################
# Deploy
########################################################################################

FROM alpine:3.20
RUN apk add --no-cache curl jq
COPY --from=build /app/build/gaiad /gaiad
COPY entrypoint.sh /entrypoint.sh
COPY app.toml /etc/gaia/app.toml
CMD ["/entrypoint.sh"]
