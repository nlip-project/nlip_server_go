# NLIP Go Server

Welcome to the NLIP Go Server! This project is a basic **Go** implementation of NLIP server side protocol.

This package provides an example implementation in Go using the Echo web framework.

## Quickstart

1. Clone the project
```bash
git clone git@github.com:nlip-project/nlip_server_go.git
```

2. View and modify the included `.env_example` file. This is an example `.env` configuraiton needed for the project. Rename it to `.env`. Inside this file:
   - Required variables are necessary for NLIP to run as intended. `PORT`, `CERT_FILE` and  `KEY_FILE` (for HTTPS), `EXECUTABLE LOCATION`, `UPLOAD_PATH` must be provided.
   - Optional variables are not necessary. Variables included in the `.env` file are used to streamline the deployment and running of the server smoother on a MacOS server. These can be safely removed.
   - This version uses  a Custom OAuth server and Google's OAuth server for authentication. The Custom client can be found in the same organization repository, under `nlip-iam`. Please follow the instructions there to set up the custom OAuth server.

3. After you are sure that the `.env` file is correctly set up, you can use the script located at `scripts/deploy.sh` to automate your deployment process.
   - **IMPORTANT**: If you are not using any of the **OPTIONAL** variables, you need to go into `scripts/deploy.sh` to remove the section between `#### TODO` and `#### END_TODO`. Not doing so will cause the script to fail.

## Creating a launch service

#### MacOS

Lorem Ipsum

#### Linux

Lorem Ipsum

## Documentation

Lorem Ipsum

## Endpoints

#### `/nlip/`

Accepted formats:
1. Text
```
{
    "format": “text”,
    "subformat": “english”,
    "content": "Tell me a fun fact",
}
```
2. Binary
```
{
    "format": “text”,
    "subformat": “english”,
    "content": "Describe this picture",
    "submessages": [
        {
            "format": "binary",
            "subformat": "jpeg",
            "content": <base-64-encoding>
        }
    ]
}
```

## Pitfalls

Lorem Ipsum

## TODO

Lorem Ipsum
