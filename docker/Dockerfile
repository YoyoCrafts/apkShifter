FROM centos:7 AS builder

#FROM centos:7

RUN yum -y update


RUN yum  -y install java zlib.i686 libstdc++.so.6


WORKDIR /app




ENV TZ=Asia/Shanghai

COPY config /app/config
COPY main /app/main


# 执行 main 程序
CMD ["./main"]
