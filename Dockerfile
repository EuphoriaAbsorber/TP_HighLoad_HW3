FROM golang:latest AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build -o main main.go


FROM ubuntu:latest
RUN apt-get -y update && apt-get install -y tzdata

ENV PGVER 14
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER art WITH SUPERUSER PASSWORD '12345';" &&\
    createdb -O postgres dbproject_base &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root
COPY --from=builder /app /app

EXPOSE 5000
ENV PGPASSWORD 12345
CMD service postgresql start && psql -h localhost -d dbproject_base -U art -p 5432 -a -q -f /app/db/db.sql && /app/main
#docker build -t art .
#docker run -p 5000:5000 --name art -t art
#./technopark-dbms-forum func -u http://localhost:5000/api -r report.html -k
#docker run -d --memory 2G --log-opt max-size=5M --log-opt max-file=3 --name art -p 5000:5000 art
