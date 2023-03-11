FROM alpine:latest

COPY ./dist/main /main

EXPOSE 9000

CMD [ "/main" ]
