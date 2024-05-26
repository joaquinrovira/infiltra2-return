ARG BIN_NAME="infiltra2"
ARG RUNTIME_USER="user"
ARG WORKDIR="/app"

FROM golang:alpine as build

# https://github.com/moby/moby/issues/37622#issuecomment-412101935
ARG BIN_NAME
ARG WORKDIR

WORKDIR ${WORKDIR}
COPY . .
RUN go run ./build -outdir . -outname ${BIN_NAME}

FROM alpine

ARG BIN_NAME
ARG RUNTIME_USER
ARG WORKDIR

RUN adduser -h ${WORKDIR} -s /bin/sh -D ${RUNTIME_USER}

USER ${RUNTIME_USER}
WORKDIR ${WORKDIR}

# Copy executable binary
COPY --from=build --chown=${RUNTIME_USER}:${RUNTIME_USER} ${WORKDIR}/${BIN_NAME} .

# Include static directory
COPY --chown=${RUNTIME_USER}:${RUNTIME_USER} static ./static

ENV WORKDIR ${WORKDIR}
ENV BIN_NAME ${BIN_NAME}
ENTRYPOINT ${WORKDIR}/${BIN_NAME}