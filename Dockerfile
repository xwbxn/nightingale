FROM python:2
#FROM ubuntu:21.04

WORKDIR /app
ADD n9e /app
ADD etc /app/etc/
ADD https://monset.oss-cn-hangzhou.aliyuncs.com/wait /wait
RUN mkdir -p /app/pub && chmod +x /wait
ADD pub /app/pub/
RUN chmod +x n9e

EXPOSE 19000
EXPOSE 18000

CMD ["/app/n9e", "-h"]
