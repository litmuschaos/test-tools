# This Dockerfile contains the hardened alpine image with all the
# litmus experiment dependencies installed.
# It is also made non-root, sudo-enabled with default litmus directory.
FROM alpine:3.16.9

LABEL maintainer="LitmusChaos"

# make a pipe fail on the first failure
SHELL ["/bin/sh", "-o", "pipefail", "-c"]

# ensure we only use apk repositories over HTTPS (altough APK contain an embedded signature)
RUN echo "https://alpine.global.ssl.fastly.net/alpine/v$(cut -d . -f 1,2 < /etc/alpine-release)/main" > /etc/apk/repositories \
	&& echo "https://alpine.global.ssl.fastly.net/alpine/v$(cut -d . -f 1,2 < /etc/alpine-release)/community" >> /etc/apk/repositories

ENV GLIBC_REPO=https://github.com/sgerrand/alpine-pkg-glibc
ENV GLIBC_VERSION=2.30-r0
ARG TARGETARCH
ARG LITMUS_VERSION

RUN rm -rf /var/lib/apt/lists/*

# Install generally useful things
RUN apk --update add \
        sudo \
	      htop\
        bash\
        make\
        git \
        curl\
        iproute2\
        stress-ng\
        openssh-client\
#        libc6-compat \
        sshpass \
        expat-doc

RUN set -ex && \
     apk --update add libstdc++ curl ca-certificates && \
     for pkg in glibc-${GLIBC_VERSION} glibc-bin-${GLIBC_VERSION}; \
         do curl -sSL ${GLIBC_REPO}/releases/download/${GLIBC_VERSION}/${pkg}.apk -o /tmp/${pkg}.apk; done && \
     apk add --allow-untrusted /tmp/*.apk && \
     rm -v /tmp/*.apk && \
     /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib

# Change default shell from ash to bash
RUN sed -i -e "s/bin\/ash/bin\/bash/" /etc/passwd

#Installing Kubectl
ENV KUBE_LATEST_VERSION="v1.21.2"
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/${TARGETARCH}/kubectl -o     /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

#Installing crictl binaries
RUN curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.28.0/crictl-v1.28.0-linux-${TARGETARCH}.tar.gz --output crictl-v1.28.0-linux-${TARGETARCH}.tar.gz && \
    tar zxvf crictl-v1.28.0-linux-${TARGETARCH}.tar.gz -C /usr/local/bin

#Installing promql cli binaries
RUN curl -L https://github.com/chaosnative/promql-cli/releases/download/3.0.0/promql_linux_${TARGETARCH} --output /usr/local/bin/promql && chmod +x /usr/local/bin/promql

#Installing nsutil cli binaries
RUN curl -L https://github.com/litmuschaos/test-tools/releases/download/${LITMUS_VERSION}/nsutil-linux-${TARGETARCH} --output /usr/local/bin/nsutil && chmod +x /usr/local/bin/nsutil

#Installing pause cli binaries
RUN curl -L https://github.com/litmuschaos/test-tools/releases/download/${LITMUS_VERSION}/pause-linux-${TARGETARCH} --output /usr/local/bin/pause && chmod +x /usr/local/bin/pause

#Installing dns_interceptor cli binaries
RUN curl -L https://github.com/litmuschaos/test-tools/releases/download/${LITMUS_VERSION}/dns_interceptor --output /usr/local/bin/dns_interceptor && chmod +x /usr/local/bin/dns_interceptor

COPY --from=docker:19.03 /usr/local/bin/docker /usr/local/bin/

# Installing toxiproxy binaries
RUN curl -L https://litmus-http-proxy.s3.amazonaws.com/cli/toxiproxy-cli-linux-${TARGETARCH}.tar.gz --output toxiproxy-cli-linux-${TARGETARCH}.tar.gz && \
    tar zxvf toxiproxy-cli-linux-${TARGETARCH}.tar.gz -C /usr/local/bin/ && \
    chmod +x /usr/local/bin/toxiproxy-cli
RUN curl -L https://litmus-http-proxy.s3.amazonaws.com/server/toxiproxy-server-linux-${TARGETARCH}.tar.gz --output toxiproxy-server-linux-${TARGETARCH}.tar.gz && \
    tar zxvf toxiproxy-server-linux-${TARGETARCH}.tar.gz -C /usr/local/bin/ && \
    chmod +x /usr/local/bin/toxiproxy-server

# The user the app should run as
ENV APP_USER=litmus
# The home directory
ENV APP_DIR="/$APP_USER"
# Where persistent data (volume) should be stored
ENV DATA_DIR "$APP_DIR/data"
# Where configuration should be stored
ENV CONF_DIR "$APP_DIR/conf"

# Update base system
# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates

# Add custom user and setup home directory
RUN adduser -s /bin/true -u 1000 -D -h $APP_DIR $APP_USER \
  && mkdir "$DATA_DIR" "$CONF_DIR" \
  && chown -R "$APP_USER" "$APP_DIR" "$CONF_DIR" \
  && chmod 700 "$APP_DIR" "$DATA_DIR" "$CONF_DIR" \
# change to 0(root) group because openshift will run container with arbitrary uid as a member of root group
  && chgrp -R 0 "$APP_DIR" "$DATA_DIR" "$CONF_DIR" \
  && chmod -R g=u "$APP_DIR" "$DATA_DIR" "$CONF_DIR" \
  && echo "$APP_USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/litmus \
  && chmod 0440 /etc/sudoers.d/litmus

# Remove existing crontabs, if any.
RUN rm -fr /var/spool/cron \
	&& rm -fr /etc/crontabs \
	&& rm -fr /etc/periodic

# Remove all but a handful of admin commands.
RUN find /sbin /usr/sbin \
  ! -type d -a ! -name apk -a ! -name ln \
  -delete

# Remove world-writeable permissions except for /tmp/
RUN find / -xdev -type d -perm +0002 -exec chmod o-w {} + \
	&& find / -xdev -type f -perm +0002 -exec chmod o-w {} + \
	&& chmod 777 /tmp/ \
  && chown $APP_USER:root /tmp/

# Remove unnecessary accounts, excluding current app user and root
RUN sed -i -r "/^($APP_USER|root|nobody)/!d" /etc/group \
  && sed -i -r "/^($APP_USER|root|nobody)/!d" /etc/passwd

# Remove interactive login shell for everybody
RUN sed -i -r 's#^(.*):[^:]*$#\1:/sbin/nologin#' /etc/passwd

# Disable password login for everybody
RUN while IFS=: read -r username _; do passwd -l "$username"; done < /etc/passwd || true

# Remove apk configs. -> Commented out because we need apk to install other stuff
RUN find /bin /etc /lib /sbin /usr \
 -xdev -type f -regex '.*apk.*' \
 ! -name apk \
 -exec rm -fr {} +

# Remove temp shadow,passwd,group
RUN find /bin /etc /lib /sbin /usr -xdev -type f -regex '.*-$' -exec rm -f {} +

# Ensure system dirs are owned by root and not writable by anybody else.
RUN find /bin /etc /lib /sbin /usr -xdev -type d \
  -exec chown root:root {} \; \
  -exec chmod 0755 {} \;

# Remove suid & sgid files
RUN find /bin /etc /lib /sbin /usr -xdev -type f -a \( -perm +4000 -o -perm +2000 \) -delete

# Remove dangerous commands
RUN find /bin /etc /lib /sbin /usr -xdev \( \
  -iname hexdump -o \
  -iname chgrp -o \
  -iname ln -o \
  -iname od -o \
  -iname strings -o \
  -iname su -o \
  # -iname sudo \
  \) -delete

# Remove init scripts since we do not use them.
RUN rm -fr /etc/init.d /lib/rc /etc/conf.d /etc/inittab /etc/runlevels /etc/rc.conf /etc/logrotate.d

# Remove kernel tunables
RUN rm -fr /etc/sysctl* /etc/modprobe.d /etc/modules /etc/mdev.conf /etc/acpi

# Remove root home dir
RUN rm -fr /root

# Remove fstab
RUN rm -f /etc/fstab

# Remove any symlinks that we broke during previous steps
RUN find /bin /etc /lib /sbin /usr -xdev -type l -exec test ! -e {} \; -delete

# default directory is /litmus
WORKDIR $APP_DIR
USER ${APP_USER}