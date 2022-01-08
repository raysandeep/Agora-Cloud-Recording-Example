docker build -t videoapp .
docker run -d --name videoapp --env-file .env -p 80:3000 videoapp