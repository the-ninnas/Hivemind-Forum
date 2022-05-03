# Base image
FROM golang

# Create directory inside the image
WORKDIR /forum

# Copies all the files to /forum
COPY . /forum

# Build the application
RUN go build -o main cmd/*

EXPOSE 8080

CMD [ "/forum/main"]