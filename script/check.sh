#!/bin/bash

DIR="/usr/local/monitor/es-index-monitor"
APP="es-index-monitor"
CFG="cfg.yaml"

chmod +x ${DIR}/${APP}

function stop(){
  ps aux |grep ${APP} | egrep -v "${DIR}|grep" | awk '{print $2}' | xargs kill
}

function start(){
  num=`ps aux | grep ${APP} | egrep -v "${DIR}|grep" | wc -l`
  if [[ $num < 1 ]];then
    cd ${DIR}
    ./${APP} -c ./${CFG}  &>  ./log &
    echo "start ok"
    exit
  fi
  echo "${APP} already running"
}


case $1 in
  start)
     start;;
  stop)
     stop;;
esac
