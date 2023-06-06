FROM scratch
COPY nsq_forward /
ENTRYPOINT ["/nsq_forward"]
