#!/usr/bin/env bash

# Currently this should not be changed; stay tuned
INSTALL_DIR="/opt"

function verify_root {
  if [[ $(whoami) != 'root' ]]; then
    echo "You must be root to install this program"
    exit 1
  fi
}

function verify_requirements {
  export CURL_EXEC=$(command -v curl)
  if [[ $CURL_EXEC == '' ]]; then
    echo "Please install curl and retry the installation script"
    exit 1
  fi

  export TAR_EXEC=$(command -v tar)
  if [[ $TAR_EXEC == '' ]]; then
    echo "Please install tar and retry the installation script"
    exit 1
  fi
}

function prepare_install_root {
  mkdir -p "$INSTALL_DIR"
}

function download_release {
  RELEASE_URL=$($CURL_EXEC -s https://api.github.com/repos/Cesura/avologo/releases/latest | grep browser_download_url | cut -d'"' -f4)
  export PACKAGE_NAME=${RELEASE_URL##*/}

  echo "Downloading latest release..."
  $CURL_EXEC -L -s -O "$RELEASE_URL"
}

function extract_release {
  echo "Extracting archive..."
  $TAR_EXEC xf "$PACKAGE_NAME" -C "$INSTALL_DIR"
}

function copy_config {
  echo "Copying sample config file..."
  cp -n "$INSTALL_DIR/avologo/avologo.conf.example" /etc/avologo.conf
}

function prepare_systemd {
  if [[ $(command -v systemctl) != "" ]]; then
    echo "Preparing systemd..."
    cp $INSTALL_DIR/avologo/systemd/* /etc/systemd/system/
    systemctl daemon-reload
  fi
}

function success_message {
  echo "Done!"
  SHORT_PKG=${PACKAGE_NAME%.tar.gz}
  echo "You have successfully installed $SHORT_PKG"
  echo "Please configure avologo in /etc/avologo.conf and then run:"
  echo -e "\tsystemctl start avologo-server"
  echo "or"
  echo -e "\tsystemctl start avologo-client"
}

# Entry point
verify_root
verify_requirements
download_release
prepare_install_root
extract_release
copy_config
prepare_systemd
success_message
