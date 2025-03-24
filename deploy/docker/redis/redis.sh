# redis standalone node
# mkdir -p C:\docker\redis\data
# mkdir -p C:\docker\redis
docker run -p 6379:6379 --name redis -v C:\docker\redis\data:/data -v C:\docker\redis\conf\redis.conf:/etc/redis/redis.conf -d redis redis-server /etc/redis/redis.conf