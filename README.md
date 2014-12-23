TaobaoMobileImageResizer
========================

[![Build Status](https://secure.travis-ci.org/missdeer/TaobaoMobileImageResizer.png)](https://travis-ci.org/missdeer/TaobaoMobileImageResizer)

[![Gobuild Download](http://gobuild.io/badge/github.com/missdeer/TaobaoMobileImageResizer/downloads.svg)](http://gobuild.io/github.com/missdeer/TaobaoMobileImageResizer)

Usage
----
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
