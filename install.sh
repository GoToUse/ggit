#!/usr/bin/env bash
function registerGgit() {
    current_pwd=$(pwd ".")
    ggit_path="${current_pwd}/main"
    go build main.go && ln -sf ${ggit_path} /usr/local/bin/ggit
    if [[ ! -f "${ggit_path}" ]]; then
      echo "${ggit_path} does not exist."
      exit 1
    fi
    if [[ ! -f "/usr/local/bin/ggit" ]]; then
      echo "/usr/local/bin/ggit does not exist."
      exit 1
    fi
    echo "Build&register ``ggit`` successfully, done! You can use ``ggit`` command."
}

registerGgit