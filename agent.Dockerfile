FROM alpine
EXPOSE 7100

COPY make/release/linux/amd64/agent /make/release/agent

ENTRYPOINT ["/make/release/agent"]
CMD ["serve"]
