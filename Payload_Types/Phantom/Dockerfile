FROM python:3.11-slim

WORKDIR /Mythic/

COPY requirements.txt .
RUN pip install -r requirements.txt

# Install Go for building the agent
RUN apt-get update && apt-get install -y \
    wget \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install Go
RUN wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz && \
    rm go1.21.5.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Install garble for obfuscation
RUN go install mvdan.cc/garble@latest

COPY . .

CMD ["python", "mythic_service.py"]