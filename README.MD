*NB! Buttons with a * in their name are not operatale and are placeholders for future functions or aesthetic purposes.

Image upload:
* Create a category
* Go to "Create Post" http://localhost:8080/post to make a post with an image.

Usage with Docker:
* Type "sh docker-build.sh" to your terminal
* Go to localhost:8080 from your browser
* To delete docker image and container, type "sh docker-remove.sh' to your terminal

Usage with terminal:
* go run cmd/*

Usage of OAUTH login:
* GITHUB: 
Go to github.go, change githubClientID and githubClientSecret on line 14 & 15. How to get clientID and Secret? 
- Log in to github.com
- Go to settings > Developer settings > OAuth Apps > New OAuth App
- Application name: whatever..
- Homepage URL: http://localhost:8080/
- Authorization callback URL: http://localhost:8080/githubcallback

Authors: Anna-Maria, Gert, Mihkel, Silver
