FROM m.daocloud.io/docker.io/library/ubuntu:22.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    python3 \
    python3-pip && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /home/user

COPY hello.py /home/user/hello.py

ENV APP_ENV=production

CMD ["python3", "hello.py"]
