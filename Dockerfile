FROM alpine:3.18

WORKDIR /opt/campustrader/

ADD campustrader.tar.gz .

RUN chmod +x ./CampusTrader

# 5. 声明监听端口
EXPOSE 8080

CMD ["./CampusTrader"]
