FROM ubuntu:24.04

# Enable 32-bit architecture and update package lists
RUN dpkg --add-architecture i386 && apt-get update

# Install necessary libraries, tools, and SSH server
RUN apt-get install -y \
    libc6:i386 \
    libstdc++6:i386 \
    gcc-multilib \
    gdb \
    openssh-server

# Set up SSH server
RUN mkdir /var/run/sshd \
    && echo 'root:secure_password' | chpasswd \
    && sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# Expose port 22 for SSH
EXPOSE 22

# Set the working directory
WORKDIR /app

# Copy your 32-bit binary into the container
COPY shutdown_procedure /app/
COPY shutdown_procedure.cpp /app/

# Start SSH server on container start
CMD ["/usr/sbin/sshd", "-D"]
