package file

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

// 文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)

	return len(content), err
}

// 文件扩展名
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// 检测是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// 检测是否有权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// 是否存在目录，不存在新建
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// 创建目录
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 必须打开最大限度试图打开文件
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

// 下载文件
func DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// 文件名
	fileName := Md5Str(url) + ".jpg"
	dir, _ := os.Getwd()
	// 保存路径
	src := filepath.Join(dir, "runtime", "download")
	err = IsNotExistMkDir(src)
	if err != nil {
		return "", err
	}
	// 文件路径
	filePath := filepath.Join(src, fileName)
	err = os.WriteFile(filePath, body, 0777)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// 创建缩略图
func CreateThumb(imagePath string, width, height int) (string, error) {
	src, err := imaging.Open(imagePath)
	if err != nil {
		fmt.Println("failed to open image:", err)
		return "", err
	}
	// 重新画尺寸
	src = imaging.Resize(src, width, height, imaging.Lanczos)
	dir, _ := os.Getwd()
	// 保存路径
	thumbPath := filepath.Join(dir, "runtime", "thumb")
	err = IsNotExistMkDir(thumbPath)
	if err != nil {
		return "", err
	}
	// 文件名
	fileName := Md5Str(imagePath+strconv.Itoa(width)+"_"+strconv.Itoa(height)) + ".jpg"
	filePath := filepath.Join(thumbPath, fileName)
	err = imaging.Save(src, filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// 检测是否http资源
func CheckIsHttpResource(name string) bool {
	if strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "http://") {
		return true
	}
	return false
}

// 删除文件
func RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

// md5加密
func Md5Str(val string) string {
	m := md5.New()
	m.Write([]byte(val))
	return hex.EncodeToString(m.Sum(nil))
}
