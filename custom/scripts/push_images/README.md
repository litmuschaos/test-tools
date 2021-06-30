## Table of content
- [Supported Tunables](https://github.com/litmuschaos/test-tools/blob/master/custom/scripts/push_images/README.md#supported-tunables)
- [Check the available commands](https://github.com/litmuschaos/test-tools/blob/master/custom/scripts/push_images/README.md#check-the-available-commands)
- [List Down all the Images](https://github.com/litmuschaos/test-tools/blob/master/custom/scripts/push_images/README.md#list-down-all-the-images)
- [Pull the LitmusChaos Images into your machine](https://github.com/litmuschaos/test-tools/blob/master/custom/scripts/push_images/README.md#pull-the-litmuschaos-images-into-your-machine)
- [Push the LitmusChaos Images to your image registry](https://github.com/litmuschaos/test-tools/blob/master/custom/scripts/push_images/README.md#push-the-litmuschaos-images-to-your-image-registry)


## LitmusChaos Images

- Push docker images contains all the images that are used to execute a litmuschaos generic experiment using litmus portal. For more information please check [LitmusChaos Docs](https://litmusdocs-beta.netlify.app/docs/introduction).


**LitmusChaos images in the script**

 <table>
    <tr>
      <th> Portal Images </th>
      <th> Backend and monitoring Images  </th>
      <th> Workflow and other Images </th>
  </tr>
  <tr>
    <td> 
      <ul>
         <li>litmuschaos/litmusportal-frontend</li>
         <li>litmuschaos/litmusportal-server</li>
         <li>litmuschaos/litmusportal-event-tracker</li>
         <li>litmuschaos/litmusportal-auth-server</li>
         <li>litmuschaos/litmusportal-subscriber</li>
      </ul>
    </td>
    <td>
      <ul>
         <li>litmuschaos/chaos-operator</li>
         <li>litmuschaos/chaos-runner</li>
         <li>litmuschaos/go-runner</li>
         <li>litmuschaos/chaos-exporter</li>
      </ul>    
    </td> 
    <td>
      <ul>
         <li>litmuschaos/k8s:latest</li>
         <li>litmuschaos/litmus-checker:latest</li>
         <li>litmuschaos/workflow-controller:v2.11.0</li>
         <li>litmuschaos/argoexec:v2.11.0</li>
         <li>litmuschaos/mongo:4.2.8</li>         
      </ul>      
    </td>
  </tr>
 </table>

## Get LitmusChaos Images In Your Repository

- For pulling the litmus image and pushing into your registry please follow the given steps:

```bash
wget https://raw.githubusercontent.com/litmuschaos/test-tools/master/custom/scripts/push_images/litmus_image_push.sh

chmod +x litmus_image_push.sh
```
#### Supported Tunables

 <table>
    <tr>
      <th> Variables </th>
      <th> Description </th>
      <th> Specify </th>
      <th> Notes </th>
  </tr>
  <tr>
    <td> TARGET_REPONAME </td>
    <td> Provide the name of target repo-name for the image. Example <code>litmuschaos</code> </td>
    <td> Mandatory </td>
    <td> No default value is provided</td>
  </tr>
  <tr>
    <td> TARGET_IMAGE_REGISTRY </td>
    <td>  Provide the name of target image registry. </td>
    <td> Optional </td>
    <td> Default to <code>docker.io</code></td>
  </tr>
  <tr>
    <td> LITMUS_PORTAL_TAG </td>
    <td> Provide the tag for portal components like litmuschaos-forntend and others</td>
    <td> Optional </td>
   <td> If LITMUS_PORTAL_TAG is not provided ,By-default It will select the latest release tag </td>
  </tr>
  <tr>
    <td> LITMUS_BACKEND_TAG </td>
    <td> Provide the tag for litmus backend components like chaos-operator,chaos-runner,go-runner and others</td>
    <td> Optional </td>
   <td> If LITMUS_BACKEND_TAG is not provided ,By-default It will select the latest release tag </td>
  </tr>
  <tr>
    <td> LITMUS_IMAGE_REGISTRY </td>
    <td> Provide the litmuschaos image registry </td>
    <td> Optional </td>
    <td> If not provided, it will use <code>docker.io</code> as default registry.</td>
  </tr>
 </table>

#### Check the available commands:

```bash
$ ./litmus_image_push.sh -h

Usage:       ./litmus_image_push.sh ARGS (list|pull|tag|push)

list:        "./litmus_image_push.sh list" will list all the images used by the litmus portal.     


pull:        "./litmus_image_push.sh pull" will pull the litmus images with the given image tag. 
              The value of tag can be provided by exporting following ENVs:
              - LITMUS_PORTAL_TAG: Tag for the portal component like 'litmusportal-frontend' etc
              - LITMUS_BACKEND_TAG: Tag for backend component chaos-operator, chaos-runner, go-runner etc
              - IMAGE_REGISTRY: Name of litmuschaos image registry. Default is docker.io
              The default images tags are the latest tags released.

push:         "./litmus_image_push.sh push" will push the images to the given target image registry with the give repo-name
              The value of target images can be set by exporting following ENVs:
              - TARGET_IMAGE_REGISTRY: Give the target registry name. Default is set to "docker.io"
              - TARGET_REPONAME: Give the target image repo-name. This is mandatory to provide.               

```

#### List Down all the Images

Format:
```
./litmus_image_push.sh list
```
Example:
```bash
$ ./litmus_image_push.sh list

portal component images ...
1. docker.io/litmuschaos/litmusportal-frontend:2.0.0-Beta8
2. docker.io/litmuschaos/litmusportal-server:2.0.0-Beta8
3. docker.io/litmuschaos/litmusportal-event-tracker:2.0.0-Beta8
4. docker.io/litmuschaos/litmusportal-auth-server:2.0.0-Beta8
5. docker.io/litmuschaos/litmusportal-subscriber:2.0.0-Beta8

backend component images ...
6. docker.io/litmuschaos/chaos-operator:1.13.6
7. docker.io/litmuschaos/chaos-runner:1.13.6
8. docker.io/litmuschaos/chaos-exporter:1.13.6
9. docker.io/litmuschaos/go-runner:1.13.6

workflow controller images ...
10. docker.io/litmuschaos/k8s:latest
11. docker.io/litmuschaos/litmus-checker:latest
12. docker.io/litmuschaos/workflow-controller:v2.11.0
13. docker.io/litmuschaos/argoexec:v2.11.0
14. docker.io/litmuschaos/mongo:4.2.8

```

#### Pull the LitmusChaos Images into your machine

Format:
```
./litmus_image_push.sh pull
```

Example: 

```bash
$ export LITMUS_BACKEND_TAG=1.13.5
$ ./litmus_image_push.sh pull
 Pulling docker.io/litmuschaos/litmusportal-frontend:2.0.0-Beta8
2.0.0-Beta8: Pulling from litmuschaos/litmusportal-frontend
9b794450f7b6: Pull complete 
f8fd2f03c36e: Pull complete 
3cb08bcbe78b: Pull complete 
052722e881c1: Pull complete 
2fd682de3a2e: Pull complete 
80ce77084de4: Pull complete 
1b5bab11387e: Pull complete 
a9916e96b367: Pull complete 
a5846cd3892d: Pull complete 
e389376c0ba7: Pull complete 
Digest: sha256:f602ad199ca4918e88483802b65b23719af3c3cfc5ee9dbf42834ea63343616d
Status: Downloaded newer image for litmuschaos/litmusportal-frontend:2.0.0-Beta8
docker.io/litmuschaos/litmusportal-frontend:2.0.0-Beta8

 Pulling docker.io/litmuschaos/litmusportal-server:2.0.0-Beta8
2.0.0-Beta8: Pulling from litmuschaos/litmusportal-server
540db60ca938: Pull complete 
ee72d46a452d: Pull complete 
cc7256672b8e: Pull complete 
5ea7a1e3d4e3: Pull complete 
Digest: sha256:bb8cc8cc96b848cc8e6390ccbbd1fcdc926fe1eeaa95808c0a6afc99a5b7d49e
Status: Downloaded newer image for litmuschaos/litmusportal-server:2.0.0-Beta8
docker.io/litmuschaos/litmusportal-server:2.0.0-Beta8

 Pulling docker.io/litmuschaos/litmusportal-event-tracker:2.0.0-Beta8
2.0.0-Beta8: Pulling from litmuschaos/litmusportal-event-tracker
df20fa9351a1: Pull complete 
9d89a673b031: Pull complete 
9a2277ec6d3b: Pull complete 
Digest: sha256:a7af56c759bc598cb94ab5725c5bbecf69a3b262d496db1664c45220746a93cc
Status: Downloaded newer image for litmuschaos/litmusportal-event-tracker:2.0.0-Beta8
docker.io/litmuschaos/litmusportal-event-tracker:2.0.0-Beta8

 Pulling docker.io/litmuschaos/litmusportal-auth-server:2.0.0-Beta8
2.0.0-Beta8: Pulling from litmuschaos/litmusportal-auth-server
540db60ca938: Already exists 
b05e72f6d6a4: Pull complete 
96dabbd5c86a: Pull complete 
Digest: sha256:1b75c8156e39b058a305ddccc7afe9a39ae7e8b0218f18b8dac47f599854c744
Status: Downloaded newer image for litmuschaos/litmusportal-auth-server:2.0.0-Beta8
docker.io/litmuschaos/litmusportal-auth-server:2.0.0-Beta8

 Pulling docker.io/litmuschaos/litmusportal-subscriber:2.0.0-Beta8
2.0.0-Beta8: Pulling from litmuschaos/litmusportal-subscriber
df20fa9351a1: Already exists 
7350b870b03c: Pull complete 
cb56a658641c: Pull complete 
Digest: sha256:70fa773e2776c609f78f4b657e17945259fd91b7a5df77369bd9760cefcf6e5c
Status: Downloaded newer image for litmuschaos/litmusportal-subscriber:2.0.0-Beta8
docker.io/litmuschaos/litmusportal-subscriber:2.0.0-Beta8


 Pulling docker.io/litmuschaos/chaos-operator:1.13.5
1.13.5: Pulling from litmuschaos/chaos-operator
8f403cb21126: Pull complete 
65c0f2178ac8: Pull complete 
d241b7cec9b6: Pull complete 
b01766494add: Pull complete 
f805af9d2af4: Pull complete 
Digest: sha256:446cc6c92f4fdc1ab1d203005a1b7a77f2d554251c1a651fec6db64c03551c5a
Status: Downloaded newer image for litmuschaos/chaos-operator:1.13.5
docker.io/litmuschaos/chaos-operator:1.13.5

 Pulling docker.io/litmuschaos/chaos-runner:1.13.5
1.13.5: Pulling from litmuschaos/chaos-runner
8f403cb21126: Already exists 
65c0f2178ac8: Already exists 
93446a7f9cb9: Pull complete 
Digest: sha256:a331c0170509956499468e959256de8f68be01a6557ee095b2b151d6a9215fa5
Status: Downloaded newer image for litmuschaos/chaos-runner:1.13.5
docker.io/litmuschaos/chaos-runner:1.13.5

 Pulling docker.io/litmuschaos/chaos-exporter:1.13.5
1.13.5: Pulling from litmuschaos/chaos-exporter
540db60ca938: Already exists 
90c3ae3c2987: Pull complete 
8b9b1a8f1fbf: Pull complete 
Digest: sha256:05319917b53d1d038419f63039a97dd9535036a000b095182134d8a6451f49e2
Status: Downloaded newer image for litmuschaos/chaos-exporter:1.13.5
docker.io/litmuschaos/chaos-exporter:1.13.5

 Pulling docker.io/litmuschaos/go-runner:1.13.5
1.13.5: Pulling from litmuschaos/go-runner
540db60ca938: Already exists 
9096e06f66ff: Pull complete 
5d88d0ed31a2: Pull complete 
1d2e856478f4: Pull complete 
e5e6f7ae8611: Pull complete 
202e5b265ffb: Pull complete 
e94aedf735d5: Pull complete 
d1f06d6c103c: Pull complete 
d2b2f65251b9: Pull complete 
47140343b874: Pull complete 
fa9284414937: Pull complete 
aa1f2ddc2bc9: Pull complete 
db9b992a082c: Pull complete 
Digest: sha256:e01be5bdbb83f7a6996a7633c09f9345ac6b99e0f20e0218b77948b17533a26c
Status: Downloaded newer image for litmuschaos/go-runner:1.13.5
docker.io/litmuschaos/go-runner:1.13.5


 Pulling docker.io/litmuschaos/k8s:latest
latest: Pulling from litmuschaos/k8s
188c0c94c7c5: Pull complete 
bcd76da4b3e7: Pull complete 
0cafdb09230a: Pull complete 
438d8a6bde50: Pull complete 
984ce75a6ffc: Pull complete 
ad5f83d0f934: Pull complete 
3b792e60a5cd: Pull complete 
9c186b046686: Pull complete 
aaf5d53c0c47: Pull complete 
Digest: sha256:768491682ca99e14498ddb62973a95a664a0525ca18e7ad0eaab2621d3124f5b
Status: Downloaded newer image for litmuschaos/k8s:latest
docker.io/litmuschaos/k8s:latest

 Pulling docker.io/litmuschaos/litmus-checker:latest
latest: Pulling from litmuschaos/litmus-checker
540db60ca938: Already exists 
3364ad582f2f: Pull complete 
e6fa4af51404: Pull complete 
Digest: sha256:3016763935fdbfccd2a5ccf2867550b2a8962a78ace73d957247b18ab619be18
Status: Downloaded newer image for litmuschaos/litmus-checker:latest
docker.io/litmuschaos/litmus-checker:latest

 Pulling docker.io/litmuschaos/workflow-controller:v2.9.3
v2.9.3: Pulling from litmuschaos/workflow-controller
1f407d3f644c: Pull complete 
1633c91701ba: Pull complete 
Digest: sha256:18cb5e9e8e4143ba40c9ad208141a6c61cbd8f62dcbc13646dc9d0f597908b6f
Status: Downloaded newer image for litmuschaos/workflow-controller:v2.9.3
docker.io/litmuschaos/workflow-controller:v2.9.3

 Pulling docker.io/litmuschaos/argoexec:v2.9.3
v2.9.3: Pulling from litmuschaos/argoexec
54fec2fa59d0: Pull complete 
ea1b58a62c8f: Pull complete 
baf5dcb8fbf2: Pull complete 
f7d968f2f223: Pull complete 
7be27d15fd7b: Pull complete 
1648478b409a: Pull complete 
220fc1c30111: Pull complete 
f102822e7b4e: Pull complete 
2ce8271a1891: Pull complete 
Digest: sha256:543ee6e910f5b7b81baf5b471776c533ef23000839b266d62504254e37a8471c
Status: Downloaded newer image for litmuschaos/argoexec:v2.9.3
docker.io/litmuschaos/argoexec:v2.9.3

 Pulling docker.io/litmuschaos/mongo:4.2.8
4.2.8: Pulling from litmuschaos/mongo
f08d8e2a3ba1: Pull complete 
3baa9cb2483b: Pull complete 
94e5ff4c0b15: Pull complete 
1860925334f9: Pull complete 
9d42806c06e6: Pull complete 
31a9fd218257: Pull complete 
5bd6e3f73ab9: Pull complete 
f6ae7a64936b: Pull complete 
a614d629c284: Pull complete 
477320af2dcc: Pull complete 
b8aab702fcf5: Pull complete 
b94c6a2dc294: Pull complete 
8cf889bdb7c6: Pull complete 
Digest: sha256:12070904286ea8d1c647a78512432cf74a874c75ebda033abdc01eb409a96092
Status: Downloaded newer image for litmuschaos/mongo:4.2.8
docker.io/litmuschaos/mongo:4.2.8
```

#### Push the LitmusChaos Images to your image registry

Format:
```bash
$ export TARGET_REPONAME=<target-repo-name>
$ ./litmus_image_push.sh push 
```

Example:
```bash
$ export TARGET_REPONAME=uditgaurav
$ export TARGET_IMAGE_REGISTRY=docker.io
$ ./litmus_image_push.sh push

The push refers to repository [docker.io/uditgaurav/litmusportal-frontend]
9d244c78ec15: Mounted from litmuschaos/litmusportal-frontend 
bb036397b05e: Mounted from litmuschaos/litmusportal-frontend 
029d0c195d83: Pushed 
db639965f176: Layer already exists 
19e5eaa18644: Layer already exists 
697f2aa6662e: Layer already exists 
1f9e2810747e: Layer already exists 
a3355a4d5656: Layer already exists 
50a03d8e0394: Layer already exists 
2b2bcc6e6724: Layer already exists 
2.0.0-Beta8: digest: sha256:1a49146424f1a246236de6bbae0d48dc34ebce04b1eede7c1341f5cf0bd6e928 size: 2402

The push refers to repository [docker.io/uditgaurav/litmusportal-server]
d4c0d26cf219: Mounted from litmuschaos/litmusportal-server 
1d6c2a1990e6: Mounted from litmuschaos/litmusportal-server 
815f9195e346: Mounted from litmuschaos/litmusportal-server 
b2d5eeeaba3a: Layer already exists 
2.0.0-Beta8: digest: sha256:1c73c3ca1c573018266a4064de7a314934a9243b9ef8e04157e11def57b99b55 size: 1156

The push refers to repository [docker.io/uditgaurav/litmusportal-event-tracker]
a68cd48cb97a: Mounted from litmuschaos/litmusportal-event-tracker 
560bd3ca6190: Mounted from litmuschaos/litmusportal-event-tracker 
50644c29ef5a: Layer already exists 
2.0.0-Beta8: digest: sha256:7c9e6966560a15a5f32a72fecc243113f6732305209047508ea67f74913283f9 size: 948

The push refers to repository [docker.io/uditgaurav/litmusportal-auth-server]
3813f7a52caa: Mounted from litmuschaos/litmusportal-auth-server 
b91cd2c5136a: Mounted from litmuschaos/litmusportal-auth-server 
b2d5eeeaba3a: Layer already exists 
2.0.0-Beta8: digest: sha256:ad36967c7849c8d2f7e141103cea3908b4bd10e31ab17d4cd9e846960615a1dd size: 947

The push refers to repository [docker.io/uditgaurav/litmusportal-subscriber]
d14193e0b3c6: Mounted from litmuschaos/litmusportal-subscriber 
589fe741e599: Mounted from litmuschaos/litmusportal-subscriber 
50644c29ef5a: Layer already exists 
2.0.0-Beta8: digest: sha256:c9770a89975525dabd039a6caea9064c7e1e4a09944809c909590805f0fc3958 size: 948

The push refers to repository [docker.io/uditgaurav/chaos-operator]
0c8769ae6ff6: Mounted from litmuschaos/chaos-operator 
6272ab06526e: Mounted from litmuschaos/chaos-operator 
e4ddbf1f92bf: Mounted from litmuschaos/chaos-operator 
144a43b910e8: Mounted from litmuschaos/chaos-runner 
4a2bc86056a8: Mounted from litmuschaos/chaos-runner 
1.13.5: digest: sha256:46f894e2ceb1b516503269e0769b5b65b38e42048a874eaaaf592f688e6557d8 size: 1363

The push refers to repository [docker.io/uditgaurav/chaos-runner]
48d6a05f0d98: Mounted from litmuschaos/chaos-runner 
144a43b910e8: Layer already exists 
4a2bc86056a8: Layer already exists 
1.13.5: digest: sha256:e4466b4f94b9c201d994f9242bf34d7fdb679c92788127495cded54084fe0223 size: 949

The push refers to repository [docker.io/uditgaurav/chaos-exporter]
1577184deaf0: Mounted from litmuschaos/chaos-exporter 
4d3251350b3b: Mounted from litmuschaos/chaos-exporter 
b2d5eeeaba3a: Layer already exists 
1.13.5: digest: sha256:4d20b0632e092d48ede05f34d8a9d631289ff7f3c09b52f9e5eacd68f53980d6 size: 948

The push refers to repository [docker.io/uditgaurav/go-runner]
271879a892e9: Mounted from litmuschaos/go-runner 
955ea34328df: Mounted from litmuschaos/go-runner 
80d266cd9492: Mounted from litmuschaos/go-runner 
31091d26e39a: Mounted from litmuschaos/go-runner 
61dd62c9f863: Mounted from litmuschaos/go-runner 
085ac66cb5b2: Mounted from litmuschaos/go-runner 
620a8c683711: Mounted from litmuschaos/go-runner 
cd5a67aaf636: Mounted from litmuschaos/go-runner 
8310c90d696a: Mounted from litmuschaos/go-runner 
1ed9400dbe50: Mounted from litmuschaos/go-runner 
6f0201089563: Mounted from litmuschaos/go-runner 
1de62c4d61ca: Mounted from litmuschaos/go-runner 
b2d5eeeaba3a: Layer already exists 
1.13.5: digest: sha256:6a6808306d22d37123221ba5abe24a73677f3dc185faa87f504a5afbf978bdc4 size: 3058


The push refers to repository [docker.io/uditgaurav/k8s]
7a09dff94066: Mounted from litmuschaos/k8s 
6502d376cefe: Mounted from litmuschaos/k8s 
9a74f7bf71fd: Mounted from litmuschaos/k8s 
02ca90b54e8c: Mounted from litmuschaos/k8s 
e0c00b80a90e: Mounted from litmuschaos/k8s 
85c50ec826d8: Mounted from litmuschaos/k8s 
597318d985b0: Mounted from litmuschaos/k8s 
69cae67cc70c: Mounted from litmuschaos/k8s 
ace0eda3e3be: Mounted from litmuschaos/k8s 
latest: digest: sha256:768491682ca99e14498ddb62973a95a664a0525ca18e7ad0eaab2621d3124f5b size: 2216

The push refers to repository [docker.io/uditgaurav/litmus-checker]
20d04d96b019: Mounted from litmuschaos/litmus-checker 
c91a7a28cc2e: Mounted from litmuschaos/litmus-checker 
b2d5eeeaba3a: Layer already exists 
latest: digest: sha256:9b619e174d1cb3e4c8202c843f51e4b01d8ea5ebef5ac72205fa42f487bf1698 size: 948

The push refers to repository [docker.io/uditgaurav/workflow-controller]
a199faf3f870: Layer already exists 
ec0a2776976b: Layer already exists 
v2.9.3: digest: sha256:017f7e89e6faca3c0897fa68391fe03b9e83ae475764ded7ca0deef94cdc2676 size: 739

The push refers to repository [docker.io/uditgaurav/argoexec]
ed60a3b92d3c: Layer already exists 
75515e229396: Layer already exists 
9f67644d8e91: Layer already exists 
5906eba2497c: Layer already exists 
b2fe9b4eb3ed: Layer already exists 
71d40540ea58: Layer already exists 
22d1b6b8a5c2: Layer already exists 
f2b050ecb00c: Layer already exists 
c2adabaecedb: Layer already exists 
v2.9.3: digest: sha256:507d8f99c02134c785aea9ca682ea19f746c2d2d9eacbd6cda2036fe582ebc20 size: 2209

The push refers to repository [docker.io/uditgaurav/mongo]
04e5e1d7c68a: Layer already exists 
c09bcffe6c1d: Layer already exists 
a041da64c34b: Layer already exists 
4d1c1210233d: Layer already exists 
d640316cede8: Layer already exists 
cca0bdd11f6b: Layer already exists 
23ebc01e0bc4: Layer already exists 
dae8c5ccdf7e: Layer already exists 
dfb25e85dc08: Layer already exists 
001e4a80973b: Layer already exists 
2ba5b91ca2b0: Layer already exists 
2f37d1102187: Layer already exists 
79bde4d54386: Layer already exists 
4.2.8: digest: sha256:14468b12f721906390c118a38c33caf218c089b751b2f205b2567f99716ae1e9 size: 3032

```
