FROM python:2
#FROM ubuntu:21.04

WORKDIR /app
ADD n9e /app
ADD etc /app/etc/
ADD https://monset.oss-cn-hangzhou.aliyuncs.com/wait /wait
RUN chmod +x /wait && chmod +x n9e

EXPOSE 17000

CMD ["/app/n9e", "-h"]
