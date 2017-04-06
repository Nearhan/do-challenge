# Coding Challenge


## Design Rationale







## Dockerfile run


I've included two docker files.


The default Dockerfile assumes you have go installed correctly.
This uses the scatch image and produces a very small image because all
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


### Do not have Go installed


1. Build the image

```
docker build -f Dockerfile.extra . -t do-challenge
```

2. Run the container

```
docker build -f Dockerfile.extra . -t do-challenge
```





