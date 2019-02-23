FROM ubuntu
RUN apt-get update && \
	apt-get -y dist-upgrade && \
	apt-get -y install nginx nodejs npm && \
	apt-get -y autoremove && \
	apt-get clean

COPY ./ /var/tmp/build/
WORKDIR /var/tmp/
RUN cd ./build && \
	npm install && \
	npm run build && \
	mv index.html /var/www/html/ && \
	mv ./dist /var/www/html/ && \
	rm -f /var/www/html/index.nginx-debian.html && \
	cd .. && rm -fr ./build

WORKDIR /var/www

EXPOSE 80
COPY nginx.conf /etc/nginx/sites-available/default
VOLUME /var/www/html/data
CMD /usr/sbin/nginx -g 'daemon off;'

