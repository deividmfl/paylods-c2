FROM mcr.microsoft.com/dotnet/sdk:6.0

# Install advanced evasion tools and dependencies
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    mono-complete \
    msbuild \
    upx-ucl \
    openssl \
    wget \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Install Python cryptography tools
RUN pip3 install pycryptodome

WORKDIR /Mythic/

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY . .

# Apply advanced obfuscation and evasion
RUN python3 advanced_obfuscator.py
RUN if [ -f "agent_code/build_evasive.sh" ]; then \
        cd agent_code && chmod +x build_evasive.sh && ./build_evasive.sh; \
    fi

CMD ["python3", "mythic_service.py"]