# 一、每周任务

## 01.24-01.29

module_install配置检查优化【施连生】

openim收发消息接口封装及出文档【施连生】

## 01.17-01.21
openim测试支撑【施连生】
module_install优化配置检查逻辑【施连生】

## 01.04-01.07
完成openim测试程序并测试【施连生】

## 12.27-12.31
module_install问题处理及优化【施连生】
完成openim测试程序并测试【施连生】

## 12.20-12.25

完成module_install重写及测试【施连生】
dgraph调研【施连生】

## 2021.12.06-12.11

完成采购网爬虫（需对接weilink）并部署【施连生】

中国政府采购网关键词命中情况：

# 二、任务详情

## Linux下编译rtcstreamer

测试机器：`137.175.19.178  root/rP76GHfMH9gn`

系统：`ubuntu`

官方教程：[libmediasoupclient install tutorial](https://mediasoup.org/documentation/v3/libmediasoupclient/installation/) 

rtcstreamer里可以把license那块去掉，这样就不依赖xdev的liccheck xwcrypto之类的库了

### 一、编译准备

#### 初始化

```bash
vi ~/.bash_profile
# 添加如下内容
DEPOT_TOOL_PATH=/www/depot_tools
export PATH=$DEPOT_TOOL_PATH:$PATH
export NO_AUTH_BOTO_CONFIG=$DEPOT_TOOL_PATH/proxy.boto
export GYP_DEFINES="build_with_libjingle=1 build_with_chromium=0 libjingle_objc=1 OS=linux"
export GYP_GENERATORS="ninja"
export GYP_CROSSCOMPILE=1

# 之后source
source ~/.bash_profile
```

#### 安装gcc

```bash
sudo apt update
sudo apt install build-essential -y
# 结果
# gcc --version
# gcc (Ubuntu 7.5.0-3ubuntu1~18.04) 7.5.0
```

#### 安装pkg-config

```bash
apt install pkg-config
```

```bash
# 否则会出现问题：
ERROR at //build/config/linux/pkg_config.gni:103:17: Script returned non-zero exit code.
    pkgresult = exec_script(pkg_config_script, args, "value")
```

#### 安装cmake

```bash
mkdir /home/tmp/ -p
cd /home/tmp/
wget https://cmake.org/files/v3.10/cmake-3.10.3-Linux-x86_64.tar.gz
tar -zxf cmake-3.10.3-Linux-x86_64.tar.gz -C /usr/local/
vim ~/.bash_profile

## 在末尾添加如下
export PATH=/usr/local/cmake-3.10.3-Linux-x86_64/bin:$PATH

# quit
source ~/.bash_profile
# 测试
cmake --version
```

### 二、编译webrtc

> 获得libwebrtc.a

[webrtc编译linux](https://note.youdao.com/ynoteshare/index.html?id=958677a46821faf6e23b084c536a7746&type=note&_time=1644399231737)

```bash
cd /www/webrtc-checkout/src
gn gen out/m94 --args='target_os="linux" target_cpu="x64" is_debug=false is_component_build=false is_clang=false rtc_include_tests=false rtc_use_h264=true use_rtti=true use_custom_libcxx=false treat_warnings_as_errors=false use_ozone=true _GLIBCXX_USE_CXX11_ABI=0'
# export CPLUS_INCLUDE_PATH=/usr/include/
ninja -C out/m94

```

#### 直接下载libwebrtc.a【主要使用】

[Release Release m94.4606.4.0 · shiguredo-webrtc-build/webrtc-build (github.com)](https://github.com/shiguredo-webrtc-build/webrtc-build/releases/tag/m94.4606.4.0)

```
wget https://github.com/shiguredo-webrtc-build/webrtc-build/releases/download/m94.4606.4.0/webrtc.ubuntu-18.04_x86_64.tar.gz
tar -zxf webrtc.ubuntu-18.04_x86_64.tar.gz
mv webrtc webrtc.u
mkdir -p gortc-rtcstreamer/thirdparty/lib
cp webrtc.u/lib/libwebrtc.a gortc-rtcstreamer/thirdparty/lib
```

### 三、编译libmediasoupclient

> 获得libmediasoupclient.a
>
> 使用clang++进行编译

```bash
cd /www/libmediasoupclient-3.3.0
### 准备
# 编辑CMakeLists.txt
vim /www/libmediasoupclient-3.3.0/CMakeLists.txt
# 在第二行添加
SET (CMAKE_C_COMPILER             "/usr/bin/clang")
SET (CMAKE_CXX_COMPILER             "/usr/bin/clang++")
SET (CMAKE_CXX_FLAGS                "-stdlib=libc++ -std=c++11")
# quit
### 开始编译
cd /www/libmediasoupclient-3.3.0
cmake . -Bbuild \
  -DLIBWEBRTC_INCLUDE_PATH:PATH=/home/webrtc.u/include/ \
  -DLIBWEBRTC_BINARY_PATH:PATH=/www/webrtc-checkout/src/out/m94/obj
make -C build/

cp /www/libmediasoupclient-3.3.0/build/libmediasoupclient.a /home/gortc-rtcstreamer/thirdparty/lib/
# rm -f /home/gortc-rtcstreamer.nolicence/thirdparty/lib/libmediasoupclient.a; cp /www/libmediasoupclient-3.3.0/build/libmediasoupclient.a /home/gortc-rtcstreamer.nolicence/thirdparty/lib/
# make install -C build/
```

### 四、编译gortc-rtcstreamer

#### 记录

源码：https://go.pfgit.cn/xmedia/gortc-rtcstreamer.git

````
https://go.pfgit.cn/xmedia/gortc-rtcstreamer.git
````

bash命令记录

```
cd ..;rm -rf gortc-rtcstreamer;cp -r gortc-rtcstreamer.bak gortc-rtcstreamer;cd gortc-rtcstreamer;
vim src/CMakeLists.txt
```

#### 准备

```bash
cd /home/
# 下载gortc-rtcstreamer
cd gortc-rtcstreamer/
cp -r gortc-rtcstreamer gortc-rtcstreamer.nolicence
apt install -y dos2unix
dos2unix module_install.sh
mkdir -p /home/gortc-rtcstreamer.nolicence/thirdparty/lib
# cp /home/webrtc.u/lib/libwebrtc.a /home/gortc-rtcstreamer.nolicence/thirdparty/lib # 已经拷贝就不用执行
# vim src/core/XWMain.h
# vim src/core/XWMain.cpp

# vim src/CMakeLists.txt # 进行下面的设置
# module_build.sh设置
# cmake  -DCMAKE_TOOLCHAIN_FILE=/home/linux.toolchain.cmake  .
```

#### 设置CMakeLists.txt

```bash
SET(CMAKE_CXX_COMPILER "/usr/bin/clang++")
#SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -O2 -g0 -std=gnu++14 -stdlib=libstdc++")
SET (CMAKE_CXX_FLAGS                "${CMAKE_CXX_FLAGES} -stdlib=libstdc++ -std=c++11 -stdlib=libstdc++ -pthread")
SET (CMAKE_CXX_FLAGS                "${CMAKE_CXX_FLAGES} -stdlib=libstdc++ -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions")
SET (CMAKE_CXX_COMPILER             "/usr/bin/clang++")
SET (CMAKE_CXX_FLAGS                "-WALL")
SET (CMAKE_CXX_FLAGS_DEBUG          "-g")
SET (CMAKE_CXX_FLAGS_MINSIZEREL     "-Os -DNDEBUG")
SET (CMAKE_CXX_FLAGS_RELEASE        "-s -O2")
SET (CMAKE_CXX_FLAGS_RELWITHDEBINFO "-O2 -g")


## 依赖库设置
set(link_libs 
    xwshare_static
    mysqlclient
    mediasoupclient
    sdptransform
    webrtc
)
```

backup

```bash
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions -lcrypto -D_GLIBCXX_USE_CXX11_ABI=0")
```



#### 编译

```bash
cd /home/;rm -rf gortc-rtcstreamer;cp -r gortc-rtcstreamer.nolicence gortc-rtcstreamer;cd gortc-rtcstreamer;
sh module_build.sh
# sh module_build.sh &> result.log

cd /home/;rm -rf gortc-rtcstreamer;cp -r gortc-rtcstreamer.gnu gortc-rtcstreamer;cd gortc-rtcstreamer;
```





### problem

#### libclientmysql.so问题-已解决

```bash
mkdir -p /home/tmp/
cd /home/tmp/
wget http://launchpadlibrarian.net/212189159/libmysqlclient18_5.6.25-0ubuntu1_amd64.deb
dpkg -i libmysqlclient18_5.6.25-0ubuntu1_amd64.deb
ln -s /usr/lib/x86_64-linux-gnu/libmysqlclient.so.18 /usr/lib/x86_64-linux-gnu/libmysqlclient.so
# 注意：不要使用xdev/lib/下的libmysqlclient.so.18
# 如果还有libssl.so.10 libcrypto.so.10的问题，进行下面的操作（未验证）
# cd /usr/lib/x86_64-linux-gnu
# ln -s /www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libcrypto.so libcrypto.so
# ln -s /www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libssl.so libssl.so
# ln -s /www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libcrypto.so libcrypto.so.10
# ln -s /www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libssl.so libssl.so.10
```

#### undefined reference to symbol 'dlclose@@GLIBC_2.2.5'

出现的问题

```
/usr/bin/ld: /home/gortc-rtcstreamer/src/../thirdparty/lib/libwebrtc.a(latebindingsymboltable_linux.o): undefined reference to symbol 'dlclose@@GLIBC_2.2.5'
//lib/x86_64-linux-gnu/libdl.so.2: error adding symbols: DSO missing from command line
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

`CMakeLists.txt`中添加下面语句

```
set (CMAKE_CXX_LINK_EXECUTABLE "${CMAKE_CXX_LINK_EXECUTABLE} -ldl ") 
```

重新编译后仍出现下面问题

```
/usr/bin/ld: /opt/xdev/lib/libxwshare_static.a(Properties.cpp.o): undefined reference to symbol '_ZNSs5eraseEN9__gnu_cxx17__normal_iteratorIPcSsEES2_@@GLIBCXX_3.4'
//usr/lib/x86_64-linux-gnu/libstdc++.so.6: error adding symbols: DSO missing from command line
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

修改`CMakeLists.txt`中添加的语句

```
set (CMAKE_CXX_LINK_EXECUTABLE "${CMAKE_CXX_LINK_EXECUTABLE} -ldl -lstdc++ ") 
```

重新编译后出现下面问题

```
/usr/bin/ld: /opt/xdev/lib/libxwshare_static.a(Properties.cpp.o): undefined reference to symbol '_ZNSs5eraseEN9__gnu_cxx17__normal_iteratorIPcSsEES2_@@GLIBCXX_3.4'
//usr/lib/x86_64-linux-gnu/libstdc++.so.6: error adding symbols: DSO missing from command line
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

修改`CMakeLists.txt`中添加的语句

```
set (CMAKE_CXX_LINK_EXECUTABLE "${CMAKE_CXX_LINK_EXECUTABLE} -lstdc++ -ldl") 
```

重新编译后出现下面问题

```
/usr/bin/ld: /opt/xdev/lib/libxwshare_static.a(Properties.cpp.o): undefined reference to symbol '_ZNSs5eraseEN9__gnu_cxx17__normal_iteratorIPcSsEES2_@@GLIBCXX_3.4'
//usr/lib/x86_64-linux-gnu/libstdc++.so.6: error adding symbols: DSO missing from command line
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

 下面修改依赖库顺序

原依赖库顺序

```
set(link_libs
    xwshare_static
    mysqlclient
    mediasoupclient
    sdptransform
    webrtc
)
```

新依赖库顺序

```
set(link_libs
    xwshare_static
    mysqlclient
    sdptransform
    mediasoupclient
    webrtc
)
```





#### webrtc编译问题

问题1

```bash
### 问题
cc1plus: warning: unrecognized command line option ‘-Wno-class-memaccess’
cc1plus: warning: unrecognized command line option ‘-Wno-packed-not-aligned’
### 使用的命令
cd /www/webrtc-checkout/src
gn gen out/m94 --args='target_os="linux" target_cpu="x64" is_debug=false is_component_build=false is_clang=false rtc_include_tests=false rtc_use_h264=true use_rtti=true use_custom_libcxx=false use_ozone=true rtc_enable_protobuf=false rtc_use_x11=false'
ninja -C out/m94 
### 解决办法：升级gcc，8.x以上
sudo apt-get update 
sudo apt-get install gcc-8
sudo apt-get install g++-8
sudo update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-8 100
sudo update-alternatives --config gcc
sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-8 100
sudo update-alternatives --config g++
```

#### XWMain-MAIN_PROBLEM

```
CMakeFiles/rtcstreamer.dir/core/XWMain.cpp.o: In function `XWMain::Init(int, char**)':
/home/gortc-rtcstreamer/src/core/XWMain.cpp:33: undefined reference to `mediasoupclient::Device::Load(nlohmann::basic_json<std::map, std::vector, std::string, bool, long, unsigned long, double, std::allocator, nlohmann::adl_serializer>, mediasoupclient::PeerConnection::Options const*)'

##############################
nm -C out/m94/lib/libwebrtc.a |grep __cxx11|head
##############################
需要的是std::，但是提供的是std::__cxx11::
```



#### cannot find libssl.so.10 libcrypto.so.10

```bash
####
/opt/xdev/lib/libmysqlclient.so: undefined reference to `SSL_get_verify_result@libssl.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `OPENSSL_add_all_algorithms_noconf@libcrypto.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `SSL_CTX_ctrl@libssl.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `SSL_accept@libssl.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `SSL_new@libssl.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `ERR_get_error_line_data@libcrypto.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `SSL_CTX_set_cipher_list@libssl.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `SSL_CIPHER_get_name@libssl.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `ERR_get_error@libcrypto.so.10'
/opt/xdev/lib/libmysqlclient.so: undefined reference to `DH_free@libcrypto.so.10'
#############
/www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libssl.a
/www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libssl.so
/www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libcrypto.so

ln -s /www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libcrypto.so libcrypto.so
#############
cd /usr/lib/x86_64-linux-gnu
ln -s libssl.so.1.0.0          libssl.so.10
ln -s libcrypto.so.1.0.0    libcrypto.so.10

rm -f /usr/lib/x86_64-linux-gnu/libssl.so.10
rm -f /usr/lib/x86_64-linux-gnu/libcrypto.so.10

ln -s libcrypto.so libcrypto.so.10
ln -s libssl.so libssl.so.10
ln -s libcrypto.so.10 libcrypto.so.1.0.0
```



### backup

```bash
# cp /opt/xdev/xwshare/SysConfigReader.* /opt/xdev/xwshare/BaseThread.* /opt/xdev/xwshare/XWLog.* /opt/xdev/xwshare/ModuleConfigHelper.* /opt/xdev/xwshare/Properties.* /home/gortc-rtcstreamer/src/
# cp -r /opt/xdev/xwshare/* /home/gortc-rtcstreamer/src/
# vim src/core/XWMain.h
# vim src/core/XWMain.cpp

# cp /www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libcrypto.so /home/gortc-rtcstreamer.bak/thirdparty/lib/
# cp /www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/lib/x86_64-linux-gnu/libcrypto.a /home/gortc-rtcstreamer.bak/thirdparty/lib/
# cp /usr/local/lib/libmediasoupclient.a /home/gortc-rtcstreamer.bak/thirdparty/lib/
# cp /home/webrtc/lib/libwebrtc.a /home/gortc-rtcstreamer.bak/thirdparty/lib/

# vim src/CMakeLists.txt
# SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions -lcrypto -D_GLIBCXX_USE_CXX11_ABI=0")


# 2.14.09
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions -lcrypto -D_GLIBCXX_USE_CXX11_ABI=0 -stdlib=stdlib++")
```



## libmediasoupclient

> 在centos上进行编译

[libmediasoupclient install tutorial](https://mediasoup.org/documentation/v3/libmediasoupclient/installation/) 

```
https://mediasoup.org/documentation/v3/libmediasoupclient/installation/
https://github.com/versatica/libmediasoupclient/archive/refs/tags/3.3.0.zip
```

```shell
gn gen out/m94 --args='is_debug=false is_component_build=false is_clang=true rtc_include_tests=false rtc_use_h264=true use_rtti=true use_custom_libcxx=true treat_warnings_as_errors=false use_ozone=true rtc_enable_protobuf=false rtc_use_x11=false'
export CPLUS_INCLUDE_PATH=/www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/include/x86_64-linux-gnu/

gn gen out/m94 --args='is_debug=false is_component_build=false is_clang=false rtc_include_tests=false rtc_use_h264=true use_rtti=true use_custom_libcxx=true treat_warnings_as_errors=false use_ozone=true rtc_enable_protobuf=false rtc_use_x11=false'
export CPLUS_INCLUDE_PATH=/www/webrtc-checkout/src/build/linux/debian_sid_amd64-sysroot/usr/include/x86_64-linux-gnu/
ninja -C out/m94

gn gen out/m94 --args='target_os="linux" target_cpu="x64" is_debug=false is_component_build=false is_clang=false rtc_include_tests=false rtc_use_h264=true use_rtti=true use_custom_libcxx=false treat_warnings_as_errors=false use_ozone=true rtc_enable_protobuf=false rtc_use_x11=false'
export CPLUS_INCLUDE_PATH=/usr/include/
ninja -C out/m94

nm -C out/m94/lib/libwebrtc.a |grep __cxx11|head
```

## 编译gortc-rtcstreamer-ubuntu

### 基础信息

测试机器：`137.175.19.178  root/rP76GHfMH9gn`

系统：`ubuntu`

官方教程：[libmediasoupclient install tutorial](https://mediasoup.org/documentation/v3/libmediasoupclient/installation/) 

### 编译环境

[webrtc编译linux](https://note.youdao.com/ynoteshare/index.html?id=958677a46821faf6e23b084c536a7746&type=note&_time=1644399231737)

#### 初始化

```bash
vi ~/.bash_profile
# 添加如下内容
DEPOT_TOOL_PATH=/www/depot_tools
export PATH=$DEPOT_TOOL_PATH:$PATH
export NO_AUTH_BOTO_CONFIG=$DEPOT_TOOL_PATH/proxy.boto
export GYP_DEFINES="build_with_libjingle=1 build_with_chromium=0 libjingle_objc=1 OS=linux"
export GYP_GENERATORS="ninja"
export GYP_CROSSCOMPILE=1

# 之后source
source ~/.bash_profile
```

#### 安装gcc

默认安装gcc-7.5.0

```shell
sudo apt update
sudo apt install build-essential -y
```

#### 安装cmake

文件地址：[https://cmake.org/files/](https://cmake.org/files/)

```
https://cmake.org/files/
```

下载编译

方式一：（不推荐使用）

```bash
mkdir /home/tmp/ -p
cd /home/tmp/
wget https://cmake.org/files/v3.10/cmake-3.10.3.tar.gz
tar -zxf cmake-3.10.3.tar.gz
cd cmake-3.10.3
./bootstrap
# maybe need to update g++
# sudo apt-get install g++
# other problem
# sudo apt-get install libssl-dev
make -j4
make install
cmake --version
```

方式二：

```bash
mkdir /home/tmp/ -p
cd /home/tmp/
wget https://cmake.org/files/v3.10/cmake-3.10.3-Linux-x86_64.tar.gz
tar -zxf cmake-3.10.3-Linux-x86_64.tar.gz -C /usr/local/
vi ~/.bash_profile
# 在末尾添加如下
export PATH=/usr/local/cmake-3.10.3-Linux-x86_64/bin:$PATH
# quit
source ~/.bash_profile
# 测试
cmake --version
```

### 编译libmediasoupclient

```bash
cd /www/webrtc-checkout/src
gn gen out/m94 --args='target_os="linux" target_cpu="x64" is_debug=false is_component_build=false is_clang=false rtc_include_tests=false rtc_use_h264=true use_rtti=true use_custom_libcxx=false treat_warnings_as_errors=false use_ozone=true rtc_enable_protobuf=false rtc_use_x11=false'
# export CPLUS_INCLUDE_PATH=/usr/include/
ninja -C out/m94
```





## 编译gortc-rtcstreamer

### 基础信息记录

源码：https://go.pfgit.cn/xmedia/gortc-rtcstreamer.git

````
https://go.pfgit.cn/xmedia/gortc-rtcstreamer.git
````

g++参数

```
g++ -MMD -MF obj/net/dcsctp/tx/rr_send_queue/rr_send_queue.o.d -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -DABSL_ALLOCATOR_NOTHROW=1 -I../.. -Igen -I../../third_party/abseil-cpp -fno-ident -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -D__DATE__= -D__TIME__= -D__TIMESTAMP__= -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions --sysroot=../../build/linux/debian_sid_amd64-sysroot -fvisibility-inlines-hidden -Wnon-virtual-dtor -Woverloaded-virtual -c ../../net/dcsctp/tx/rr_send_queue.cc -o obj/net/dcsctp/tx/rr_send_queue/rr_send_queue.o


-D__DATE__= -D__TIME__= -D__TIMESTAMP__= -Wall
-DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -DABSL_ALLOCATOR_NOTHROW=1

SET(CMAKE_CXX_COMPILER /opt/rh/devtoolset-8/root/usr/bin/g++)
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions -lcrypto -D_GLIBCXX_USE_CXX11_ABI=0")
```

### 环境

测试机器：`137.175.19.178  root/viKNP86aAblA`

### 初始化

```
cd ..;rm -rf gortc-rtcstreamer;cp -r gortc-rtcstreamer.bak gortc-rtcstreamer;cd gortc-rtcstreamer;vim src/CMakeLists.txt

rm -f /home/gortc-rtcstreamer/thirdparty/lib/libwebrtc.a;cp /home/libwebrtc.a.old /home/gortc-rtcstreamer/thirdparty/lib/libwebrtc.a
```



CMakeList.txt

```
SET(CMAKE_CXX_COMPILER /opt/rh/devtoolset-8/root/usr/bin/g++)
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions")

SET(CMAKE_CXX_COMPILER /opt/rh/devtoolset-8/root/usr/bin/g++)
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions -lcrypto -D_GLIBCXX_USE_CXX11_ABI=1")

SET(CMAKE_CXX_COMPILER /opt/rh/devtoolset-9/root/usr/bin/g++)
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions -Wl,-Bdynamic")

-stdlib=libc++
```

### 编译

#### 1. 将所有的文件转换成unix下格式



#### 2. 安装cmake

cmake最低版本`3.8`

下载cmake包

```
wget https://cmake.org/files/v3.15/cmake-3.15.7.tar.gz
```

解压cmake包并进入目录

```
tar -zxf cmake-3.15.7.tar.gz
cd cmake-3.15.7
```

编译安装cmake

```
./bootstrap
gmake
gmake install
```

---

如果之前并没有安装过cmake，下面的部分就不需要继续了

查看编译后的cmake版本

```text
 /usr/local/bin/cmake --version
```

移除原来的cmake版本，新建软连接

```text
 yum remove cmake -y
 ln -s /usr/local/bin/cmake /usr/bin/
```

终端查看版本

```text
cmake --version
```

#### 3. 编译

替换`src/CMakeList.txt`内容：

```shell
# 将下面的
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp")
# 替换为
SET(CMAKE_CXX_COMPILER /opt/rh/devtoolset-9/root/usr/bin/g++)
SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fopenmp -MMD -MF -DUSE_UDEV -DUSE_AURA=1 -DUSE_GLIB=1 -DUSE_NSS_CERTS=1 -DUSE_OZONE=1 -DUSE_X11=1 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D_FORTIFY_SOURCE=2 -DCR_SYSROOT_HASH=43a87bbebccad99325fdcf34166295b121ee15c7 -DNDEBUG -DNVALGRIND -DDYNAMIC_ANNOTATIONS_ENABLED=0 -DWEBRTC_ENABLE_PROTOBUF=1 -DWEBRTC_INCLUDE_INTERNAL_AUDIO_DEVICE -DRTC_ENABLE_VP9 -DWEBRTC_HAVE_SCTP -DWEBRTC_USE_H264 -DWEBRTC_LIBRARY_IMPL -DWEBRTC_NON_STATIC_TRACE_EVENT_HANDLERS=0 -DWEBRTC_POSIX -DWEBRTC_LINUX -fno-strict-aliasing --param=ssp-buffer-size=4 -fstack-protector -funwind-tables -fPIC -pipe -pthread -m64 -march=x86-64 -msse3 -Wno-builtin-macro-redefined -Wall -Wno-unused-local-typedefs -Wno-maybe-uninitialized -Wno-deprecated-declarations -Wno-comments -Wno-packed-not-aligned -Wno-missing-field-initializers -Wno-unused-parameter -O2 -fdata-sections -ffunction-sections -fno-omit-frame-pointer -g0 -fvisibility=hidden -std=gnu++14 -Wno-narrowing -Wno-class-memaccess -fno-exceptions")
```



### 遇到问题

#### 编译cmake

错误信息

```
Error when bootstrapping CMake:
Cannot find a C++ compiler that supports both C++11 and the specified C++ flags.
Please specify one using environment variable CXX.
The C++ flags are "".
They can be changed using the environment variable CXXFLAGS.
See cmake_bootstrap.log for compilers attempted.
```

解决：

将gcc版本降下来，gcc一般安装在`/usr/bin/`下，将备份的gcc(一般都是gcc-v4.8.5)还原，然后继续执行

编译完cmake之后记得将gcc版本还原



### Wrong Message

```

```



## elasticsearch修改用户名和密码

[elasticsearch启动常见问题 - 妖怪梧桐 - 博客园 (cnblogs.com)](https://www.cnblogs.com/sitongyan/p/11263753.html)

### 文档

中文文档链接：[《Elasticsearch中文文档》 | Elasticsearch 技术论坛 (learnku.com)](https://learnku.com/docs/elasticsearch73/7.3)

官方文档链接：[Elasticsearch Guide [7.16]](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)



### 命令调用

#### 使用docker部署的情况

参考链接：[如何给ElasticSearch设置用户名和密码 - 知乎 (zhihu.com)](https://zhuanlan.zhihu.com/p/163337278?ivk_sa=1024320u)

添加用户



添加密码：

```
./bin/elasticsearch-setup-passwords interactive
```

使用curl访问es

```bash
curl localhost:9200 -u elastic:{password}
## 如
curl localhost:9200 -u elastic:{123456}
```

修改密码

```bash
curl -XPOST -u elastic "localhost:9200/_security/user/elastic/_password" -H 'Content-Type: application/json' -d'{"password" : "abcd1234"}'
```

## module_install优化

### 调整检查安装配置所在安装目录格式

格式：

```
安装包名：{project}-{module}-install-{date}.tar.gz
解压后目录名：{project}-{module}-install
```

其中将包名或目录名通过`-`进行分段，包为4段，目录为3段，且第三段必须是`install`

### [EAGLE-95](http://jira.xw:9090/browse/EAGLE-95)

> 模块安装的时候，检查globe 配置情况

1. 检查所有的IP配置，以是不是以`_IP`后缀结尾的
2. 如果配置值是`localhost`或者`127.0.0.1`，则直接报`WARN`
3. 如果配置值不是2中的值，则检查ip是不是有效的ip，如果不是，则报`WARN`

### [EAGLE-89](http://jira.xw:9090/browse/EAGLE-89)

> 安装模块优化全局配置检查逻辑，将private配置添加到globe

1. 将private中的配置取出
2. 与globe中的配置进行对比
3. 如果globe中不存在配置key，添加到globe中
4. 如果globe中存在配置key，但是value为空，添加到globe中
5. 检查更新的globe中的private项，是否存在空，存在，则报错，退出程序

### module_install.sh杀安装程序

问题描述：

​	使用安装包进行安装  
​	解压后，到解压目录下运行里面的module_install.sh，出现将自己的进程杀掉的情况  

探究问题：【以modelprocess】安装程序为例

1.研究该moudle_install.sh将要杀掉的进程是什么进程

```
命令：echo $(ps aux | grep modelprocess)
使用go run main.go install **.tar.gz 进行安装 里面含有module_install.sh
root 43418 5.8 12.0 2626348 225780 	? 		Sl 	13:45 0:21 java -Xms58m -Xmx58m -jar -Dspring.config.location=/usr/local/modelprocess/conf/application.yml /usr/local/modelprocess/bin/modelprocess.jar 
root 47176 6.5 0.4 	901352 	7672 	pts/0 	Sl+ 13:51 0:00 go run main.go install titan-modelprocess-install-1.tar.gz 
root 47214 0.3 0.0 	703916 	824 	pts/0 	Sl+ 13:51 0:00 /tmp/go-build1641831010/b001/exe/main install titan-modelprocess-install-1.tar.gz 
root 47221 0.0 0.0 	113284 	1224 	pts/0 	S+ 	13:51 0:00 bash -c cd /tmp/module_install/titan-modelprocess-install;sh module_install.sh install 
root 47255 0.0 0.0 	112824 	964 	pts/0 	S+ 	13:51 0:00 grep modelprocess
想要杀掉的进程列表：43418 47176 47214 47221
失败进程：47214
是否杀完：否

使用module_install install **.tar.gz进行安装
root 51628 17.8 12.2 	2626348 228532 	? 		Sl 13:58 0:20 java -Xms58m -Xmx58m -jar -Dspring.config.location=/usr/local/modelprocess/conf/application.yml /usr/local/modelprocess/bin/modelprocess.jar 
root 52626 0.0 	0.0 	113284 	1408 	pts/0 	S+ 13:59 0:00 /bin/bash /opt/xdev/bin/module_install install titan-modelprocess-install-1.tar.gz 
root 52630 0.5 	0.1 	127472 	3476 	pts/0 	S+ 13:59 0:00 /opt/xdev/bin/install_tool/module_install install titan-modelprocess-install-1.tar.gz 
root 52640 0.0 	0.0 	113284 	1228 	pts/0 	S+ 13:59 0:00 sh -c cd /usr/local/etc/module_install/titan-modelprocess-install;sh module_install.sh install 
root 52667 0.0 	0.0 	112824 	964 	pts/0 	S+ 13:59 0:00 grep modelprocess
进程列表：51628 52626 52630 52640
失败进程：52630
是否杀完：是
```

原因：都是想要将运行的程序杀死，但是python比较厉害，自己挂了还能继续跑脚本  
解决办法：在main函数中接收kill信号，然后忽略，当执行kill当前程序时，并不会kill当前程序，该程序将会继续执行下去


## openim测试程序
第二版：单机收发
代码地址：https://go.pfgit.cn/attachments/c317fa20-169b-4e65-b981-3e302447a253
使用说明：https://go.pfgit.cn/shiliansheng/openim/src/master/%e4%bd%bf%e7%94%a8%e8%af%b4%e6%98%8e.md

## openIM环境搭建
```
git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive
#git clone https://github.com/OpenIMSDK/Open-IM-SDK-Core.git
docker-compose up -d
./docker_check_service.sh 
./check_all.sh
```

## dgraph调研

### docker
监听端口：
	gRPC server started.  Listening on port 9080
	HTTP server started.  Listening on port 8080
	CONN: Connecting to zero:5080
	Server Listening on port 20000

## 完成module_install重写及测试

module_install.sh的在 http://samba.xw/xianwei/release/bigdata/titan/
安装模块代码：放在tools/module_install下 https://172.16.0.58:7777/svn/xw/xdev

出现一个问题，就是如果使用go执行脚本命令，脚本中途出错，退出，则go程序也直接退出

## edumeet

github链接：https://github.com/edumeet/edumeet

## 爬虫
http://172.16.5.11:28100/task/{task_id}
查询链接
```
原始查询链接
url := 	"http://search.ccgp.gov.cn/bxsearch?searchtype=2&page_index=1&dbselect=bidx&" + 
			"kw=" + keywords + "&start_time=" +	begintime.Format(timeFormat) + 
			"&end_time=" + endtime.Format(timeFormat) + "&timeType=2"
指定查询时间： 将timeType修改为6
url := 	"http://search.ccgp.gov.cn/bxsearch?searchtype=2&page_index=1&dbselect=bidx&" + 
			"kw=" + keywords + "&start_time=" +	begintime.Format(timeFormat) + 
			"&end_time=" + endtime.Format(timeFormat) + "&timeType=6"
```
## 12.6
合并采购网爬虫到脚手架中
```
# tobe
# 下划线
# 说明
# 都在common中的config中读取
# 查找后使用on duplicate key update ts = now()
# 使用iris查看时，使用一个线程阻塞住
```