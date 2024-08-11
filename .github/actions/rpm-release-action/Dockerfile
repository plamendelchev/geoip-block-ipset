FROM almalinux:9

# Install depedencies
RUN dnf --refresh -y install rpm-build git tar 'dnf-command(builddep)'

COPY --chmod=0744 entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
