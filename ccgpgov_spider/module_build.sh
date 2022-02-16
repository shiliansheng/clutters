#!/bin/bash
# modify 2019.11.25
set -evu

cur_date=$(date +%Y%m%d)

workdir=$(cd $(dirname $0) && pwd)
cd ${workdir}

# todo 设置模块名
module=ccgpgov
project=spider
modulename=${project}-${module}

packagedir=${modulename}-install
rm -rf ${packagedir}*
mkdir -p ${packagedir}

# todo 打包过程，如：编译，下载包，拷贝文件等。
cd cmd
go build -o ${module} .
if [ $? != 0 ];then
   echo "ERROR:build failed"
   exit 1
fi

cd ${workdir}
mkdir -p ${packagedir}/bin
cp -f cmd/${module} ${packagedir}/bin/
cp -f build/* ${packagedir}/
cp -f ChangeLog ${packagedir}/
cp -rf configs ${packagedir}/
cp -rf views ${packagedir}/

# 最终打包，一般不用修改
cd ${workdir}
#svn info > ${packagedir}/ModuleInfo
find ${packagedir} -name '*.svn' | xargs -i rm -rf {}
find ${packagedir} -name '*.sh' -o -name '*.py' -o -name '*.bin' | xargs -i chmod -v +x {}
find ${packagedir} -name '*.guard' -o -name '*.xagent' -o -name '*.sh' -o -name '*.py' -o -name '*.conf' -o -name '*.toml' | xargs -i dos2unix {}
tar zcvf ${packagedir}-${cur_date}.tar.gz ${packagedir}
echo "$modulename build ok."
