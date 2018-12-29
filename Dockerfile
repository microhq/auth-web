FROM alpine:3.2
ADD templates /templates
ADD auth-web /auth-web
WORKDIR /
ENTRYPOINT [ "/auth-web" ]
