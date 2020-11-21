# wechat_dat_to_image
微信本地图片格式转化

PC端微信收到的图片存储在本地，我们可以打开微信的设置->通用设置->文件管理->打开文件夹，在FileStorage/Image日期目录下找到图片文件，但是它经过了简单的加密并且后缀为.dat。

通过下面这个算法可以将dat文件还原成原始图片，目前支持jpg，gif，png三种格式，其他格式应该也同理。

先直接上代码，golang代码如下：
```
const (
	jpg0 = 0xFF
	jpg1 = 0xD8
	gif0 = 0x47
	gif1 = 0x49
	png0 = 0x89
	png1 = 0x50
)

func Dat2Image(datpath string) (string, error) {
	b, err := ioutil.ReadFile(datpath)
	if err != nil {
		return "", err
	}
	if len(b) < 2 {
		return "", errors.New("image size error")
	}

	j0 := b[0] ^ jpg0
	j1 := b[1] ^ jpg1
	g0 := b[0] ^ gif0
	g1 := b[1] ^ gif1
	p0 := b[0] ^ png0
	p1 := b[1] ^ png1
	var v byte
	var ext string
	if j0 == j1 {
		v = j0
		ext = "jpg"
	} else if g0 == g1 {
		v = g0
		ext = "gif"
	} else if p0 == p1 {
		v = p0
		ext = "png"
	} else {
		return "", errors.New("unknown image format")
	}

	for i := range b {
		b[i] = b[i] ^ v
	}

	imgpath := datpath[0:len(datpath)-len(ext)] + ext
	err = ioutil.WriteFile(imgpath, b, os.ModePerm)
	return imgpath, err
}
```
成功后会在dat所在的目录生成相同名字的图片文件（后缀是具体图片格式）

解释如下：
1. 前六个常量根据名字可以知道分别是jpg，gif，png图片十六进制文件的前两个字符，你可以用ultraedit等软件打开相应格式的图片，通过十六进制方式查看前两个字符就是上面的常量，这个数值是固定了

2. 这个加密只是简单的用同一个字符对图片的每个字符进行异或，我们只要找到是哪个值就可以了，比如这个变量的名字是v

3. 这里只列举了三种常用的格式，我们假定图片一定是这三种中的一种，这样就很好办了，因为这三种图片的前两个字符是已知的，而dat文件的前两个字符也是已知的，当我们将对应的字符异或后如果相等就能确定图片确实是这种格式并且异或后的值就是我们要找的那个变量v（这里要好好理解下），如果不相等就继续换下一种格式尝试，总共只有那么几种格式很容易就试出来了

4. 最后用得到的v对dat文件的每个字符进行异或就可以还原出原始图片

下面用golang写了一个简单的转换程序，只要将wechat_dat_to_image.exe在需要转换的dat文件所在目录下运行就可以了，在github源码里可以下载到这个程序。

[github源码](https://github.com/tujiaw/wechat_dat_to)  
[https://github.com/tujiaw/wechat_dat_to](https://github.com/tujiaw/wechat_dat_to)