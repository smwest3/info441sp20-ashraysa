docker rm -f gateway

docker pull ashmann7/etvgateway

# docker network rm customNetwork

docker network create customNetwork

docker run -d \
--network customNetwork \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=/etc/letsencrypt/live/api.eradicatethevape.live/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.eradicatethevape.live//privkey.pem \
-e SESSIONKEY="mystring" \
-e REDISADDR="redisServer:6379" \
-e DSN="root:password@tcp(db:3306)/db" \
-e MARKETPLACEADDR="marketplace:5200" \
-e THREADSADDR="threads:5300" \
-e PROGRESSADDR="progress:5400" \
-p 443:443 \
--name gateway \
ashmann7/etvgateway

exit

