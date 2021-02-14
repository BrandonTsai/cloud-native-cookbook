FROM nginx:1.14.2


RUN chgrp -R 0 /etc/nginx/ /var/cache/nginx /var/run /var/log/nginx  && \ 
  chmod -R g+rwX /etc/nginx/ /var/cache/nginx /var/run /var/log/nginx

# users are not allowed to listen on priviliged ports
RUN sed -i.bak 's/listen\(.*\)80;/listen 8081;/' /etc/nginx/conf.d/default.conf
EXPOSE 8080

# comment user directive as master process is run as user in OpenShift random UID
RUN sed -i.bak 's/^user/#user/' /etc/nginx/nginx.conf

EXPOSE 8080
