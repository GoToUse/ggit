#!/usr/bin/env bash
function registerGgit() {
  OS=$1
  lowerOS=$(echo "${OS}" | tr '[:upper:]' '[:lower:]')
  current_pwd=$(pwd ".")
  ggit_path="${current_pwd}/main"
  CGO_ENABLED=0 GOOS=${lowerOS} GOARCH=amd64 go build main.go && ln -sf ${ggit_path} /usr/local/bin/ggit
  if [[ ! -f "${ggit_path}" ]]; then
    echo "${ggit_path} does not exist."
    exit 1
  fi
  if [[ ! -f "/usr/local/bin/ggit" ]]; then
    echo "/usr/local/bin/ggit does not exist."
    exit 1
  fi
  echo "Build&register $()ggit$() successfully, done! You can use $()ggit$() command."
}

function getOS() {
  os=$(uname -s)
  if [[ ${os} -eq "Darwin" ]]
  then
    echo "I'm a macos system."
    registerGgit ${os}
  elif [[ ${os} -eq "Linux" ]] 
  then
    echo "I'm a linux system."
    registerGgit ${os}
  fi
}

# registerGgit
getOS