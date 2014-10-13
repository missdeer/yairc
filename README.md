TaobaoMobileImageResizer
========================
[![Gobuild Download](http://gobuild.io/badge/github.com/missdeer/TaobaoMobileImageResizer/downloads.svg)](http://gobuild.io/github.com/missdeer/TaobaoMobileImageResizer)

Usage
----
常规用法：命令行后加入文件名或目录名，程序将所有文件，目录及子目录中的文件都生成对应的小于620*960分辨率的图片。

./TaobaoMobileImageResizer filename1 [filename2 filename3 ...] 

./TaobaoMobileImageResizer directory1 [directory2 directory3 ...]

实时监控用法：命令行第1个参数为-w或--watch，程序将监控后面参数指定的所有目录及子目录，如果有新的文件生成或有文件被修改，则立即生成该文件对应的小于620*960分辨率的图片。这方法适用于修图过程。同时，此方法会在监控前先把现有的文件都扫描并生成对应的图片，相当于常规用法的功能。

./TaobaoMobileImageResizer -w directory1 [directory2 directory3 ...]
