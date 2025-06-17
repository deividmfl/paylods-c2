FROM mcr.microsoft.com/dotnet/sdk:6.0

# Install Python and required dependencies
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    mono-complete \
    msbuild \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /Mythic/

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY . .

CMD ["python3", "mythic_service.py"]