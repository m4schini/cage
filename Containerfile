# =============================================================================
# Dev Container — Debian + Nix Shell
# =============================================================================
#
# Build example:
#   podman run -it \
#       --volume=./shell.nix:/home/developer/shell.nix:ro \
#       --volume=/home/aurora/.zshrc:/home/developer/.zshrc \
#       --volume=nixstore:/nix:Z \
#       --rm \
#       my-devcontainer
#
# Run example:
#   podman run --rm -it my-devcontainer
# =============================================================================

ARG DEBIAN_VERSION=bookworm-slim
FROM debian:${DEBIAN_VERSION}

# -----------------------------------------------------------------------------
# Build arguments
# -----------------------------------------------------------------------------

# UID/GID for the non-root developer user
ARG DEV_UID=1000
ARG DEV_GID=1000
ARG DEV_USER=developer

# -----------------------------------------------------------------------------
# 1. System bootstrap
# -----------------------------------------------------------------------------
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    set -eux; \
    apt-get update -qq; \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        xz-utils \
        git \
        bash \
        procps \
        sudo;

# -----------------------------------------------------------------------------
# Create the non-root developer user
# -----------------------------------------------------------------------------
RUN set -eux; \
    groupadd --gid "${DEV_GID}" "${DEV_USER}"; \
    useradd \
        --uid "${DEV_UID}" \
        --gid "${DEV_GID}" \
        --shell /bin/bash \
        --create-home \
        "${DEV_USER}"; \
    echo "${DEV_USER} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${DEV_USER}; \
    chmod 0440 /etc/sudoers.d/${DEV_USER}

# -----------------------------------------------------------------------------
# Install Nix (single-user mode, no daemon required)
# -----------------------------------------------------------------------------
USER ${DEV_USER}
WORKDIR /home/${DEV_USER}
RUN curl -fsSL https://nixos.org/nix/install | sh -s -- --no-daemon
RUN . "/home/${DEV_USER}/.nix-profile/etc/profile.d/nix.sh"
RUN /home/${DEV_USER}/.nix-profile/bin/nix --version # Smoke-test
ENV NIXPKGS_ALLOW_UNFREE=1
RUN touch /home/${DEV_USER}/shell.nix

# -----------------------------------------------------------------------------
# Label
# -----------------------------------------------------------------------------
LABEL org.opencontainers.image.title="cage"
LABEL org.opencontainers.image.description="Isolate agentic ai"

LABEL org.opencontainers.image.vendor="cage"
LABEL security.non-root="true"
LABEL security.no-new-privileges="true"

WORKDIR /home/${DEV_USER}/workspace
USER ${DEV_USER}

# Default to an interactive bash shell; override as needed.
ENV PATH="/home/${DEV_USER}/.nix-profile/bin:${PATH}"
ENTRYPOINT ["nix-shell"]
CMD ["../shell.nix"]