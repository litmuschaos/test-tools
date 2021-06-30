#!/bin/bash
# This script is used to pull litmus images required to run generic experiments
# using litmus portal and push into your image registry
set -e

setup(){

declare -ga portal_images=("litmusportal-frontend" "litmusportal-server" "litmusportal-event-tracker"
                       "litmusportal-auth-server" "litmusportal-subscriber")
declare -ga backend_images=("chaos-operator" "chaos-runner" "chaos-exporter" "go-runner")

declare -ga workflow_images=("k8s:latest" "litmus-checker:latest" "workflow-controller:v2.11.0" "argoexec:v2.11.0" "mongo:4.2.8")


if [[ -z "${LITMUS_BACKEND_TAG}" ]];then
  LITMUS_BACKEND_TAG=$(get_latest_backend_release)
fi

if [[ -z "${LITMUS_PORTAL_TAG}" ]];then
  LITMUS_PORTAL_TAG=$(get_latest_portal_release)
fi

if [[ -z "${LITMUS_IMAGE_REGISTRY}" ]];then
  LITMUS_IMAGE_REGISTRY="docker.io"
fi

if [[ -z "${TARGET_IMAGE_REGISTRY}" ]];then
  TARGET_IMAGE_REGISTRY="docker.io"
fi

}

list_all(){

setup
i=1
echo
echo "portal component images ..."
for val in ${portal_images[@]}; do
  echo "${i}. ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}:${LITMUS_PORTAL_TAG}"
  i=$((i+1))
done
echo

echo "backend component images ..."
for val in ${backend_images[@]}; do
  echo "${i}. ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}:${LITMUS_BACKEND_TAG}"
  i=$((i+1))
done
echo

echo "workflow controller images ..."
for val in ${workflow_images[@]}; do
  echo "${i}. ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}"
  i=$((i+1))
done
echo

}

pull_all(){

setup

for val in ${portal_images[@]}; do
  echo " Pulling ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}:${LITMUS_PORTAL_TAG}"
  docker pull ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}:${LITMUS_PORTAL_TAG}
  echo
done
echo

for val in ${backend_images[@]}; do
  echo " Pulling ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}:${LITMUS_BACKEND_TAG}"
  docker pull ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}:${LITMUS_BACKEND_TAG}
  echo
done
echo

for val in ${workflow_images[@]}; do
  echo " Pulling ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}"
  docker pull ${LITMUS_IMAGE_REGISTRY}/litmuschaos/${val}
  echo
done
echo

}

tag_and_push_all(){

setup

if [[ -z "${TARGET_REPONAME}" ]];then
  printf "Please provide the target repo-name by setting TARGET_REPONAME ENV. like:
  TARGET_REPONAME=\"litmuschaos\"\n"
  exit 1
fi

echo
for val in ${portal_images[@]}; do
  IMAGEID=$( docker images -q litmuschaos/${val}:${LITMUS_PORTAL_TAG} )
  docker tag ${IMAGEID} ${TARGET_IMAGE_REGISTRY}/${TARGET_REPONAME}/${val}:${LITMUS_PORTAL_TAG}
  docker push ${TARGET_IMAGE_REGISTRY}/${TARGET_REPONAME}/${val}:${LITMUS_PORTAL_TAG}
  echo
done

for val in ${backend_images[@]}; do
  IMAGEID=$( docker images -q litmuschaos/${val}:${LITMUS_BACKEND_TAG} )
  docker tag ${IMAGEID} ${TARGET_IMAGE_REGISTRY}/${TARGET_REPONAME}/${val}:${LITMUS_BACKEND_TAG}
  docker push ${TARGET_IMAGE_REGISTRY}/${TARGET_REPONAME}/${val}:${LITMUS_BACKEND_TAG}
  echo
done
echo

for val in ${workflow_images[@]}; do
  IMAGEID=$( docker images -q litmuschaos/${val} )
  docker tag ${IMAGEID} ${TARGET_IMAGE_REGISTRY}/${TARGET_REPONAME}/${val}
  docker push ${TARGET_IMAGE_REGISTRY}/${TARGET_REPONAME}/${val}
  echo
done
echo

}

get_latest_backend_release() {
  curl --silent "https://api.github.com/repos/litmuschaos/litmus-go/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

get_latest_portal_release() {
  curl --silent "https://api.github.com/repos/litmuschaos/litmus/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}


print_help(){
cat <<EOF

Usage:       ${0} ARGS (list|pull|push)

list:        "${0} list" will list all the images used by the litmus portal.     


pull:        "${0} pull" will pull the litmus images with the given image tag. 
              The value of tag can be provided by exporting following ENVs:
              - LITMUS_PORTAL_TAG: Tag for the portal component like 'litmusportal-frontend' etc
              - LITMUS_BACKEND_TAG: Tag for backend component chaos-operator, chaos-runner, go-runner etc
              - LITMUS_IMAGE_REGISTRY: Name of litmuschaos image registry. Default is docker.io
              The default images tags are the latest tags released.

push:         "${0} push" will push the images to the given target image registry with the give repo-name
              The value of target images can be set by exporting following ENVs:
              - TARGET_IMAGE_REGISTRY: Give the target registry name. Default is set to "docker.io"
              - TARGET_REPONAME: Give the target image repo-name. This is mandatory to provide.               

EOF

}


case ${1} in
  list)
    list_all
    ;;
  pull)
    pull_all 
    ;;
  push)
    tag_and_push_all
    ;;
  *)
    print_help
    exit 1
esac
