package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func raw_data_yml() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x2a, 0xc8,
		0xcf, 0xc9, 0xc9, 0xcc, 0x4b, 0xf7, 0xc9, 0x4f, 0x4e, 0x2c, 0xc9, 0xcc,
		0xcf, 0xb3, 0xe2, 0x52, 0x50, 0x28, 0x49, 0xad, 0x28, 0x01, 0xd1, 0x0a,
		0x0a, 0xa9, 0x79, 0x10, 0x5a, 0x41, 0x21, 0x23, 0xbf, 0xb4, 0xa8, 0xd8,
		0x4a, 0x41, 0xc9, 0x35, 0x38, 0xc0, 0xd1, 0xef, 0x4c, 0xb3, 0xbf, 0x8f,
		0x12, 0x54, 0x02, 0xa4, 0xdf, 0x0a, 0x4c, 0x16, 0x43, 0xb4, 0x14, 0xa3,
		0x69, 0xc9, 0xc8, 0x2f, 0x4a, 0x2c, 0xc6, 0x50, 0x9c, 0xcf, 0x05, 0xb2,
		0xa8, 0x28, 0x33, 0x3d, 0x3d, 0xb5, 0xa8, 0x18, 0xdd, 0x32, 0x5d, 0x85,
		0xfc, 0xbc, 0x54, 0x38, 0xbb, 0xa4, 0x3c, 0x1f, 0xc1, 0xce, 0x28, 0x4a,
		0x4d, 0x45, 0xb3, 0x47, 0x57, 0xa1, 0x34, 0x0f, 0xa1, 0x22, 0x25, 0xbf,
		0x18, 0xa1, 0xba, 0x28, 0xb5, 0x98, 0x0b, 0x10, 0x00, 0x00, 0xff, 0xff,
		0xf1, 0x68, 0x26, 0xde, 0xe1, 0x00, 0x00, 0x00,
	},
		"raw/data.yml",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"raw/data.yml": raw_data_yml,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"raw": &_bintree_t{nil, map[string]*_bintree_t{
		"data.yml": &_bintree_t{raw_data_yml, map[string]*_bintree_t{}},
	}},
}}
