// Code generated by go-bindata.
// sources:
// ../../static/index.html
// ../../static/js/app.js
// ../../static/partials/api.html
// ../../static/partials/applications.html
// ../../static/partials/navbar.html
// ../../static/partials/nodes.html
// DO NOT EDIT!

package static

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _indexHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x54\x5f\x53\x9b\x4e\x14\x7d\xf7\x53\xf0\xdb\x57\x7f\x40\x10\x63\x62\x07\x9c\x89\x86\xd6\x66\x12\x13\x8d\x25\xea\xdb\x06\x2e\xb0\x64\xd9\x85\xdd\x25\x21\x7e\xfa\x2e\x49\xad\xf6\x9f\xb5\x9d\xe9\x13\x7b\xef\x9e\xfb\xe7\x1c\xe0\x78\xff\x0d\xa7\x17\xb7\xf7\xb3\xc0\xc8\x54\x41\xcf\x0e\xbc\xf6\x61\x50\xcc\x52\x1f\x01\x43\x06\x4b\x4d\x5c\x96\x3e\xa2\x5c\x60\x09\x62\x0d\x02\x9d\x1d\x18\x86\x97\x01\x8e\xdb\x83\x3e\x16\xa0\xb0\x11\x65\x58\x48\x50\x3e\xaa\x55\x62\xf6\xd1\xcb\xab\x4c\xa9\xd2\x84\xaa\x26\x6b\x1f\xdd\x99\x9f\x06\xe6\x05\x2f\x4a\xac\xc8\x92\x02\x32\x22\xce\x14\x30\x5d\xf7\x31\xf0\x21\x4e\xe1\x9b\x4a\x86\x0b\xf0\xd1\x9a\xc0\xa6\xe4\x42\xbd\x00\x6f\x48\xac\x32\x3f\x86\x35\x89\xc0\xdc\x05\xff\x1b\x84\x11\x45\x30\x35\x65\x84\x29\xf8\xce\x53\x23\x45\x14\x85\xb3\x31\xbf\xc1\xc6\x7c\xb7\xbf\x67\xef\x53\xfb\x6b\x4a\xd8\xca\x10\x40\x7d\x24\xd5\x96\x82\xcc\x00\xf4\xa0\x4c\x40\xe2\xa3\x76\x71\xf9\xce\xb6\x0b\xdc\x44\x31\xb3\x96\x9c\x2b\xa9\x04\x2e\xdb\x20\xe2\x85\xfd\x35\x61\xbb\x96\x6b\x9d\xd8\x91\x94\xcf\x39\xab\x20\x1a\x25\x25\xd2\x8b\x29\x48\x05\x51\x5b\x3d\x23\xc3\x6e\xff\xd8\x74\xaa\x7e\x71\x3b\x9a\x0e\xe6\x4d\x3f\x77\x06\xf5\x21\xee\x2e\x86\x21\x9b\x91\x23\xba\x7a\x9f\x6c\x36\xc1\x00\xf7\xb3\xe1\x30\xce\x1f\x68\x39\x86\xb4\xc9\xf2\x70\x12\x38\x49\x9a\x2f\x66\x1f\x8a\xd5\xa3\xec\x69\x25\x04\x97\x92\x0b\x92\x12\xe6\x23\xcc\x38\xdb\x16\xbc\x96\xe8\x5f\x93\x32\x55\x06\x05\xbc\x46\x2d\x19\x2f\x8e\xae\x3a\x0e\x9d\x54\x39\x5e\x9d\xaf\x1a\x97\xda\x93\xd3\x00\x67\xf5\xa6\x9c\x27\x70\xb5\x0e\x4f\xdc\x51\x17\x1e\x99\x5b\x3f\x3c\xe2\xf2\xb6\x53\xf7\x82\x7b\x79\x37\xc9\xaf\xc3\xc3\x4e\xc0\xba\xe2\x77\xd4\x64\x24\x48\xa9\x0c\x29\xa2\x67\x2a\x38\xc7\x8d\x95\x72\x9e\x52\xc0\x25\x91\x3b\x1a\x6d\xce\xa6\x64\x29\xed\xbc\xaa\x41\x6c\x6d\xc7\x72\x1c\xcb\xfd\x12\xed\x18\xe4\xba\xa9\x67\xef\x1b\xbe\xd2\xfd\xad\x42\xe5\xdf\xbf\xfc\xfc\xa7\x02\x75\x8a\xf9\x72\x34\x0c\x2e\xf5\xa7\x9a\x14\xf5\xf9\xf9\xf5\xec\x64\x70\x7c\x2d\x4a\x51\x75\xa7\x61\xb2\x70\x7b\xb3\x9b\x1b\x37\xef\x06\xe3\xaa\x91\xd2\xd9\x86\xd5\x54\x31\x28\xd9\x65\x38\x3b\xc5\xa3\x5e\x33\xff\xb5\x40\x6f\xe0\xf2\xba\x52\xfa\xa7\xaf\x29\x16\x9a\x88\x63\x75\xad\xce\x53\xfc\x27\x62\xfd\xd5\x00\x53\xf0\x5a\xc1\x5b\xc6\xe8\x4a\xed\x46\x3f\xa0\x3c\xfb\xc9\x8e\xbc\x25\x8f\xb7\xfb\x42\xc3\x8b\xc9\xba\xf5\xaf\xd6\x3f\x34\x5c\x47\x7b\xec\x1e\xa2\x6b\x5a\xc3\xfb\x1c\x00\x00\xff\xff\x78\x44\x55\x33\x00\x05\x00\x00")

func indexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_indexHtml,
		"index.html",
	)
}

func indexHtml() (*asset, error) {
	bytes, err := indexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.html", size: 1280, mode: os.FileMode(420), modTime: time.Unix(1457866204, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _jsAppJs = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x59\x5f\x6f\xdb\x36\x10\x7f\xf7\xa7\xe0\xbc\xa0\x94\x31\x57\x5e\x1f\x67\x23\x18\xda\x74\x28\x0a\x34\x45\x91\x3f\x4f\x41\x1e\x34\x89\xb6\x35\x48\xa2\x40\xd2\xc9\x82\xc0\xdf\x7d\x77\xd4\x3f\x52\xa2\x6c\x29\x69\xd3\x6e\x98\x5e\x6c\x91\xc7\xd3\xdd\xef\xee\x7e\x77\xb2\xef\x02\x41\x12\x2e\x02\xc9\xc4\x1d\x13\xe4\x94\x04\xd9\x66\x97\x04\xc2\x4f\x79\xb4\x4b\x98\x47\x9b\x4d\x3a\x27\x37\x13\x02\x17\xcd\x36\x17\x7c\xa7\x18\x9d\x4f\x8a\xfb\x46\xe6\x8c\x67\x4a\xf0\x24\x61\x42\x52\xbd\x77\x3b\x5b\x4d\x26\xcd\xbe\x1f\xf2\x6c\x1d\x6f\xbc\x9b\xe9\x49\x2e\xf8\x5d\x1c\xb1\xe9\x9c\xac\x77\x59\xa8\x62\x9e\x79\xd5\xda\x8c\x3c\xea\xc3\x82\xa9\x9d\xc8\x48\xb5\xec\x47\x2c\x04\x4d\x8a\x0b\x8f\x9e\x6c\x95\xca\xd1\x22\x7a\x12\xb1\x84\x6d\x02\x34\xc7\xd0\x54\x2d\x56\xaa\xf0\xaa\x17\x7d\x91\x87\xe0\x6a\x2d\x9d\x32\xb5\xe5\xd1\x9c\xe4\x81\x08\xe0\x3b\x18\x6f\x1e\xc3\x0b\x61\x8a\x02\x15\xc0\xa9\xc7\x69\x21\x3e\x5d\x92\xea\xdc\x54\x1f\x94\xb0\x72\xd3\xa8\xb8\x85\xf5\x38\x9a\x92\x25\x79\xb3\x5f\x59\xca\x2a\xb7\x6a\x73\x72\x2e\x95\x47\x17\x60\x15\xf8\x80\x8f\x99\x93\x47\xba\x65\x41\x84\x30\x2e\x1f\x29\xa2\xca\x32\xf5\xfa\xea\x21\x67\x74\x49\x68\x90\xe7\x49\x1c\x06\x68\xfb\xe2\x2f\xc9\x33\xba\xdf\xcf\x9a\x47\x18\x4f\x6b\x3f\xc9\x10\xba\x2d\x4f\xec\xfb\x22\x44\x4f\x04\x06\xf9\x4b\x81\x3d\x44\x5f\x8b\x37\x08\x5b\xbb\x16\xcc\xd6\x8e\x6f\xb9\x7e\xbf\x65\x19\x78\x6a\x38\x20\xc1\x65\x1b\x6b\xbc\x14\x4b\xf3\x04\x0c\xbe\x16\x09\xf8\x0b\xa0\xaa\x38\x48\xa4\x75\xce\xdf\xaa\x34\x29\xad\x32\xaf\xb0\x4e\x41\x38\xf9\xb6\x39\xf0\x29\x96\xea\x4c\x89\x84\x5a\x27\xf6\xb3\xa3\x06\x2e\x96\xc6\xdd\x0f\x69\x6d\xc6\x23\x36\x0a\x47\x7d\x60\x90\x49\x9f\x41\x72\xbc\x2d\x8b\x25\x7e\xfc\x20\x16\x05\x79\x3c\x2e\x6a\xf1\xb0\x60\x7d\xf9\x78\xdc\x04\xae\xb6\x4c\xdc\xc7\x92\x79\xdd\xe7\x0b\x16\xc5\x82\x85\xea\x8a\x83\x36\xbb\x24\x5a\x3a\x5b\x65\x3b\xb1\x49\xdb\xe0\xdc\x83\xfc\x6d\x72\x33\x10\xa7\xd6\xb4\x58\x90\x28\x96\x00\xc3\x03\x41\x4e\x8c\x78\x28\x27\xce\x13\x7e\xe3\xbb\x57\xbb\xae\xe9\x57\x86\x3c\xc7\x50\x57\x94\x6c\xd1\x04\xf1\x8a\xfd\x39\xd1\xbb\x16\x4d\xe8\x0d\x3f\x0f\x36\x0c\xcc\xa6\x18\xa5\xc6\x4f\x2d\xed\x6f\x58\x45\x8b\x33\x5f\xee\xc2\x90\x49\xe9\xd5\x04\x84\x3c\xd9\xa6\xe9\x52\x27\xe8\x42\x28\x50\xc2\x80\xce\x22\x3c\x70\x3c\x0d\x32\x7c\xb6\x09\xfc\x10\xdf\x1d\x35\xea\xc4\x81\x54\x04\xaa\x7b\x43\x73\x7f\x18\xa1\x39\x31\x4f\x55\x77\x07\x71\x33\xf2\xa6\x0d\x20\x40\x67\x59\xec\x7f\x60\x0a\xad\xc6\x7a\xa0\x49\x9c\xc6\x0a\xda\xc9\x6f\x70\x81\x75\x7c\xbd\x96\x0c\xef\x7f\xdd\x0f\x44\xdb\x42\x3c\xaf\x10\xf7\x05\x93\xbb\x44\xd9\xc0\xb7\xad\x0f\x05\x83\xca\x33\x7b\x30\x68\x68\xeb\x8f\xd7\xb8\x4a\x4e\x4f\x49\xb6\x4b\x12\xe7\xd3\x3d\xfa\x73\xa1\xea\x9c\x47\x41\x02\x69\x92\xe2\xa7\x37\xf3\x41\x23\xdd\xc6\x51\xc4\x32\xff\x4f\x59\xac\x9a\xf3\x81\x4b\x99\x56\xa8\xe1\x06\x17\x12\x1e\x44\xde\x6c\xd5\x11\xda\xb7\xd6\xf6\x84\x25\x92\xb9\x4c\x73\x07\xe0\x4c\x5b\x0b\xa6\xa0\xbf\xc3\x61\x2e\xe0\x20\x7a\xdf\x67\x42\x70\x71\x08\x96\x63\xf0\x20\x32\x8c\x3a\xbc\xd3\x1e\xb9\x81\x29\xe2\x56\x3e\x9a\x34\x76\x0c\x80\xc8\x9c\x4c\xda\x89\xb0\xcb\xa3\xe3\x89\xd0\x83\xe5\xb5\x3e\x3b\x1a\xcb\x71\x38\x22\x86\x40\xd4\x6a\x20\x82\x5d\xf4\x06\x22\xd7\x9a\xe0\xda\x38\xe1\x04\x37\xa0\x60\x88\xa7\xc7\x37\x91\x02\x52\x82\x91\x07\xbe\x23\x72\x57\x7e\xb9\x0f\x32\x45\x14\x27\xa5\x2a\x4a\x7e\x41\xe0\xb0\x74\xff\xb8\xfe\x08\x37\xf4\x77\x3a\x73\x16\x99\x1b\xfd\xf7\x5a\x4d\x81\x7e\xa9\xe4\x79\x09\xfd\xd3\xf1\x84\x0e\x80\x82\x95\x71\x66\x5c\x0a\x8f\xae\x6d\x67\x40\x80\x95\x4c\x8a\xf6\x0d\x02\x1e\x98\xb7\x40\xc2\xd4\xe6\x79\x4b\xc9\x53\xc8\xb7\x8f\x7b\x6b\x49\x67\x16\x0f\xa1\x49\xb2\x85\x7e\xc1\xc5\x83\xbf\xe1\xde\xeb\x37\xb3\x55\x17\xa5\xee\xab\x47\xbb\xc5\xea\xe1\x6e\x40\x6f\xb5\x46\xbb\x97\x6d\xaa\x4d\xa4\xd0\x88\xaf\xdb\x27\xcb\x30\x69\x14\x9e\xd7\x24\x51\x85\xab\xe8\x71\xfd\xdb\xf4\xc9\x17\x6c\x93\x1a\xf8\xba\x3f\x6a\x57\xff\x6f\x90\x8e\x06\xe9\x4a\x82\x36\x8c\x75\x6b\x1c\x09\xe3\x7f\xa4\x37\xf6\xd5\xc9\xa8\xe6\x88\x4a\x40\xf7\xdd\xe0\xee\xa8\x91\xaf\xdb\xa2\x71\x7c\x5c\x1e\xff\x3b\xdb\x62\x19\x0b\x09\x5e\x22\x01\x8f\xcb\xd7\xcb\xe2\x14\xd2\xee\xbb\x87\xf7\x1a\xb3\x27\x23\x58\xb1\x6d\x2f\xd5\x56\x17\xe2\x6c\x08\x1f\x49\x75\x43\x6b\x7f\x20\x0a\x63\x97\xa6\xe5\xdd\x5f\x11\xaa\xab\x98\x99\x4a\xe1\xe2\xa6\x5f\x78\x7d\x96\xa9\xeb\x1c\x9a\xce\x61\x91\xf7\xfc\x3e\x03\x21\x77\xf8\x07\x55\x24\x94\x74\x19\xc3\xef\xf6\x4a\x75\xb8\xda\x0b\x42\xbc\x74\xe4\x59\xe7\x57\xdb\x9e\x2c\x6b\xc8\x51\x0e\xcf\xab\xf1\xcc\xe8\x84\xf1\xab\x90\x23\x19\x83\xd7\x86\xa9\x8b\x20\x8b\x78\x0a\x75\xf5\x36\x8a\xc4\x93\x20\xfb\xd0\x52\x82\xe0\xa1\xfb\xc3\xcb\x32\x93\x58\x10\xa5\x01\x87\x67\xd5\xe7\x35\x85\xf6\x7c\x3e\xa4\x5d\x3a\x26\xf2\x91\x8d\xd3\x98\xf1\x8e\xf9\x07\x4d\x48\xf2\x84\xf9\x09\xdf\x14\xea\xbe\xe7\xb8\xde\xe1\xef\x2d\xbf\xbf\x64\x99\x3a\x0f\xc2\x33\x9e\xc2\xfc\x1e\x49\x33\x61\xaa\xed\x36\x04\xd5\x3a\xc8\xd6\x5f\x5f\xbd\xaa\x1b\x34\xc2\xf1\x29\x90\xaa\x6a\x66\x2b\xd7\x7c\xec\xd0\x61\xd8\x57\x84\x96\xf4\x5a\x50\xab\x49\x2d\xcb\xdb\x16\xf8\x86\x63\x6d\xdc\x4b\x59\xc1\xe0\x5d\x8c\x39\x0f\x5f\x30\x99\x43\xf4\xd8\xa0\xe1\xf6\xa8\x3d\x9f\xd9\xdf\x9d\xf4\xe8\xd8\xd0\x85\xcb\xe8\xbf\x9d\xf0\xe1\xc3\x52\x3b\x64\x98\x92\x73\xd2\x07\x5b\xab\x12\x0c\x74\xb0\x28\xce\x0b\x75\x57\xfc\x1d\xc3\xc3\x4f\xee\xc9\xe3\x46\x9a\xaa\x40\xb4\xf0\xe1\xb1\xa6\xcb\x9a\xed\xea\x2a\xcb\xd0\x55\x64\xc3\x27\x76\xbb\xba\xf1\xa3\x57\xa6\xfc\x97\xd0\x62\x80\x3e\x59\x67\xc6\x1f\x91\x35\x22\xd4\xd4\x82\x9b\x41\x8c\xcc\x7b\x79\x22\x29\xda\x74\x0f\x8d\x14\x39\x69\xd8\xd7\x4e\x05\x63\xab\x1a\xc1\x4f\xcd\xdc\x5b\x0d\x4d\xe2\xba\xd7\x9b\x0f\xfb\x96\xef\x43\xbd\xa0\xbf\xf8\x6b\x51\xf5\x0f\x41\xf1\xef\xa0\xf1\xcb\xcc\x3f\x01\x00\x00\xff\xff\x61\x19\x29\x00\xf4\x1f\x00\x00")

func jsAppJsBytes() ([]byte, error) {
	return bindataRead(
		_jsAppJs,
		"js/app.js",
	)
}

func jsAppJs() (*asset, error) {
	bytes, err := jsAppJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "js/app.js", size: 8180, mode: os.FileMode(420), modTime: time.Unix(1460967341, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _partialsApiHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xbc\x54\xcb\x6e\xdb\x30\x10\xbc\xfb\x2b\x58\xe5\xe0\x04\xa8\x44\xa0\xe9\x29\xa0\x04\x14\x3d\xb5\x87\xd4\x48\xfb\x03\x6b\x69\x23\x12\xa5\x48\x75\x49\x19\x30\x02\xff\x7b\x49\xc9\xb2\xa5\x58\x70\x1f\x87\x5c\x4c\x7a\x39\x3b\x3b\x9c\xa1\x2d\x4c\x9d\x2a\x53\xea\xae\x42\xe6\xa8\xcc\x93\x75\x0b\xe4\x15\x68\xc7\x0d\xec\xb6\x40\x99\xf4\x8d\x5e\x27\x85\xe0\x67\x64\xb1\x5a\x89\x4a\xed\x58\xa9\xc1\xb9\x3c\x29\xad\xf1\xa0\x0c\x52\x52\xac\x18\x13\xf2\x43\xf1\xf5\xfb\xb7\xc7\xf4\x69\xf3\x99\x7d\xda\x7c\x11\x3c\x14\x62\x3d\x76\x04\x0e\xc2\x16\xc1\xe7\xc9\x2d\xb4\xea\x11\x1a\x7c\xcf\xc2\xe6\x8e\x29\x13\x57\xd7\x53\x0c\x24\x2f\x2f\xec\x08\x61\x87\x03\x6b\xd0\x4b\x5b\xb9\x91\x6d\x81\x6f\x40\x0c\x94\xc3\x7e\x64\xcd\x8e\xcd\x47\xf2\x48\xff\xb1\x10\xc0\x24\xe1\x73\x9e\x48\xef\x5b\xf7\xc0\x79\x6d\x2b\x5b\x66\x96\x6a\x5e\x2b\x2f\xbb\x6d\x56\xda\x86\x6f\xc9\x96\x00\xc4\xb5\x25\x70\x48\x3b\xa4\x9b\x99\xae\x70\xc1\x2c\x14\xce\xb3\x43\x2d\x99\x4b\xbf\x38\x17\x1c\x82\x9d\x41\xc2\xa8\x66\xe2\x65\x0b\x06\x35\xeb\x3f\xd3\x0a\x9f\xa1\xd3\xfe\xa4\x7a\x01\x99\x4a\x84\x4a\x99\x7a\x82\x89\xb7\xbb\x9f\x83\xbc\xf2\x1a\x93\xe2\x09\x7f\x75\xe8\xfc\x03\x13\x2e\xd4\x4f\x90\x4e\xeb\x94\x54\x2d\x7d\xaf\x7b\x50\x9a\x01\xd5\x5d\x83\xc6\xff\xd8\xb7\x38\x46\x30\xb4\xc5\x67\x10\x5c\x7b\x85\xdb\xfc\xac\x37\xe0\x25\x7b\x97\xb3\x75\x78\x2d\xb7\x57\xdc\xbd\x1c\x32\x36\x1f\x0e\x37\x57\x15\x2c\x09\x3c\xf7\x46\x5b\xef\x04\x8f\x22\x8b\xd3\x22\xef\x27\xee\xf1\x60\xdf\x35\x33\xb7\xb6\xda\xcf\x9d\x6c\x09\x17\x66\xc6\xe7\xdd\x0f\x8c\xc7\xcb\xf4\xf3\x2f\x6f\x16\xb0\x6b\xad\x71\xf8\x57\x09\x87\xdf\x8d\xde\xff\x29\xde\x1e\xf4\x5f\xd9\xce\x3a\x67\xc1\xbe\x1e\x7c\x21\xea\x2d\x23\xed\x07\xfe\x5b\x9e\xd3\xad\xa4\xfe\xaf\x6d\xa8\x0c\xcb\xef\x00\x00\x00\xff\xff\x13\x14\x53\x6f\x52\x05\x00\x00")

func partialsApiHtmlBytes() ([]byte, error) {
	return bindataRead(
		_partialsApiHtml,
		"partials/api.html",
	)
}

func partialsApiHtml() (*asset, error) {
	bytes, err := partialsApiHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "partials/api.html", size: 1362, mode: os.FileMode(420), modTime: time.Unix(1458470541, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _partialsApplicationsHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x56\x4f\x6f\x9c\x3e\x10\xbd\xef\xa7\xb0\xfc\x93\x7e\x49\x0f\x2c\x6a\x95\x53\x0b\x48\x51\x9a\x43\xa5\xb6\x97\xaa\x1f\x60\xc0\xde\xc5\xaa\xb1\x91\x31\xab\xac\xa2\x7c\xf7\x8e\x6d\x60\x0d\x4b\x92\xe6\xd0\xaa\xaa\x72\x48\xf0\x9f\xc7\x78\xe6\xf9\xcd\x5b\x32\xb5\x4f\x84\xaa\x64\xcf\x38\xe9\x4c\x95\xd3\x8b\x16\x8c\x15\x20\xbb\x54\xc1\xa1\x04\xb3\xad\x6d\x23\x2f\x68\x91\xa5\x27\x64\xb1\xd9\x64\x4c\x1c\x48\x25\xa1\xeb\x72\x5a\x69\x65\x41\x28\x6e\x68\xb1\x21\x24\xab\xdf\x15\xd7\x6d\x2b\x45\x05\x56\x68\xd5\x65\x29\x2e\xe0\x3a\x21\x6e\xd3\x42\x29\xf9\xf8\x62\x98\xf8\xff\x49\xad\x0f\x43\x00\x87\xaa\x39\xb0\x30\xf6\x33\x17\xf0\xf6\xfb\xa7\x2c\xc5\x61\xb4\xfa\x15\x1a\xbe\x58\x9b\x42\xf3\x3b\x9b\x18\xb1\xaf\x2d\x2d\xae\xab\x21\x91\x11\xe9\x46\xd3\x01\x99\x2d\x35\x3b\x9e\x42\x18\x82\x85\x1a\xde\x72\xb0\x39\x85\xb6\x25\x42\x11\x7c\x74\x74\x84\x38\x10\x2b\xee\xef\xdd\xea\x16\x7c\x62\xe4\xe1\x01\x63\xb2\x55\x84\xc2\x24\xd7\xf6\xd7\x32\x9d\xf6\x11\x01\x23\xa0\xb4\x8a\xe0\x5f\xd2\x35\xfe\xc1\xf8\x0e\x7a\x69\x29\xa9\x0d\xdf\xe5\x34\xfd\x2f\x85\x88\xed\x74\x99\x17\x2d\x6e\x99\xb0\x59\x0a\xb3\xe0\x65\x6f\xad\x56\x8f\x9d\x00\x6a\x8f\x77\xe1\x78\xa8\x30\xf0\x8f\x9c\x32\x2e\xb9\xe5\x97\x18\xf4\x0d\x2d\x3e\xfa\xc9\x2c\x62\x5c\x1c\x8e\xcd\x44\xf3\x48\x2d\x0e\xdd\x2d\xa3\x70\x70\x1c\x49\xc7\x1d\xb7\x37\xba\x6f\x29\x31\x5a\xf2\x9c\x0e\x13\x30\x02\x12\x09\x25\x97\x39\xdd\x6e\xb7\xa3\x2e\x86\xb4\xed\xb1\x45\x68\x98\xd0\x65\x11\x13\x3f\xa7\xf4\x2b\x83\x97\xc9\x2f\x31\xf7\x1b\x3f\x22\x11\x63\x59\x1a\xe2\x84\x2c\x31\xb5\x62\x33\x3c\x66\x1a\x6f\x34\x03\x49\x76\xc0\x38\x25\x82\xe5\x94\x23\xa9\x5f\xdc\x1a\x75\xfa\x15\x8a\xf1\xbb\x9c\x26\x6f\xc7\x32\x18\x76\x90\xde\xc7\x75\x48\xce\xca\x23\xc6\x39\xfa\xb7\x3e\xbb\xa5\xd0\x2e\xcb\x43\x92\xf1\xdd\x21\x92\xae\xfa\x86\xab\x51\x1c\xe7\x70\xd7\x7d\xa7\xfd\x35\x84\xd3\xfa\xd4\x5b\xcf\xf2\x58\x49\xdd\x61\x91\x0c\x2c\x60\x2a\x5d\x23\xa6\x40\xf3\x6b\xb9\xf1\xb8\x22\xeb\x5a\x50\x61\xa3\x16\x8c\x71\x85\x8a\x36\x3d\x6e\xfc\x6f\x45\xc3\xbb\x0f\x59\xea\x00\x45\x4c\xf3\x90\x44\x7d\x35\x4f\xd3\x0a\x2b\x07\x76\xe7\x2c\x39\x01\x93\xf3\x7e\xab\xaf\x4e\x92\xf3\x17\xf6\x58\xfd\x4e\x84\x71\xf5\x3b\x6d\x9a\x59\x37\x44\x6f\xb8\xbd\x41\x91\x31\x04\x41\xbe\x6c\x82\xfb\x39\x75\x2d\x4d\x63\xf7\xc3\xab\x0a\xb4\x50\xef\x49\xef\xb3\xd4\xcf\x16\x11\x84\x6a\x7b\x3b\x90\xee\xba\x9e\xce\x0e\x1d\xe2\x04\x02\xc2\x09\xa8\x60\xac\xc0\x91\x3d\x1a\xc9\xdc\x22\xe2\xb2\xdd\x34\x2e\xec\x19\x4e\x76\x5a\xdb\xb9\x26\x22\x48\xdb\x4b\xac\x87\xef\x30\x59\x67\x4e\x91\x1f\x08\x74\x1c\x6e\x8c\xc6\x37\xf1\x42\xfc\xc8\xdf\xc5\x3c\x91\x17\xb5\xe9\x9a\xd0\x0a\x2f\xae\x15\xcd\xfc\x4a\xe4\xd6\x88\x06\xcc\x31\x36\x80\xbe\x65\x30\xf9\xd7\x37\x38\xe0\x2f\x50\xed\x8a\xea\x96\x47\x44\x95\x4c\xc3\x17\xb8\x42\x30\x9a\x57\x5f\xf8\x33\xbe\xb0\x66\xe6\x7f\xa7\x29\x84\xcf\x97\xdf\x63\x0b\xc1\x11\x9f\x32\x86\x57\x87\x3b\x83\xfc\x73\x0e\x37\x7c\xe2\x04\x87\x7b\xfa\x2b\x67\xc1\xdd\xba\xcf\xfd\x0c\x00\x00\xff\xff\xce\x40\x8a\x4e\x18\x0c\x00\x00")

func partialsApplicationsHtmlBytes() ([]byte, error) {
	return bindataRead(
		_partialsApplicationsHtml,
		"partials/applications.html",
	)
}

func partialsApplicationsHtml() (*asset, error) {
	bytes, err := partialsApplicationsHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "partials/applications.html", size: 3096, mode: os.FileMode(420), modTime: time.Unix(1457866204, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _partialsNavbarHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x51\xc1\x4e\xc3\x30\x0c\xbd\xf3\x15\x91\x77\xd8\xa9\xca\x1d\xb5\x95\x38\x22\x21\x84\x80\x1f\xf0\x1a\xb7\xb3\x14\x25\x55\x9a\x4e\x4c\x88\x7f\xc7\x5d\xd3\x30\x56\x10\x70\xe0\x94\xf8\xf9\xd9\xcf\xcf\x2e\x1d\x1e\x54\x63\x71\x18\x2a\x90\xef\x0e\x83\x9a\x9f\xc2\x50\x8b\xa3\x8d\x50\x5f\x29\x55\x1a\xce\xac\xc6\xbb\x88\xec\x28\x14\xad\x1d\xd9\x9c\xf2\x9f\x19\xa9\xc1\x9e\xd0\x50\x48\x79\x61\xec\xc6\x18\xbd\x53\xf1\xd8\x53\x05\x73\x00\x17\x25\xd1\x77\x9d\x25\xd5\x78\x6b\xb1\x1f\xc8\x80\x32\x18\x31\xc1\x93\xf4\x8c\x2f\x30\x86\x8e\x62\x05\x9b\x54\xfd\x91\xc6\xc0\x58\xd0\x4b\x8f\xce\x90\xa9\xa0\x45\x2b\x68\x1e\x64\x10\x7c\x11\x1e\x42\xe1\x9d\x3d\x42\xfd\x3c\x4b\x4b\x2b\xee\x30\xb2\x77\xa5\x9e\x78\x5f\x16\xb1\xec\xa0\x10\x45\xa8\xff\x8d\xa4\xe7\x05\xe5\x18\x2f\x36\xb5\x0b\xe2\x0d\xd4\x3e\x50\x5b\x81\x86\xfa\xce\x3f\xa2\x7a\xa2\x70\xa0\x50\x6a\x4c\x37\xd1\x72\x94\xf5\x79\x96\x35\xa9\xd5\xda\xd8\xe4\xfe\x19\xcc\x13\x8c\xf6\x6c\x84\xa5\x56\x9e\xcc\x10\x8e\x65\xe5\xba\x22\xd1\x5e\xb1\x89\x7c\xa0\x6b\xd5\x63\x47\xaa\xaa\xd4\x16\xfb\xde\x72\x73\xda\xee\xb0\x7d\x13\xcf\xb8\x18\xd8\xe8\xf3\x1c\xd4\x37\x67\xd1\xe4\xa7\xd4\x96\x7f\xaf\xe3\xbc\xa1\x95\xc0\x09\x84\xfa\x7e\x7a\x2e\x5b\x96\x7a\xb4\x3f\x1a\x5d\xbe\x81\xbb\x7d\xfc\x93\x6b\x5e\x9b\x65\xf1\xf8\x70\xfb\xfd\x1c\xf9\x76\xe9\x53\x6a\x51\xaf\xdf\x03\x00\x00\xff\xff\x41\x3c\x84\x3a\xae\x03\x00\x00")

func partialsNavbarHtmlBytes() ([]byte, error) {
	return bindataRead(
		_partialsNavbarHtml,
		"partials/navbar.html",
	)
}

func partialsNavbarHtml() (*asset, error) {
	bytes, err := partialsNavbarHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "partials/navbar.html", size: 942, mode: os.FileMode(420), modTime: time.Unix(1458333988, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _partialsNodesHtml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x5b\xeb\x6e\xdb\xc8\x15\xfe\x9f\xa7\x98\x72\x8b\xac\x17\x58\x89\x6b\x63\xd3\x14\xa9\xc4\xc2\x2b\x25\xc0\x22\x96\x13\x48\xbb\x40\x80\xa2\x3f\x46\xe4\x48\x1a\x98\x1c\xb2\xe4\x50\x96\x60\xf8\xdd\x7b\xe6\x46\x0e\x6f\xb2\x24\xcb\x76\xd0\xba\x40\xaa\xb9\x1c\xce\x9c\xef\x5c\xbe\x39\xe4\x78\x07\x6c\xd9\xa3\xcc\x0f\xf3\x80\xa0\x2c\xf5\x87\xce\x8f\x09\x4e\x39\xc5\x61\xe6\x32\xbc\x9e\xe3\xb4\xbf\xe2\x51\xf8\xa3\xe3\x0d\xdc\x52\xd2\x7b\xf3\x66\x10\xd0\x35\xf2\x43\x9c\x65\x43\xc7\x8f\x19\xc7\x94\x91\xd4\xf1\xde\x20\x34\x58\x5d\x78\xd7\x71\x40\xb2\x81\x0b\x2d\x31\xc0\xf1\x3c\x24\x46\x58\x75\xe4\xff\xf7\x56\xf1\x5a\x3f\x24\xa4\x56\x04\x07\xaa\x2d\x7b\xde\x98\xac\x3f\xfe\xf9\xfb\xc0\x85\xa6\x35\x7a\x99\x24\xed\xa3\x9f\xc9\xb6\x36\x5a\x6c\x49\x36\xbc\x97\xd2\xe5\x8a\x3b\xde\xa5\xcf\x69\xcc\xb2\x52\x52\xb4\x8a\x8d\x07\x7c\x1e\x07\xdb\x72\x89\x14\x01\xe8\x94\x24\x04\xf3\xa1\xc3\x00\x14\xa2\x0c\x89\xdf\xcc\x31\x42\x42\x2c\xf0\xee\xee\xe4\x70\x3f\x90\x3a\xa3\xfb\x7b\x58\x36\x68\x17\xc1\x12\xc0\x43\x22\x80\xa6\x4d\xa4\x0d\x52\x31\x0f\x12\xd8\x08\xcc\x39\x43\xf0\xaf\x97\x45\xf2\x27\x20\x0b\x9c\x87\xdc\x11\x78\xfc\x90\xfa\x37\x43\x27\xc2\xfe\x28\x0a\xb2\x33\xb1\xe3\xcf\x68\x01\x1e\x27\x3f\x39\xde\xe4\x72\x84\xfc\x38\x8a\x30\x0b\xc0\x48\xf8\xd8\xc5\x33\x92\x65\x60\x67\xb9\x38\xac\x3a\x53\x5d\xe4\xa2\xcb\xdf\xbe\x1e\xbc\xec\x2a\x25\x8b\xa1\xe3\xfe\xe0\x4a\xcb\xbb\x0d\x5b\x3b\x28\xc0\x1c\xf7\x78\xbc\x5c\x86\x04\x80\xc5\x01\x0e\xcd\x18\x4e\x97\x04\x7c\xf7\x43\xb4\x9d\xc8\x61\xef\x63\x40\x79\x5d\x83\x79\xce\x39\x68\xd7\xa1\x06\x66\x4b\x88\x52\x0b\x5c\x40\x42\xc2\x89\xc1\x36\x96\xbd\xca\x92\xb6\xdb\xa0\x9d\x16\x91\x66\xa2\x0b\x9a\x22\x01\x20\x8f\xa0\x6d\x65\x92\xd8\x6f\x99\xc6\x79\xe2\xa0\x34\x16\x58\x74\x07\xa7\x14\xf7\x42\x3c\x27\xe1\xd0\xe9\xf7\xfb\x26\x65\xb4\xde\x7c\x9b\x80\xa8\xea\x38\x75\x14\x2d\xce\xf1\x53\x88\x67\x72\x06\xba\x8f\x64\x4b\x9a\x73\xe0\xaa\x05\x94\x7a\xa0\x93\xf7\x46\xff\x54\x72\x5d\x1a\x17\xa2\x25\x20\x0e\xa2\xc1\xd0\x21\x60\x4e\x65\x59\x91\xd3\x94\x05\x64\x33\x74\x7a\xe7\x46\xff\x00\x98\x24\x5e\xda\x00\x42\x12\xcc\xb7\xb0\x8e\xf2\xc7\x95\x18\x52\xb4\x51\xdf\xa4\x67\x9e\xd5\x2b\xc5\x7e\x1e\x11\x66\xe2\xbd\x29\x2e\x58\xa8\x9c\x6f\x93\x10\x79\x5e\xf0\xcd\x83\x06\xf4\xc3\x38\x23\x3a\x8c\x02\x9a\x45\xb4\x58\xa8\xea\x8f\x91\x94\xf3\x06\x59\x82\x99\x9a\x58\xd1\x20\x20\x0c\x92\x34\xcd\x61\xe2\x2d\xa7\x11\xc9\xfe\x31\x70\x85\x80\x67\x9b\x59\x2b\xb1\xfa\xb5\xaa\x26\xa7\x3c\xd4\xd6\xad\x5a\x49\x84\x2e\x6a\xa1\x91\xd5\xaf\x65\xb0\x49\x8f\x75\x19\x40\x84\x9f\x0d\x7f\x11\xa7\x51\x25\x11\xac\x27\xc4\x9c\x8e\x45\x5b\x04\x84\x24\x6e\x04\xf3\x43\x47\x29\xe1\xd8\x07\x01\x78\x4b\x59\xc6\xd1\x54\xfd\x61\xe0\xca\x7e\x6d\x15\xca\x92\x9c\x6b\xcb\x0b\x36\x73\x2a\x1b\xeb\x95\x94\x15\xcc\x2e\x10\xc1\x80\x43\xd8\xdc\xb2\x40\x95\xff\x6c\xfc\xc7\x22\x02\xea\xdd\x81\x08\x66\x4f\x82\x48\xee\xd2\x44\x24\x86\x77\x20\x1a\xb8\xb6\xcf\x1e\x70\xf7\x22\x8e\xb9\x15\xef\xb6\x40\x92\x87\x00\x8a\x2c\x40\x5f\x71\x92\x58\x1c\x47\x81\x6a\x49\x9a\xc6\xf0\x1c\x84\x9a\x6c\xc9\x20\xab\xaa\x71\x10\xf3\xb4\xa5\x90\x27\xd3\xa6\x25\x1b\xf6\x59\x39\x49\x69\x84\xd3\xad\xcd\x69\x79\x12\xe0\x92\x93\x67\x78\x0d\x05\xc7\x4a\xa0\xca\xea\x7b\x58\x50\x8a\xe6\x01\x84\xa7\xc8\xf3\x95\xf2\x9e\x87\xf2\x2a\x07\xd4\xcb\xf1\x9c\x2a\x35\xba\x58\x41\x15\xaa\x8f\x67\x05\xb3\x4b\x8d\x15\xf4\xf0\xc9\x79\xee\x95\xb9\x8f\x43\xf4\xca\xdc\x27\x65\x6e\x5d\x8d\x3e\x3d\x73\xeb\x77\x92\x57\xea\x7e\x1e\xea\x16\xdf\x00\x90\xb6\xb9\xc8\x50\x59\xb9\x66\xf6\x1b\xf2\x8b\xf2\xf9\x65\x10\xa4\x3b\x08\x5d\x4c\x77\xb0\xc5\x01\x19\x24\xdb\x9b\x0c\xc9\xa4\x55\xaf\xeb\x56\xe4\xc3\x9b\xe9\x14\xde\xb3\xe3\x48\xef\x77\xc6\x32\xc8\x80\x25\x61\x24\xc5\xbc\x25\xc3\x0e\xe5\xaa\x02\xa5\x45\x56\xd2\x01\x72\xf4\xe4\xec\x3b\xdb\x4d\xbf\xb3\x53\xf1\xef\xac\x4e\xc0\x59\xdf\x8c\x9e\x16\x13\xbb\xbd\xd9\x85\xe9\x5a\x4d\x3f\x1e\x53\xb1\x4f\x05\x93\x19\x3d\x2d\xa6\xc5\x88\xf1\x3f\x93\x2e\x48\x9f\x60\x16\x9d\xe5\xc9\x4f\x7b\x80\x62\x79\x34\x27\xe9\x0e\x58\x66\xab\x0a\x2a\x3d\x78\x7a\x50\xe3\xf8\x96\xed\x86\x15\x80\xc4\xc9\x80\xa9\xed\x1a\xd0\xe4\xf0\x69\xc1\xbd\x58\xed\x69\xb8\x1a\x4e\x15\x9a\x89\xcf\x55\xc1\xff\x48\x09\x9a\x99\x02\x74\x1f\x60\xaf\x75\xdb\xce\x37\xee\x99\xf9\xd0\x9b\x3d\x6d\xf1\x26\xbe\x56\xeb\xcf\xd2\xaf\x05\xdc\xf3\x14\x70\x13\xec\x17\x57\x01\x88\xc7\x68\x2e\xca\x39\x20\x52\x68\xb6\x5c\x77\x1c\x5d\xcc\xe5\xa1\x99\x67\x78\x8d\xe0\x5f\x0f\x3c\x9b\x55\x49\x34\xa4\xda\x2d\x49\x4a\x84\x0e\x58\x5c\xe2\x38\x15\x4e\x90\x81\x29\x97\xb9\xc3\x3e\xa7\x6b\xf2\x01\xfd\x25\x5b\xc5\xb7\x33\x10\xbf\xaf\xb3\x10\xb6\xaf\x2a\xb4\xd0\xa4\x8c\xb0\x33\x73\x1b\x22\xbf\xb8\x32\xc8\x57\x64\xdf\x8b\x58\xc6\xa8\xdf\x24\xb8\x21\x3d\x54\x71\xa9\xba\x21\x22\xd0\xfe\x4e\xc6\x40\xff\x0a\x67\x1c\x0d\x87\x88\x01\x6f\xdc\xdf\xb7\x3c\x52\x47\x7b\x34\x58\x11\x37\xf2\x8e\x86\xa0\x50\x6c\x0a\x9a\x02\x89\xec\x44\x36\x70\xf3\xf0\x81\x7a\xbc\xa6\xc4\x82\x92\x30\xc8\x08\x58\x13\x47\x22\xbd\x72\xbe\x1d\x6d\x7d\x08\x3c\xaf\x86\x6c\x10\x12\xa8\x80\x03\x6f\x0c\x12\x48\x8a\xc0\xd6\x6a\xa8\x21\x69\x1d\x92\xc5\x7a\x5d\xc5\x1a\xb9\x45\xe5\x92\x68\x8d\xc3\x9c\xb4\x1f\x2e\x1d\xd5\x40\x44\x21\xc5\x7e\x81\x5f\x0c\x8c\x73\xfe\x6e\xd7\x51\x5a\xea\x52\xf1\xac\xb1\xbd\x7d\x16\x59\xb4\xd6\x1f\x9b\xe7\xfa\x13\xbc\x19\x8f\xd4\x12\x35\xd5\xf4\xff\x92\x10\xfb\x64\x15\x87\xc0\x2f\x43\xe7\x63\x94\x00\x2e\x88\xc9\x1b\x42\x12\x08\xce\x15\x70\x70\x9e\xa6\x22\x55\x25\x4c\x07\xb9\x35\x67\xb8\xc6\x1b\xb5\x71\x91\xab\xa0\x9b\x50\x74\xe8\xa4\x24\x09\x29\xb1\xd4\x72\xbc\x29\x0c\x6d\x3f\xa0\x2f\x9f\x9b\x87\xff\x71\xd5\x40\x3d\x2a\xc8\x9a\xfa\x64\x06\x59\x92\x67\xdd\x81\x21\x85\x90\x92\xda\x2b\x36\xec\x55\x3b\xc2\xa3\x61\xe6\x4a\x0c\xf8\x2b\xe2\xdf\xcc\xe3\x4d\x51\x2a\x59\xeb\x1d\xe6\x62\xb2\xd6\xcf\xb9\x2d\x5b\x4a\xc6\x9e\x92\xff\xe4\x04\xb2\xb0\x0a\x53\x4e\xd5\x31\x9a\xe8\x7d\x8c\x73\x0b\x85\x8c\x73\x7f\xc3\x1c\x4a\x18\x68\x00\x11\x35\xa4\xfa\x7a\x16\x18\xff\xe7\x9a\x36\x33\xba\x64\x70\x62\x4f\x70\xba\xa4\xac\xe3\x69\x35\xd9\x2c\x6e\x4e\x15\x3f\x21\x65\x37\x97\xe3\x69\x67\xe8\x5c\xc1\x3c\xba\x0c\x70\x22\x18\x53\xd2\x2c\x52\x6f\xda\x1d\x41\x24\x74\x10\x1e\xef\x5a\xb7\x8d\xe9\xcc\x8e\x56\xfc\xc1\x3e\x53\xd8\x06\x8a\x75\xb1\xa3\x68\x76\x31\x8f\x7c\x74\xef\x77\x91\x83\x42\xef\x4a\x61\xe8\x8f\x8d\x36\x2a\x96\x8b\x5e\x85\xdf\x2e\xde\xbd\x6b\xb0\x46\xab\x9b\xf6\xb7\x02\xdf\x7c\x8d\x6f\x45\x39\xf4\xc7\x37\xd9\x78\x41\x13\x68\x0d\x94\x05\x8c\x5e\x4f\x6e\x00\x51\x17\x33\x12\x66\x13\x0c\xc9\x36\xd2\x1d\xc8\x97\xec\x66\xa7\x29\x64\x9d\x27\x2f\x9e\xca\xe7\x6f\xda\xad\xd1\xf6\x7c\xa1\x83\xf5\x87\x2c\x54\xfc\x15\x4b\x9b\x61\x46\x2b\xb1\x38\xe2\x29\xf6\x6f\xd0\x7c\x8b\xfe\x2a\x0b\x6c\xc7\xe8\x2f\x26\x21\xaf\xd5\xa8\xf8\xbb\x8b\xf6\x0d\x3b\xc9\xf3\x28\x67\x29\x9d\xfe\xa5\x36\xfd\xb7\xa3\x8e\x34\x59\x26\x95\x7a\x68\xfb\x34\xf4\x6b\xf3\xa0\x65\x56\xcf\x16\xd6\x05\x73\xbb\x09\x77\xf8\xa7\xe3\xb1\x47\x05\x8a\x00\x32\x62\x3c\x2c\xe2\x44\x86\x09\x1a\x29\x3f\xbf\x60\xe6\x4c\x49\x90\x33\x78\x75\xf5\xb7\xda\x2f\x52\x4b\xcb\xfc\xaa\xdf\xa8\x95\x4e\x9b\x49\x6c\x0e\x67\x15\x94\x74\x12\x1e\x8a\x17\x08\xba\x04\x5e\x63\xc4\x1f\x76\x7d\x1f\xc6\xb9\x96\x1a\xaa\x8f\x7b\xaa\x79\xb4\x49\x5a\x06\xf7\x3b\xe2\xad\x13\x58\x6b\x28\x53\x56\x1f\xf3\xcd\xd9\xbe\x8e\x35\xe1\xc6\xcb\xd1\x67\xf4\xf6\x6d\x43\xc2\x9c\x1d\x1d\xd3\x92\x48\xc5\xdc\x3f\x91\xf3\xe5\xb3\x83\x3e\x20\xe7\x13\xa6\x60\xce\x27\x3b\xf0\xd3\xcd\xc5\x8c\xf0\x86\x58\x79\xe2\x4f\xbf\x7d\xc5\x29\x8e\xa4\x10\x54\x56\xbb\x4e\xfa\x16\x7f\x58\x31\x07\x3b\x15\x27\x27\xac\x7a\x81\xc6\xa6\x7a\xe8\x8e\xb8\x27\x8a\x37\xd8\x5d\xe2\xe9\x8f\xaf\xc0\x12\x9c\xb2\xa5\x1c\xab\x1e\xec\xb6\xbe\x7b\x04\x5f\x7b\xe8\xed\x61\x95\xf3\xf1\x34\x5e\x2c\x40\x0f\x61\x95\xf3\xd2\x2a\x48\x8d\x7e\x27\xc6\x29\xb5\xd4\xc6\xb1\x06\x2a\xc6\xf9\xfb\xe9\x6c\xb3\x48\x45\x25\x0f\x6c\xa0\xe2\xe5\x93\xe9\xbe\x9c\x49\x0a\x15\xf4\x67\xec\xb2\x5b\x0d\x90\xbf\xbd\x7f\xff\xfe\x62\xdf\x30\x39\xec\x7d\xc3\xe8\xd2\xc6\x43\x85\x9e\x9a\x88\x6a\x2c\x53\x4c\x5b\xb1\xde\x2d\x52\x38\xf8\x59\xf9\x88\x91\x5b\xad\x7c\x27\x23\x5d\x17\x22\xc7\x93\x91\xbf\xfa\x5d\x16\x66\x45\x79\x20\xbb\xcf\x1e\x58\x25\x14\x70\x99\xd2\x48\x97\x02\xba\x23\xc3\xea\xfc\x09\x78\x47\xda\xb9\x0c\x67\xef\xe5\x92\xcb\xb2\x81\x50\x42\x9f\xf9\x55\xf5\x1e\x93\x5e\xfb\x58\x03\x96\x17\x2f\xa7\x13\xa8\xea\xf7\x78\xc3\x7c\x7a\x43\x4c\xa4\x3e\xea\x8b\xb2\x6a\x3e\xd9\x01\x04\x0b\x4a\xe8\x78\xf3\x9d\x40\x97\xfa\xe8\x1b\x8b\xcd\xd1\xd0\x0f\x23\xd5\x6b\x8b\x75\x9a\xb4\x5a\xc9\x52\xf9\x5b\x44\xe6\x97\x0a\x7b\x5a\x82\x86\x5f\xa7\xe2\x2a\x07\xa4\x86\x43\xf4\xcb\x73\x56\x75\x7f\x50\x30\xdb\xf2\xa1\xd2\xce\x92\x7a\x4c\x61\x57\xd9\xac\xb2\xea\x0b\x1c\xd5\xd6\xee\xfd\x31\x09\xf1\xd6\x54\x2c\xb6\x96\xcf\x11\x54\x15\x55\x0e\xf8\xf4\xfb\x7f\x7a\x5f\xda\xe1\xe6\x7d\xaf\x51\xed\xfb\x18\xf5\x1f\xe4\x58\x61\xf1\xe8\x9b\xd5\xff\x06\x00\x00\xff\xff\x08\x2a\xba\x86\xe5\x35\x00\x00")

func partialsNodesHtmlBytes() ([]byte, error) {
	return bindataRead(
		_partialsNodesHtml,
		"partials/nodes.html",
	)
}

func partialsNodesHtml() (*asset, error) {
	bytes, err := partialsNodesHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "partials/nodes.html", size: 13797, mode: os.FileMode(420), modTime: time.Unix(1460967341, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"index.html": indexHtml,
	"js/app.js": jsAppJs,
	"partials/api.html": partialsApiHtml,
	"partials/applications.html": partialsApplicationsHtml,
	"partials/navbar.html": partialsNavbarHtml,
	"partials/nodes.html": partialsNodesHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"index.html": &bintree{indexHtml, map[string]*bintree{}},
	"js": &bintree{nil, map[string]*bintree{
		"app.js": &bintree{jsAppJs, map[string]*bintree{}},
	}},
	"partials": &bintree{nil, map[string]*bintree{
		"api.html": &bintree{partialsApiHtml, map[string]*bintree{}},
		"applications.html": &bintree{partialsApplicationsHtml, map[string]*bintree{}},
		"navbar.html": &bintree{partialsNavbarHtml, map[string]*bintree{}},
		"nodes.html": &bintree{partialsNodesHtml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

