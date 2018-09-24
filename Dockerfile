FROM scratch
ADD  it2s-coap /
EXPOSE 5683
ENTRYPOINT ["./it2s-coap"]
