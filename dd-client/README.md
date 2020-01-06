# Generic DD-client 
- Purpose:
```
To write data on volume mount point of Kubernetes applications.
```
- Compatibility:
``` 
This client is only compatible for application which have dd-utility installed.
```
Pre-requisites:
```
- Mandatory ENVs such as [BLOCK_SIZE, COUNT, NAMESPACE, MOUNT_POINT, APP_LABEL, RETRY_DURATION, RETRY_COUNT] should be present in job spec (dd-client.yml) before launching client 
``` 