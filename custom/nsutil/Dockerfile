FROM arm64v8/ubuntu:latest

# Install gcc for compiling C code
RUN apt-get update && apt-get install -y gcc

# Set the working directory
WORKDIR /app

# Copy your C code into the container
COPY nsutil.c /app

# Compile the C code
RUN gcc -shared -fPIC nsutil.c -o nsutil.so
