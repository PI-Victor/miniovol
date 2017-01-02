FROM alpine

RUN mkdir -p /run/docker/plugins

COPY _out/bin/miniovol miniovol

CMD ["miniovol"]
