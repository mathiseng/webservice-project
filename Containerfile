
# guidelines where found here https://devtodevops.com/podman-build-from-dockerfile/
FROM ubuntu

COPY ./artifact.bin ./artifact.bin

ENV HOST 0.0.0.0

CMD ["./artifact.bin"]