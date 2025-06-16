FROM itsafeaturemythic/mythic_go_base:latest

WORKDIR /Mythic/

COPY [".", "."]

RUN apt-get update && \
    apt-get install -y upx-ucl && \
    go install mvdan.cc/garble@latest

CMD ["python3", "mythic_service.py"]