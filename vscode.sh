#!/usr/bin/env bash

SHELL_FOLDER=$(cd $(dirname ${BASH_SOURCE[0]});pwd)

fileUrl="https://vscode.download.prss.microsoft.com/dbazure/download/stable/7c6fdfb0b8f2f675eb0b47f3d95eeca78962565b/code-stable-x64-1743783505.tar.gz"
fileName=`echo $fileUrl|awk -F "/" '{print $NF}'`

_main() {
    which sudo >/dev/null && SUDO="sudo"

    cd /tmp \
    && wget -c $fileUrl -O $fileName \
    && chmod 666 ${fileName} \
    && ${SUDO} tar --no-same-owner -zxvf $fileName -C /opt/

    if [ $? -ne 0 ];then
        echo "install faild"
        exit 1
    fi
    ${SUDO} cp ${SHELL_FOLDER}/icon/vscode.png /opt/VSCode-linux-x64/vscode.png \
    && ${SUDO} cp ${SHELL_FOLDER}/desktop/code.desktop /usr/share/applications/code.desktop \
    && ${SUDO} cp ${SHELL_FOLDER}/desktop/code-url-handler.desktop /usr/share/applications/code-url-handler.desktop \
    && ${SUDO} ln -fs /opt/VSCode-linux-x64/code /usr/local/bin/code \
    && ${SUDO} ln -fs /opt/VSCode-linux-x64/code-tunnel /usr/local/bin/code-tunnel \
    && echo "install vscode latest success"
}


_main
