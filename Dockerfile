FROM ubuntu:18.04

ADD golang /usr/src/myapp/golang
ADD php /usr/src/myapp/php
ADD kafka /usr/src/myapp/kafka
copy . /usr/src/myapp/

RUN apt-get -y update && apt-get install -y php7.0
RUN apt-get -y install openjdk-8-jre

RUN php -r "readfile('http://getcomposer.org/installer');" | php -- --install-dir=/usr/bin/ --filename=composer;
WORKDIR /usr/src/myapp/php
RUN composer install

EXPOSE 8000

WORKDIR /usr/src/myapp
ENTRYPOINT ["/bin/sh", "startup.sh"]
