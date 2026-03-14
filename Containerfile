FROM docker.io/library/ubuntu:latest

RUN apt update && apt install curl -y

USER ubuntu

RUN curl -fsSL https://claude.ai/install.sh | bash

ENTRYPOINT ["claude"]