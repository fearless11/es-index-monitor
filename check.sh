#!/bin/bash

# desc: for project start|stop|pack
# * * * * * /usr/local/monitor/log/check.sh start &> /dev/null

DIR="/usr/local/monitor/log"
APP="es-index-monitor"
CFG="cfg.yaml"

chmod +x ${DIR}/${APP} &> /dev/null

function stop(){
  ps aux |grep ${APP} | egrep -v "${DIR}|grep" | awk '{print $2}' | xargs kill
}

function start(){
  num=`ps aux | grep ${APP} | egrep -v "${DIR}|grep" | wc -l`
  if [[ $num < 1 ]];then
    cd ${DIR}
    cp -rf ./log ./log.old
    ./${APP} -c ./${CFG}  &>  ./log &
    echo "start ok"
    exit
  fi
  echo "${APP} already running"
}

function pack() {
  export GOOS=linux
  go build .
  tar zcvf es-index-monitor-`date +%F`.tgz es-index-monitor cfg.yaml check.sh 
  rm -rf es-index-monitor
}

if [ $# != 1 ];then
  echo "usage: ./check.sh start|stop|pack"
  exit 
fi

case $1 in
  start)
     start;;
  stop)
     stop;;
  pack)
     pack;;
esac
