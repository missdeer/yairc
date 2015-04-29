YAIRC - Yet Another Image Resizer & Cutter, formerly known as TaobaoMobileImageResizer
========================

[![Build Status](https://secure.travis-ci.org/missdeer/yairc.png)](https://travis-ci.org/missdeer/yairc)

=======

功能
----
- 生成针对淘宝移动端适配的图片
- 生成所有尺寸的iOS launch image和app icon（WIP）

使用方法
----
目前只提供命令行使用方法。

####常规用法：命令行后加入文件名或目录名，程序将所有文件，目录及子目录中的文件都生成对应的小于620*960分辨率的图片。
	./yairc filename1 [filename2 filename3 ...]
	./yairc directory1 [directory2 directory3 ...]

####实时监控用法：命令行第1个参数为`-w`或`--watch`，程序将监控后面参数指定的所有目录及子目录，如果有新的文件生成或有文件被修改，则立即生成该文件对应的小于620*960分辨率的图片。这方法适用于修图过程。同时，此方法会在监控前先把现有的文件都扫描并生成对应的图片，相当于常规用法的功能。
	./yairc -w directory1 [directory2 directory3 ...]

####切图用法：命令行第1个参数为`-c`或`--cut`，程序将第2个及以后参数指定的所有文件，目录及子目录中的文件对比大小，如图图片按比例缩小为620像素宽，高大于960像素的情况，将自动把图片按620*960再切分，切分得到的最后一个图片的高度将可能小于960像素。
	./yairc -c filename1 [filename2 filename3 ...]
	./yairc -c directory1 [directory2 directory3 ...]

TODO
----
* 支持Android App的splash image生成
* 也许支持Android App的app icon生成

预编译包下载
----

下载不同平台对应的可执行文件（推荐）：

[Darwin x86](https://github.com/missdeer/yairc/raw/prebuilt/yairc-darwin-386)

[Darwin amd64](https://github.com/missdeer/yairc/raw/prebuilt/yairc-darwin-amd64)

[DragonflyBSD x86](https://github.com/missdeer/yairc/raw/prebuilt/yairc-dragonfly-386)

[DragonflyBSD amd64](https://github.com/missdeer/yairc/raw/prebuilt/yairc-dragonfly-amd64)

[OpenBSD x86](https://github.com/missdeer/yairc/raw/prebuilt/yairc-openbsd-386)

[OpenBSD amd64](https://github.com/missdeer/yairc/raw/prebuilt/yairc-openbsd-amd64)

[NetBSD x86](https://github.com/missdeer/yairc/raw/prebuilt/yairc-netbsd-386)

[NetBSD amd64](https://github.com/missdeer/yairc/raw/prebuilt/yairc-netbsd-amd64)

[FreeBSD x86](https://github.com/missdeer/yairc/raw/prebuilt/yairc-freebsd-386)

[FreeBSD amd64](https://github.com/missdeer/yairc/raw/prebuilt/yairc-freebsd-amd64)

[FreeBSD arm6](https://github.com/missdeer/yairc/raw/prebuilt/yairc-freebsd-arm)

[Linux x86](https://github.com/missdeer/yairc/raw/prebuilt/yairc-linux-386)

[Linux amd64](https://github.com/missdeer/yairc/raw/prebuilt/yairc-linux-amd64)

[Linux arm6](https://github.com/missdeer/yairc/raw/prebuilt/yairc-linux-arm)

[Windows x86](https://github.com/missdeer/yairc/raw/prebuilt/yairc-windows-386), 下载后将文件名添加后缀.exe

[Windows amd64](https://github.com/missdeer/yairc/raw/prebuilt/yairc-windows-amd64),下载后将文件名添加后缀.exe 

