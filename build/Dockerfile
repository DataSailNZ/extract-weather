FROM alpine

COPY ./extract-weather /usr/local/bin/extract-weather

RUN apk add --no-cache tzdata
ENV TZ="Pacific/Auckland"
RUN date

RUN echo "0 6 * * * /usr/local/bin/extract-weather" >> /etc/crontabs/root

CMD ["crond", "-f", "-l", "4", "-L", "/dev/stdout"]
