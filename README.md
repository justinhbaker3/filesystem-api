# Filesystem API
An API used to interact with a remote filesystem

## Building and running
To create a docker image named filesystem-api:
```
make build-docker
```

To run the above image:
```
make run-docker
```

Then, navigate to localhost:8080 for a list of files and directories the API has access to. By default, the docker container only exposes files and directories present within the container (everything in this repository is copied into the image when it is built). 

Users who want to mount a local filesystem should do so at their own risk. To run the container with a mount to your current local directory:
```
docker run -p 8080:8080 --mount type=bind,source="$(pwd)",target=/app/filesystem,readonly -e FILESYSTEM_API_DIRECTORY=/app/filesystem filesystem-api
```

## API Documentation

GET Request:
```
curl localhost:8080/path/to/file/or/directory
```

Get Response Example (directory):
```
{
  "is_directory": true,
  "directory_content": [
    {
      "is_directory": true,
      "name": ".git",
      "size": 4096,
      "permissions": "-rwxr-xr-x"
    },
    {
      "is_directory": false,
      "name": ".gitignore",
      "size": 14,
      "permissions": "-rw-r--r--"
    },
  ]
}
```

Get Response Example (file):
```
{
    "is_directory":false,
    "file_body":"filesystem-api"
}
```

## Testing
To run all unit tests:
```
make test
```
