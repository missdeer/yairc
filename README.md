YAIRC - Yet Another Image Resizer & Cutter, formerly known as TaobaoMobileImageResizer
========================

[![Build Status](https://secure.travis-ci.org/dfordsoft/yairc.png)](https://travis-ci.org/dfordsoft/yairc)

=======

编译
----

需要支持CGO，将编译器路径添加到`PATH`或`CC`、`CXX`路径中。

`go get github.com/dfordsoft/yairc`

功能
----
- 生成针对淘宝移动端适配的图片
- 生成所有尺寸的iOS launch image和app icon

使用方法
----
目前只提供命令行使用方法。

####生成iOS launch images：准备一个足够大小的背景图片模板background.png，因为最大的iOS设备是iPad Pro 12"，将使用2048 * 2732大小的launch image，再准备一个足够大的前景图片模板foreground.png，建议至少512 * 512。程序会自动按比例缩放和剪裁图片。

```bash
./yairc -b background.png -f foreground --launchImage
```

####生成iOS app icons：准备一个1024*1024大小的图片模板template.png

```bash
./yairc --appIcon template.png
```

#### 生成iOS app所需要的3倍图片文件：准备好图片文件template.png

```bash
./yairc --iOSScale --as=2x template.png
```

#### 生成HiDPI设备多分辨率图片文件：准备好图片文件template.png

```bash
./yairc --iconScale -i template.png -o destine/path
```

####生成淘宝移动版适配图片常规用法：命令行后加入文件名或目录名，程序将所有文件，目录及子目录中的文件都生成对应的小于620*960分辨率的图片。

```bash
./yairc filename1 [filename2 filename3 ...]
./yairc directory1 [directory2 directory3 ...]
```

####生成淘宝移动版适配图片实时监控用法：命令行第1个参数为`-w`或`--watch`，程序将监控后面参数指定的所有目录及子目录，如果有新的文件生成或有文件被修改，则立即生成该文件对应的小于620*960分辨率的图片。这方法适用于修图过程。同时，此方法会在监控前先把现有的文件都扫描并生成对应的图片，相当于常规用法的功能。

```bash
./yairc -w directory1 [directory2 directory3 ...]
```

####生成淘宝移动版适配图片切图用法：命令行第1个参数为`-c`或`--cut`，程序将第2个及以后参数指定的所有文件，目录及子目录中的文件对比大小，如图图片按比例缩小为620像素宽，高大于960像素的情况，将自动把图片按620*960再切分，切分得到的最后一个图片的高度将可能小于960像素。

```bash
./yairc -c filename1 [filename2 filename3 ...]
./yairc -c directory1 [directory2 directory3 ...]
```

TODO
----
* 支持Android App的splash image生成
* 支持Android App的app icon生成
