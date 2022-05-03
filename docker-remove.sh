# stop the container
docker stop forumcontainer

# delete container
docker rm forumcontainer

#delete image
docker rmi forumimage
echo
# show that they are deleted
docker ps -a
docker images