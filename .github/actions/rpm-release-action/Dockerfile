FROM almalinux:9-minimal

# Install depedencies
RUN dnf -y --refresh update && \
  dnf -y install rpm-build rpmdevtools git

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
