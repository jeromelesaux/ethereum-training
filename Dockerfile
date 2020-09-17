FROM golang:1.14


# build : 
# docker build . -f Dockerfile -t certification --build-arg GOOGLE_ID= --build-arg GOOGLE_SECRET=
#  --build-arg AWS_ACCESS_KEY_ID= --build-arg AWS_SECRET_ACCESS_KEY=

ARG GOOGLE_ID
ARG GOOGLE_SECRET
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY

ENV GOOGLE_ID=$GOOGLE_ID
ENV GOOGLE_SECRET=$GOOGLE_SECRET
ENV AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
ENV AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY


# get the source code of the app
RUN go get -v github.com/jeromelesaux/ethereum-training

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/jeromelesaux/ethereum-training

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080:8080

# Create config file 
RUN echo '{\n\
    "directorysavepath": "/tmp/tmpfiles",\n\
    "ethereumendpoint": "https://ropsten.infura.io/v3/7de903803c31428bbdd1186107a2d660",\n\
    "privatekey": "48218b47d9afba13df85e4b29e4e0bb73ae526cdebb316738832be607e7c7174",\n\ 
    "serverurl": "https://btrust.labinno.fr",\n\
    "uselocalstorage" :false\n\
}\n'\
> /etc/config.json

RUN ls /etc
# Run the executable
CMD ["ethereum-training", "-config", "/etc/config.json"]

# run docker container : 
# docker run -p 8080:8080 dockername
