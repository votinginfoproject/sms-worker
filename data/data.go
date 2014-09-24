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
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0xcc, 0x55,
		0xcf, 0x4e, 0xdc, 0x3e, 0x10, 0xbe, 0xf3, 0x14, 0x23, 0xce, 0x10, 0xee,
		0xb9, 0xf1, 0xfb, 0x29, 0xc0, 0x4a, 0xe9, 0xb2, 0x62, 0x69, 0xab, 0x1e,
		0x87, 0x64, 0x36, 0x71, 0x71, 0xec, 0xd4, 0x76, 0x08, 0xdb, 0x53, 0x25,
		0x9e, 0xa2, 0x8f, 0xd0, 0xe7, 0xe8, 0x95, 0xa7, 0xe8, 0x93, 0x74, 0x6c,
		0x27, 0x29, 0xec, 0x9f, 0xaa, 0x87, 0x16, 0x21, 0x21, 0x12, 0xdb, 0xe3,
		0x99, 0x6f, 0xbe, 0xef, 0x9b, 0x6c, 0xab, 0xa5, 0x14, 0xaa, 0xca, 0x75,
		0x81, 0x4e, 0x68, 0x95, 0x1e, 0x00, 0x38, 0xba, 0x77, 0xfe, 0x09, 0x40,
		0x2a, 0x3e, 0x01, 0x5a, 0x43, 0x2b, 0x71, 0x9f, 0xc2, 0xe1, 0x07, 0xdd,
		0x19, 0x18, 0x2e, 0x41, 0x2b, 0xb1, 0x20, 0x10, 0x36, 0x3d, 0x1c, 0xc2,
		0x6a, 0x3e, 0xb5, 0x1c, 0x75, 0x11, 0x9e, 0x71, 0x97, 0xec, 0x56, 0x12,
		0xdb, 0xa2, 0x12, 0xb6, 0x3e, 0xfe, 0xb3, 0x64, 0x63, 0xf4, 0xaf, 0xa4,
		0xce, 0x88, 0xaa, 0x22, 0x63, 0x37, 0x51, 0x1e, 0x87, 0x64, 0x1b, 0x55,
		0x8f, 0xc1, 0x86, 0xdd, 0xa2, 0x46, 0x55, 0x51, 0xce, 0xff, 0x3a, 0xac,
		0x68, 0x77, 0xa3, 0xf1, 0xd6, 0xfe, 0xf4, 0xa4, 0x2a, 0xc9, 0x50, 0xb6,
		0x2a, 0x90, 0xc7, 0xf8, 0xf8, 0xa0, 0xe5, 0xf3, 0x1d, 0x5e, 0xe3, 0x8d,
		0xee, 0xdc, 0x3e, 0x56, 0x51, 0xca, 0x14, 0x98, 0x84, 0x1f, 0x5f, 0xbe,
		0x1a, 0x82, 0xce, 0x7a, 0x1a, 0x5c, 0x4d, 0xf0, 0x4e, 0x3b, 0xff, 0x3a,
		0x53, 0x2b, 0x6d, 0x9a, 0xa0, 0x0b, 0x2c, 0x8c, 0xfe, 0x48, 0x85, 0x03,
		0xa9, 0xf5, 0x6d, 0xd7, 0x82, 0xd3, 0x5a, 0x26, 0x70, 0xad, 0x41, 0x12,
		0x1a, 0x05, 0x8d, 0x36, 0x74, 0x04, 0x77, 0xc2, 0x0a, 0x07, 0xb5, 0x73,
		0x6d, 0x7a, 0x72, 0xd2, 0xf7, 0x7d, 0x72, 0x17, 0xf2, 0x08, 0x4e, 0xd3,
		0xc6, 0xeb, 0x89, 0x36, 0xd5, 0x06, 0xf8, 0x80, 0xe1, 0x89, 0x20, 0x2f,
		0x8d, 0x65, 0x3f, 0xd9, 0x7d, 0x8d, 0x6e, 0x5a, 0x04, 0x1e, 0xb7, 0x78,
		0x67, 0xc7, 0x98, 0x02, 0xa7, 0xe5, 0xa7, 0x8e, 0x9e, 0xbc, 0x7f, 0xff,
		0x76, 0x50, 0x93, 0x6c, 0xf7, 0x91, 0xdf, 0x90, 0xea, 0x52, 0xb8, 0xe6,
		0x13, 0x58, 0x5c, 0xe6, 0x39, 0xb7, 0x01, 0x2b, 0xa1, 0x4a, 0x58, 0x6f,
		0x99, 0x32, 0x89, 0x51, 0x97, 0x67, 0x67, 0xb3, 0xff, 0x67, 0xa7, 0x39,
		0x30, 0x13, 0x40, 0x92, 0x5b, 0xf0, 0x64, 0xe8, 0xd5, 0x4a, 0x14, 0x02,
		0x25, 0xf8, 0xd6, 0x86, 0xc8, 0xab, 0xec, 0x7c, 0xb6, 0xbc, 0xce, 0xae,
		0x42, 0xa4, 0xa1, 0x4a, 0x58, 0x67, 0x22, 0x75, 0x21, 0x68, 0x00, 0x20,
		0x07, 0x2b, 0xb2, 0xc7, 0xdf, 0xa0, 0x2a, 0x11, 0xb2, 0xe5, 0xe2, 0x74,
		0xfe, 0xf8, 0x70, 0x99, 0x43, 0x8b, 0x06, 0xa1, 0xc0, 0xe6, 0x46, 0xa0,
		0x2f, 0x05, 0xa2, 0x14, 0xba, 0xc1, 0x64, 0xa3, 0xff, 0xd8, 0xc1, 0xa8,
		0xdd, 0xeb, 0xe8, 0x24, 0xdc, 0xca, 0xe6, 0xe7, 0xf9, 0x6c, 0x79, 0xe1,
		0x81, 0xc4, 0xa9, 0x0b, 0x4e, 0x1a, 0xa3, 0x92, 0xdf, 0xaa, 0xee, 0x25,
		0x9b, 0x16, 0xbe, 0xc3, 0x6d, 0xd1, 0xd7, 0x5d, 0x89, 0x07, 0x64, 0x8c,
		0x8e, 0xf7, 0x77, 0x4e, 0x56, 0x59, 0x1a, 0xb2, 0x76, 0x81, 0xc6, 0x52,
		0x76, 0xcf, 0xa8, 0x99, 0x82, 0xb7, 0x96, 0x0c, 0x03, 0x64, 0x57, 0xf1,
		0x57, 0x46, 0xb1, 0xcf, 0x1d, 0x20, 0xf7, 0x54, 0xe8, 0x4a, 0x89, 0xcf,
		0x54, 0x42, 0xa1, 0x9b, 0x86, 0x75, 0x18, 0x3a, 0xbf, 0xc8, 0xf2, 0x85,
		0x6f, 0xc0, 0x12, 0xf9, 0x19, 0x01, 0xdd, 0xfa, 0xb6, 0x6d, 0xb2, 0xa3,
		0xc0, 0x9c, 0xfa, 0x98, 0xfb, 0x3d, 0x81, 0x22, 0x1a, 0x88, 0x1f, 0x22,
		0x7c, 0x0e, 0x76, 0xfc, 0x9d, 0x28, 0xc9, 0xef, 0x43, 0x2f, 0x5c, 0xcd,
		0x93, 0x53, 0x30, 0xcf, 0x71, 0x24, 0x02, 0x93, 0xc3, 0x68, 0x25, 0x3e,
		0x45, 0x29, 0xca, 0x88, 0x6e, 0xc2, 0xc6, 0xf4, 0x31, 0xe8, 0x21, 0xe1,
		0x11, 0x4b, 0x49, 0x68, 0x79, 0xd3, 0xac, 0x01, 0x2b, 0x14, 0x6a, 0xc4,
		0xa4, 0x74, 0x36, 0x68, 0xe9, 0xa7, 0x35, 0x85, 0xb9, 0x9e, 0xc4, 0xb5,
		0x80, 0x3c, 0xd5, 0xb6, 0xa8, 0xa9, 0xec, 0x24, 0x23, 0xf4, 0x72, 0x46,
		0x94, 0x86, 0x30, 0x0c, 0x2f, 0x1f, 0x15, 0xb7, 0x1c, 0x6f, 0xa9, 0xaf,
		0xc9, 0x4f, 0xb0, 0xa1, 0x56, 0xae, 0x23, 0x5c, 0xe4, 0xb6, 0xfa, 0xb1,
		0xfe, 0x58, 0xad, 0x22, 0x45, 0x06, 0xe5, 0x7f, 0x58, 0xdc, 0x92, 0x2a,
		0x53, 0x58, 0x6a, 0x63, 0xd6, 0x47, 0xd0, 0x13, 0xff, 0xf9, 0x2f, 0x88,
		0xc2, 0x1b, 0x49, 0xcf, 0xbd, 0x38, 0x59, 0xad, 0xc4, 0xf5, 0x64, 0x4c,
		0x39, 0xfc, 0xf4, 0x6c, 0xba, 0x7b, 0xbf, 0x84, 0x93, 0xe3, 0xff, 0xa1,
		0x94, 0x63, 0x8d, 0x57, 0x26, 0xe9, 0x08, 0xeb, 0x65, 0xa5, 0x1d, 0xab,
		0xfe, 0x0d, 0x89, 0x77, 0x4f, 0xbe, 0x97, 0xfd, 0x67, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x31, 0x72, 0x0b, 0x84, 0x8b, 0x08, 0x00, 0x00,
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
