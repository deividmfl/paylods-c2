FROM python:3.11-slim

WORKDIR /Mythic/

COPY requirements.txt requirements.txt

RUN pip3 install -r requirements.txt

# Install Go for payload compilation
RUN apt-get update && apt-get install -y \
    wget \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.21
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz && \
    rm go1.21.6.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV GOCACHE=/tmp/.cache

COPY . .

CMD ["python3", "mythic_service.py"]