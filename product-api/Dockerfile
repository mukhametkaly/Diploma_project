#FROM alpine:3.9
#ADD fortekassa-shopping-cart-api fortekassa-shopping-cart-api
#ADD src/config src/config
#ENTRYPOINT ["/fortekassa-shopping-cart-api"]

FROM alpine:3
ADD src/config src/config
ADD target/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]
CMD ["/go/bin/app"]