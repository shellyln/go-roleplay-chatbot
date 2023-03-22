FROM busybox

COPY ./rpchatd /rpchatd

CMD ["/rpchatd"]
