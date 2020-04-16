FROM ubuntu:18.04 as builder

# intall gcc and supporting packages
RUN apt-get update && apt-get install -yq make gcc

WORKDIR /code

# download stress-ng sources
ARG STRESS_NG_VERSION
ENV STRESS_NG_VERSION ${STRESS_NG_VERSION:-0.10.10}
ADD https://github.com/ColinIanKing/stress-ng/archive/V${STRESS_NG_VERSION}.tar.gz .
RUN tar -xf V${STRESS_NG_VERSION}.tar.gz && mv stress-ng-${STRESS_NG_VERSION} stress-ng

#install stress
RUN apt-get install stress

# make static version
WORKDIR /code/stress-ng
RUN STATIC=1 make

# Final image
FROM scratch

COPY --from=builder /code/stress-ng/stress-ng /

ENTRYPOINT ["/stress-ng"]
