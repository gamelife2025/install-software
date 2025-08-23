#!/bin/bash

VERSION="v2.0.3"

if [ -n "$1" ];then
    VERSION=$1
fi


FILEURL=https://github.com/syncthing/syncthing/releases/download/${VERSION}/syncthing-linux-amd64-${VERSION}.tar.gz
FILENAME=syncthing-linux-amd64-${VERSION}.tar.gz


_main() {
    which sudo >/dev/null && SUDO="sudo"

    cd /tmp \
    && wget -c $FILEURL \
    && chmod 666 ${FILENAME} \
    && tar zxvf  ${FILENAME} \
    && cd ${FILENAME/.tar.gz} \
    && ${SUDO} cp syncthing /usr/local/bin/syncthing \
    && echo "install syncthing  ${VERSION} success"
}

_main