# syntax=docker/dockerfile:1

ARG ALPINE_IMAGE=docker.io/library/alpine:latest

################################################################################
# Jupyter stage
################################################################################
# hadolint ignore=DL3006
FROM $ALPINE_IMAGE

SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]

# Install jupyter prerequisites.
# hadolint ignore=DL3018,DL3059
RUN apk add --no-cache --upgrade \
    ca-certificates \
    g++ \
    gcc \
    git \
    libffi-dev \
    python3 \
    python3-dev \
    py3-pip \
    R \
    R-dev \
    zeromq \
    zeromq-dev

# hadolint ignore=DL3013,DL3059
RUN pip3 install --no-cache-dir wheel

# Install jupyter lab.
# hadolint ignore=DL3013,DL3059
RUN pip3 install --no-cache-dir \
    jupyter \
    notebook \
    jupyterlab

# Install python kernel for jupyter.
# hadolint ignore=DL3013,DL3059
RUN pip install --no-cache-dir ipykernel

# Install R kernel for jupyter.
ARG CRAN_URL=https://packagemanager.rstudio.com/cran/latest

# hadolint ignore=DL3059
RUN mkdir -p /usr/share/doc/R/html; \
    Rscript -e "install.packages('IRkernel',repos='$CRAN_URL',type='source')"; \
    Rscript -e "IRkernel::installspec(user=FALSE)"

# hadolint ignore=DL3059
RUN Rscript -e "install.packages('supernova',repos='$CRAN_URL',type='source')"

ENTRYPOINT [ "jupyter" ]
