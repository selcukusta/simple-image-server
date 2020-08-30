# Host Your Own Image Server based on MongoDB, Azure Blob Storage or Google Drive

[![Go Report Card](https://goreportcard.com/badge/github.com/selcukusta/simple-image-server)](https://goreportcard.com/report/github.com/selcukusta/simple-image-server)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/selcukusta/simple-image-server)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/selcukusta/simple-image-server/blob/master/LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fselcukusta%2Fsimple-image-server.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fselcukusta%2Fsimple-image-server?ref=badge_shield)
[![codecov](https://codecov.io/gh/selcukusta/simple-image-server/branch/master/graph/badge.svg)](https://codecov.io/gh/selcukusta/simple-image-server)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=selcukusta_gdrive-image-server&metric=alert_status)](https://sonarcloud.io/dashboard?id=selcukusta_gdrive-image-server)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=selcukusta_gdrive-image-server&metric=sqale_index)](https://sonarcloud.io/dashboard?id=selcukusta_gdrive-image-server)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=selcukusta_gdrive-image-server&metric=bugs)](https://sonarcloud.io/dashboard?id=selcukusta_gdrive-image-server)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=selcukusta_gdrive-image-server&metric=code_smells)](https://sonarcloud.io/dashboard?id=selcukusta_gdrive-image-server)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=selcukusta_gdrive-image-server&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=selcukusta_gdrive-image-server)

## Getting Started

### Supported Platforms

- **Google Drive**

- **Azure Blob Storage**

- **MongoDB GridFS**

### Supported Formats

- **image/jpeg**
- **image/png**
- **image/webp**

#### 💻 Google Drive

You need to have a Google account to set up the project.

##### Setup steps for using Google Drive API

Firstly, create a new project from [Google Developer Console](https://console.developers.google.com/). Go to the **Credentials** menu and create a new **Service Accounts** credential. It might be named as `[PROJECT_ALIAS]-xxxx-xxxxxxx.json`.

Download it and rename as `gcloud-image-server-cred.json`. Put the file to the **root** folder (_it will be used for building Docker image_).

Copy your service account mail address (_it will be used for sharing your images with the project_).

##### Setup steps for host the images

Go to your Drive page and create a folder, ie. `image-server`.

Share the folder with your service account (_was copied before_). Sharing rule will be applied to the all sub-items in the folder.

As a last step, upload any image (mime should be `image/jpeg` or `image/png`) to the folder and get the ID.

_NOTE: `ID` is not showing anywhere in the portal, it sucks! To catch it, right click your item and select `Get Shareable Link`. Copy the last part of it, and toggle off the sharable link feature._

#### 💻 Azure Blob Storage

##### Setup steps for using Azure Blob Storage

If you have an Azure Account and blob storage subscription, you have to create a new Access Key from portal or CLI. It could be like that; `DefaultEndpointsProtocol=https;AccountName=[YOUR_STORAGE_NAME];AccountKey=[YOUR_ACCOUNT_KEY]==;EndpointSuffix=core.windows.net.`

You need to add these values to the environment:

| Name               | Type     |
| :----------------- | :------- |
| `ABS_ACCOUNT_KEY`  | `string` |
| `ABS_ACCOUNT_NAME` | `string` |
| `ABS_AZURE_URI`    | `string` |

##### Setup steps for host the images

You can create a new container from Blob service > Container menus. Assume that you have a container which is named as sample-photos. It has two directories and a the picture at the last directory (_summer > hotels > swimming.jpg_).

Reach your blob with the url: `http://127.0.0.1:8080/i/abs/100/400x0/sample-photos/summer/hotels/swimming.jpg`

#### 💻 MongoDB

Before running the application set these environment variables (or use `Docker` image, run with `docker container run -d --name mongodb-instance -p 27017:27017 mongo:3.6.18-xenial` and leave them default):

| Name                   | Type     | Default Value               |
| :--------------------- | :------- | :-------------------------- |
| `MONGO_CONNECTION_STR` | `string` | _mongodb://127.0.0.1:27017_ |
| `MONGO_DB_NAME`        | `string` | _Photos_                    |
| `MONGO_MAX_POOL_SIZE`  | `uint64` | _5_                         |

Run `go run cmd/mongo-seed/main.go` command and create 3 sample record on the DB such as;

| \_id                     | chunkSize | filename    | length | metadata                       | uploadDate               |
| :----------------------- | :-------- | :---------- | :----- | :----------------------------- | :----------------------- |
| 5ec684803dd893bb72ead932 | 261120    | image-3.jpg | 710863 | {"Content-Type": "image/jpeg"} | 2020-05-21T13:39:12.585Z |
| 5ec684803dd893bb72ead931 | 261120    | image-2.jpg | 516218 | {"Content-Type": "image/jpeg"} | 2020-05-21T13:39:12.603Z |
| 5ec684803dd893bb72ead930 | 261120    | image-1.jpg | 379373 | {"Content-Type": "image/jpeg"} | 2020-05-21T13:39:12.617Z |

## Running

### Via Docker

#### Build

```bash
# Build the image with Google Drive and WebP support. gcloud-image-server-cred.json file should be included!
docker image build -t [YOUR_REPOSITORY]/w-gdrive-w-webp:1.0.0 -f w-gdrive-w-webp.Dockerfile .
# Build the image with Google Drive without WebP support. gcloud-image-server-cred.json file should be included!
docker image build -t [YOUR_REPOSITORY]/w-gdrive-wo-webp:1.0.0 -f w-gdrive-wo-webp.Dockerfile .
# Build the image without Google Drive and with WebP support.
docker image build -t [YOUR_REPOSITORY]/wo-gdrive-w-webp:1.0.0 -f wo-gdrive-w-webp.Dockerfile .
# Build the image without Google Drive and withoutt WebP support.
docker image build -t [YOUR_REPOSITORY]/wo-gdrive-wo-webp:1.0.0 -f wo-gdrive-wo-webp.Dockerfile .
```

#### Run

```bash
docker container run -p 8080:8080 --rm [YOUR_REPOSITORY]/w-gdrive-w-webp:1.0.0
docker container run -p 8080:8080 --rm [YOUR_REPOSITORY]/w-gdrive-wo-webp:1.0.0
docker container run -p 8080:8080 --rm [YOUR_REPOSITORY]/wo-gdrive-w-webp:1.0.0
docker container run -p 8080:8080 --rm [YOUR_REPOSITORY]/wo-gdrive-wo-webp:1.0.0
```

### Via Docker Compose

- `docker-compose up`

- Go to your favorite browser

- Surf to the URLs such as below.

Three containers will be bringing up. The first one is your image server application. It was written with purely **[the Go programming language](https://golang.org/)**. It can be reached with _[http://localhost:8080/version](http://localhost:8080/version)_

Another one is **[Varnish HTTP Cache](https://varnish-cache.org/)**. It will automatically cache the output of the response for 14 days. It can be reached with _[http://localhost:8081/version](http://localhost:8081/version)_

The latest one, **[nuster](https://github.com/jiangwenyuan/nuster)**. Nuster, a high-performance HTTP proxy cache server and RESTful NoSQL cache server based on HAProxy. It will automatically cache the output of the response for 14 days like Varnish. It can be reached with _[http://localhost:8082/version](http://localhost:8082/version)_

_You can choose any cache server according to your experience._

## Usage

/i/ endpoint is used for image operations, has four different usages:

### Google Drive

```
/i/gdrive/webp/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{options:opt}/{*id}

/i/gdrive/webp/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{*id}

/i/gdrive/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{options:opt}/{*id}

/i/gdrive/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{*id}
```

### Azure Blob Storage

```
/i/abs/webp/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{options:opt}/{*id}

/i/abs/webp/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{*id}

/i/abs/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{options:opt}/{*id}

/i/abs/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{*id}
```

### MongoDB

```
/i/gridfs/webp/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{options:opt}/{*id}

/i/gridfs/webp/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{*id}

/i/gridfs/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{options:opt}/{*id}

/i/gridfs/{quality:range(0,100)}/{w:range(0,5000)}x{h:range(0,5000)}/{*id}
```

## WebP Support

If you add `/webp/` path to the URL, you can get the image as webp.

## Options

| Option | Description                                                                                                                 |
| ------ | --------------------------------------------------------------------------------------------------------------------------- |
| `g`    | This will convert the given image into a grayscale image.                                                                   |
| `t`    | This will scales the image up or down, crops it to the specified width and hight and returns the transformed image.         |
| `c`    | This will cuts out a rectangular region with the specified size from the center of the image and returns the cropped image. |

## Samples

- **Original size:** `http://127.0.0.1:8080/i/gdrive/100/0x0/[YOUR_FILE_ID]`, `http://127.0.0.1:8080/i/gridfs/100/0x0/[MONGODB_OBJECT_ID]` or `http://127.0.0.1:8080/i/abs/100/0x0/[YOUR_STORAGE_PATH]`

![1](assets/1.png)

- **Resize with aspect ratio:** `http://127.0.0.1:8080/i/gdrive/100/500x0/[YOUR_FILE_ID]`, `http://127.0.0.1:8080/i/gridfs/100/500x0/[MONGODB_OBJECT_ID]` or `http://127.0.0.1:8080/i/abs/100/500x0/[YOUR_STORAGE_PATH]`

![2](assets/2.png)

- **Less quality:** `http://127.0.0.1:8080/i/gdrive/1/0x0/[YOUR_FILE_ID]`,`http://127.0.0.1:8080/i/gridfs/1/0x0/[MONGODB_OBJECT_ID]` or `http://127.0.0.1:8080/i/abs/1/0x0/[YOUR_STORAGE_PATH]`

![3](assets/3.png)

- **Resize without aspect ratio:** `http://127.0.0.1:8080/i/gdrive/100/1600x600/[YOUR_FILE_ID]`, `http://127.0.0.1:8080/i/gridfs/100/1600x600/[MONGODB_OBJECT_ID]` or `http://127.0.0.1:8080/i/abs/100/1600x600/[YOUR_STORAGE_PATH]`

![4](assets/4.png)

- **Resize with crop:** `http://127.0.0.1:8080/i/gdrive/100/1600x600/c/[YOUR_FILE_ID]`, `http://127.0.0.1:8080/i/gridfs/100/1600x600/c/[MONGODB_OBJECT_ID]` or `http://127.0.0.1:8080/i/abs/100/1600x600/c/[YOUR_STORAGE_PATH]`

![5](assets/5.png)

- **Create thumbnail with aspect ratio:** `http://127.0.0.1:8080/i/gdrive/100/0x300/t/[YOUR_FILE_ID]`, `http://127.0.0.1:8080/i/gridfs/100/0x300/t/[MONGODB_OBJECT_ID]` or `http://127.0.0.1:8080/i/abs/100/0x300/t/[YOUR_STORAGE_PATH]`

![6](assets/6.png)

- **Grayscale:** `http://127.0.0.1:8080/i/gdrive/100/900x0/g/[YOUR_FILE_ID]`, `http://127.0.0.1:8080/i/gridfs/100/900x0/g/[MONGODB_OBJECT_ID]` or `http://127.0.0.1:8080/i/abs/100/900x0/g/[YOUR_STORAGE_PATH]`

![7](assets/7.png)

## LICENSE

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fselcukusta%2Fgdrive-image-server.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fselcukusta%2Fgdrive-image-server?ref=badge_large)
