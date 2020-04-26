# microservice_example
# C-Users-jbrad-go-src-github.com-jlb922-microservice_example
# based off https://www.melvinvivas.com/my-first-go-microservice/
# TODO - incorporate https://medium.com/@adiach3nko/package-management-with-go-modules-the-pragmatic-guide-c831b4eaaf31
# Also - https://www.callicoder.com/docker-golang-image-container-example/

Building the image

From the source directory run:
docker build -t go-docker-optimized -f Dockerfile .


Check the size of the image:
docker image ls

Run just the Golang image
docker run -d -p 9090:9090 -v  go-docker-optimized

Run the Golang and Mongo images:
docker-compose up -d

Rebuild and run both images if needed
docker-compose up -d --build --force-recreate

Inspect the network specified in the docker-compose.yml
docker network inspect

Run the mongo image only, set ENV to dev to run Golang code outside image
docker-compose -f docker-compose-dev.yml up -d
export MGOHOSTNAME=localhost




