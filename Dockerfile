FROM python:2
#FROM ubuntu:21.04

WORKDIR /app
ADD n9e /app
ADD etc /app/etc/
# RUN curl -sSL https://monset.oss-cn-hangzhou.aliyuncs.com/pub.tar.gz | tar xz
RUN chmod +x n9e

EXPOSE 19000
EXPOSE 18000

CMD ["/app/n9e", "-h"]
