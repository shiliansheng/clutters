WORKDIR=$(dirname $0)

module="ccgpgov"

a_info=$(ps -ef|grep bin/${module}|grep -v grep)
export LD_LIBRARY_PATH=./lib:$LD_LIBRARY_PATH
if [ "x$a_info" = "x" ];then
	echo "start ${module}"
	cd $WORKDIR
	./bin/${module} >/dev/null 2>/dev/null
	cd -
else
	echo "${module} is running"
	exit 0
fi

