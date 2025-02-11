
# guidelines where found here https://devtodevops.com/podman-build-from-dockerfile/
FROM ubuntu

LABEL org.opencontainers.image.source="https://github.com/mathiseng/webservice-project"


COPY ./artifact.bin ./artifact.bin

ENV HOST 0.0.0.0

CMD ["./artifact.bin"]