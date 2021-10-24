export PATH=$PATH:$HOME

APP_PATH=`pwd`/unripe-bison-app

usage() {
  echo "Usage: sh app.sh [start|stop|restart|status]"
  exit 1
}

is_exist(){
  pid=`ps -ef|grep $APP_PATH|grep -v grep|awk '{print $2}'`
  if [ -z "${pid}" ]; then
    return 1
  else
    return 0
  fi
}

start(){
  is_exist
  if [ $? -eq 0 ]; then
    echo "app is already running. pid=${pid}"
  else
    nohup ${APP_PATH} >/dev/null 2>&1  &
  fi
}

stop(){
  is_exist
  if [ $? -eq "0" ]; then
    kill -9 $pid
  else
    echo "app is not running"
  fi  
}

status(){
  is_exist
  if [ $? -eq "0" ]; then
    echo "App is running. PID is ${pid}"
  else
    echo "App is NOT running."
  fi
}

restart(){
  stop
  sleep 1
  start $1
}

case "$1" in
  "start")
    start $2
    ;;
  "stop")
    stop
    ;;
  "status")
    status
    ;;
  "restart")
    restart $2
    ;;
  *)
    usage
    ;;
esac
