# CNFUT (Cloud Native File/Folder Upload Tool)

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/necais/cnfut/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/necais/cnfut/tree/main)
[![codecov](https://codecov.io/gh/necais/cnfut/branch/main/graph/badge.svg?token=GAZ72S3I2J)](https://codecov.io/gh/necais/cnfut)
 [![Sonarcloud Status](https://sonarcloud.io/api/project_badges/measure?project=cnfut&metric=alert_status)](https://sonarcloud.io/dashboard?id=cnfut) 
 [![SonarCloud Bugs](https://sonarcloud.io/api/project_badges/measure?project=cnfut&metric=bugs)](https://sonarcloud.io/component_measures/metric/reliability_rating/list?id=cnfut)
 [![SonarCloud Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=cnfut&metric=vulnerabilities)](https://sonarcloud.io/component_measures/metric/security_rating/list?id=cnfut)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=plastic)](https://opensource.org/licenses/MIT)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/necais/cnfut?style=plastic)


qa34hj7wj9d4

CNFUT is a cloud native solution for copying files and directories between public clouds and file systems. 

Supported systems are: S3, Google, Azure and Local
    
    S3 supported systems
    Google cloud storage
    Azure Blob Storage
    Local file copies
    



## Deploy to kubernetes

 [necais/cnfut](https://hub.docker.com/r/necais/cnfut) could be used for deployment in kubernetes. Instructions for sample deployment can be found 
 [here](https://docs.docker.com/get-started/kube-deploy/):

## Usage
----
  Copies file between two systems

* **URL**

  /api/v1/

* **Method:**
  
   `POST`

* **Data Params**

  | Parameter            | Explanation                                       | Type    | Default values | Example                             | Mandatory  |
   ----------------------|---------------------------------------------------|---------|----------------|-------------------------------------|------------|
  | source               | Where is the source of data located(file, folder) | String  | None           | /data/data3                         | Yes        |
  | destination          | Where to copy the data(file or folder)            | String  | None           | /data/                              | Yes        |
  | sourceType           | Source type: s3, google, azure, local             | String  | None           | local                               | Yes        |
  | destinationType      | Destination type: s3, google, azure, local        | String  | None           | s3                                  | Yes        |
  | Concurrent           | Should copy executed concurrently                 | Boolean | None           | true                                | No         |
  | Region               | Region of the resource(AWS Region)                | String  | us-east-1      | us-east-1                           | No         |
  | Bucket               | Bucket name in the cloud                          | String  | None           | data                                | No         |
  | GoogleCredentialPath | Credential file path for Google                   | String  | None           | /home/necais                        | No         |
  | S3AccessKeyId        | S3 access key                                     | String  | None           | ieuwqhefbdsnbfs                     | No         |
  | S3SecretAccessKey    | S3 secret key                                     | String  | None           | fdsgdfbcvbchghf                     | No         |
  | Endpoint             | URL to endpoint                                   | String  | None           | https://cnfut.blob.core.windows.net | No         |


* **Success Response:**

  * **Code:** 202 Accepted <br />
    **Content:** `{ status : Accepted }`
 
* **Error Response:**

  * **Code:** 400 Bad Request <br />
    **Content:** `{ error : "Fields are not provided" }`

* **Sample Call:**

  ``` 
  {
  "source": "/data/data/data3",
  "destination": "data1",
  "sourceType": "local",
  "destinationType": "azure",
  "bucket": "cnfut-container",
  "endpoint": "https://cnfut.blob.core.windows.net"
  } 
  ```
