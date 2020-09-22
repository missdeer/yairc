YAIRC - Yet Another Image Resizer & Cutter, formerly known as TaobaoMobileImageResizer
========================

编译
----

需要支持CGO，将编译器路径添加到`PATH`或`CC`、`CXX`路径中。

`go get github.com/missdeer/yairc`

### 功能

- 生成所有尺寸的mobile app应用内icon
- 生成所有尺寸的iOS launch image和app icon

### 使用方法

目前只提供命令行使用方法。

#### 生成通用mobile app应用内icon

```bash
./yairc --action=icons --input=input.png --output=output/directory/path
```

#### 生成iOS launch images：准备一个足够大小的背景图片模板background.png，因为最大的iOS设备是iPad Pro 12"，将使用2048 * 2732大小的launch image，再准备一个足够大的前景图片模板foreground.png，建议至少512 * 512。程序会自动按比例缩放和剪裁图片。

```bash
./yairc --action=launchImage --os=ios -b background.png -f foreground 
```

#### 生成iOS app icons：准备一个1024*1024大小的图片模板template.png

```bash
./yairc --action=appIcon --os=ios --input=template.png
```

#### 生成icns文件

```bash
./yairc --action=icns --input=input.png
```

TODO
----
- [] 支持Android App的splash image生成
- [] 支持Android App的app icon生成
