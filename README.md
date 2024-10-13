# Filemanager

## Description

[protos repository](https://github.com/Goose47/go-grpc-filemanager.protos)

A go gRPC microservice to manage binary files.

Service limits connections: 100 concurrent connections for unary calls. 10 concurrent connections on stream calls.
Limits are configured in config files.

## Steps to install

- adjust your configuration
- build Docker image
- run Docker container:
```
docker run --name {container-name} \
-v {storage-path}:{app-storage-path} \
-e CONFIG_PATH={path-to-config-file} \
-p 44044:44044 \
{image-name}
```

## Procedure implementations

#### ListFiles
ListFiles uses os.Stat to retrieve file names, modification and change timestamps. Change timestamp is used as an approximation to created time, as linux does not have utils to retrieve file creation date.

#### File 
Server-stream rpc that returns byte contents of specified file.

#### Upload 
Client-stream rpc that saves file, uploaded by client.

