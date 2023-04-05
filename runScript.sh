docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker-compose down
docker rmi $(docker images -f dangling=true -q)
docker volume rm $(docker volume ls -q)

docker build -t art -f Dockerfile .
docker run -d -p 5000:5000 --name test1 -t art
docker-compose up -d --build