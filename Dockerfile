# Copied from: https://github.com/xtuc/redis-proxy/blob/master/Dockerfile
FROM golang:1.17 as builder

ENV APP_HOME $GOPATH/src/github.com/manojankitha/redis-proxy

COPY ./ $APP_HOME

WORKDIR $APP_HOME

RUN make install-deps build

RUN cp -f $APP_HOME/redis-proxy /usr/bin

FROM scratch

COPY --from=builder /usr/bin/redis-proxy /usr/bin/redis-proxy

CMD ["/usr/bin/redis-proxy"]