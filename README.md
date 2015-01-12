TaobaoMobileImageResizer
========================

[![Build Status](https://secure.travis-ci.org/missdeer/TaobaoMobileImageResizer.png)](https://travis-ci.org/missdeer/TaobaoMobileImageResizer)

=======

使用方法
----
目前只提供命令行使用方法。

####常规用法：命令行后加入文件名或目录名，程序将所有文件，目录及子目录中的文件都生成对应的小于620*960分辨率的图片。
	./TaobaoMobileImageResizer filename1 [filename2 filename3 ...]
	./TaobaoMobileImageResizer directory1 [directory2 directory3 ...]

####实时监控用法：命令行第1个参数为`-w`或`--watch`，程序将监控后面参数指定的所有目录及子目录，如果有新的文件生成或有文件被修改，则立即生成该文件对应的小于620*960分辨率的图片。这方法适用于修图过程。同时，此方法会在监控前先把现有的文件都扫描并生成对应的图片，相当于常规用法的功能。
	./TaobaoMobileImageResizer -w directory1 [directory2 directory3 ...]

####切图用法：命令行第1个参数为`-c`或`--cut`，程序将第2个及以后参数指定的所有文件，目录及子目录中的文件对比大小，如图图片按比例缩小为620像素宽，高大于960像素的情况，将自动把图片按620*960再切分，切分得到的最后一个图片的高度将可能小于960像素。
	./TaobaoMobileImageResizer -c filename1 [filename2 filename3 ...]
	./TaobaoMobileImageResizer -c directory1 [directory2 directory3 ...]

TODO
----
* 通用而优雅的命令行参数解析
* 支持iOS App的launch image生成
* 支持Android App的launch image生成
* 也许支持iOS/Android App的app icon生成

预编译包下载
----

下载不同平台对应的可执行文件（推荐）：

[Darwin x86](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-darwin-386)

[Darwin amd64](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-darwin-amd64)

[DragonflyBSD x86](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-dragonfly-386)

[DragonflyBSD amd64](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-dragonfly-amd64)

[OpenBSD x86](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-openbsd-386)

[OpenBSD amd64](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-openbsd-amd64)

[NetBSD x86](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-netbsd-386)

[NetBSD amd64](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-netbsd-amd64)

[FreeBSD x86](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-freebsd-386)

[FreeBSD amd64](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-freebsd-amd64)

[FreeBSD arm6](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-freebsd-arm)

[Linux x86](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-linux-386)

[Linux amd64](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-linux-amd64)

[Linux arm6](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-linux-arm)

[Windows x86](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-windows-386), 下载后将文件名添加后缀.exe

[Windows amd64](https://github.com/missdeer/TaobaoMobileImageResizer/raw/prebuilt/TaobaoMobileImageResizer-windows-amd64),下载后将文件名添加后缀.exe 

从Gobuild下载（不推荐）：

[![Gobuild Download](http://gobuild.io/badge/github.com/missdeer/TaobaoMobileImageResizer/downloads.svg)](http://gobuild.io/github.com/missdeer/TaobaoMobileImageResizer)

