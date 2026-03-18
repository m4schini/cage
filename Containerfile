FROM docker.io/library/ubuntu:latest

RUN apt-get update
RUN apt-get reinstall -y ca-certificates
RUN update-ca-certificates
RUN apt-get install -y --no-install-recommends zsh curl

USER ubuntu

WORKDIR /home/ubuntu

ADD --chown=ubuntu:ubuntu https://raw.githubusercontent.com/grml/grml-etc-core/master/etc/zsh/zshrc .zshrc
RUN echo 'export PATH="$HOME/.local/bin:$PATH"' >> .zshrc

RUN curl -fsSL https://claude.ai/install.sh | bash

WORKDIR /home/ubuntu/workspace

CMD ["/usr/bin/zsh"]