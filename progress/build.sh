GOOS=linux go build

docker build -t ashmann7/progress .

go clean

docker push ashmann7/progress

ssh ec2-user@ec2-54-82-153-19.compute-1.amazonaws.com < deploy.sh