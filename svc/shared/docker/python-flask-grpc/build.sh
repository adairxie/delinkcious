IMAGE=ifangfang/delinkcious-python-flask-grpc:0.1
docker build . -t $IMAGE
docker push $IMAGE