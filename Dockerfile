FROM golang:1.20.0-alpine3.17 AS BUILDER

ARG GITHUB_TOKEN
ARG GITHUB_USER
ARG GOPRIV

RUN apk add gcc git
ENV GOPRIVATE=${GOPRIV}
RUN git config --global url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -tags nethttpomithttp2 -ldflags="-s -w" -o app .
RUN cd cli && go build -tags nethttpomithttp2 -ldflags="-s -w" -o profile-cli .

FROM alpine:3.17

RUN apk add ca-certificates musl-dev

RUN mkdir firebase-credential
COPY --from=BUILDER ["/build/app","/"]
COPY --from=BUILDER ["/build/cli/profile-cli","/"]
# COPY --from=BUILDER ["/build/ca-certificate.cer", "/"]
COPY --from=BUILDER ["/build/cert", "/"]
COPY --from=BUILDER ["/build/firebase.json","/firebase-credential"]

ARG BUILDDATE

ENV APP_SERVICE_VERSION=$BUILDDATE

EXPOSE 4000

CMD ["/app"]