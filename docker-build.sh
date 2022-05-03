# build image
docker build -t forumimage .

# build container
docker container run -p 8080:8080 --detach --name forumcontainer forumimage
echo

# display images
docker images

# display containers
docker ps -a