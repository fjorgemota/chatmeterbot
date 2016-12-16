FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
COPY chatmeterbot /bin/
CMD ["/bin/chatmeterbot"]
