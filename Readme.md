# Coding Challenge


## Design Rationale


Initially when I wrote the TCP server I had come up with a data structure like this 

```

{
    "Package_Name": {
        Deps : [] // dependencies for this package
        ReqBy: [] // required by other packages
    }
}

```

I changed my mind however since the Index and Remove operations become pretty large and pretty complex.
And hard to read. Also it didn't run as fast as I thought it would.


So I opted out for a simpler data model.


```
{
    "Package_Name": [] // deps for the package
}
```

Which meant that I use raw iteration to Remove and Add dependencies.
This is less efficient then using the other data structure but my code was much simpler to read
and understand. If we were able to use some libraries, I would have opted out to use leveldb
or some other embedded database for all that nice persistence stuff.


## Tests

I've added integration and unit tests.

server_test.go is where I spin up a real server and test its outputs, the "integration" tests

All other tests are simple unit tests.


## Dockerfile run


I've included two docker files.


The default Dockerfile assumes you have go installed correctly.
This uses the scratch image and produces a very small image because all
that is in the container is the compiled go binary


The other file, Dockerfile.extra makes no assumptions at all except that you have docker installed
This image is far larger.


### Have Go installed

1. First Build using this command

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main
```
2. Then build image

```
docker build -t do-challenge .
```

3. Run the container!

```
docker run -p "8080:8080" do-challenge
```


### Do NOT have Go installed


1. Build the image

```
docker build -f Dockerfile.extra . -t do-challenge
```

2. Run the container

```
docker run -p "8080:8080" do-challenge
```