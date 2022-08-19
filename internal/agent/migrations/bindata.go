package migrations

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

var __000001_create_secrets_table_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\x4e\x4d\x2e\x4a\x2d\x29\xb6\x06\x04\x00\x00\xff\xff\x27\x14\x07\xb5\x1d\x00\x00\x00")

func _000001_create_secrets_table_down_sql() ([]byte, error) {
	return bindata_read(
		__000001_create_secrets_table_down_sql,
		"000001_create_secrets_table.down.sql",
	)
}

var __000001_create_secrets_table_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x24\xc9\xb1\x0a\x42\x21\x14\x06\xe0\xdd\xa7\xf8\xc7\x82\xde\xa0\xe9\x28\x27\x92\x2c\xe3\x78\x08\x1d\x2d\x1d\x02\xa7\xf4\xfd\x09\xee\x5d\xbf\xcf\x09\x93\x32\x94\x6c\x60\xf8\x0b\x1e\x51\xc1\xd9\x27\x4d\x98\xfd\xf3\xeb\x6b\xe2\x60\x00\xe0\xdb\xf0\x22\x71\x57\x12\x3c\xc5\xdf\x49\x0a\x6e\x5c\x4e\xdb\x8d\xfa\xee\x63\x42\x39\xeb\x0e\xad\xae\x0a\x1b\xa2\x35\xc7\xf3\x3f\x00\x00\xff\xff\x37\xf3\x93\x82\x62\x00\x00\x00")

func _000001_create_secrets_table_up_sql() ([]byte, error) {
	return bindata_read(
		__000001_create_secrets_table_up_sql,
		"000001_create_secrets_table.up.sql",
	)
}

var __000002_create_dates_columns_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x4e\x4d\x2e\x4a\x2d\x29\xe6\x72\x09\xf2\x0f\x50\x70\xf6\xf7\x09\xf5\xf5\x53\x48\x2e\x4a\x4d\x2c\x49\x4d\x89\x4f\x2c\xb1\xe6\xe2\x22\xa4\xb8\xb4\x20\x85\x78\xc5\x29\xa9\x39\xa9\x30\xc5\x80\x00\x00\x00\xff\xff\x86\x0e\x9b\x9e\x86\x00\x00\x00")

func _000002_create_dates_columns_down_sql() ([]byte, error) {
	return bindata_read(
		__000002_create_dates_columns_down_sql,
		"000002_create_dates_columns.down.sql",
	)
}

var __000002_create_dates_columns_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x4e\x4d\x2e\x4a\x2d\x29\xe6\x72\x74\x71\x51\x70\xf6\xf7\x09\xf5\xf5\x53\x48\x2e\x4a\x4d\x2c\x49\x4d\x89\x4f\x2c\x51\x08\x71\x8d\x08\xb1\xe6\xe2\x22\xa0\xa1\xb4\x20\x05\x59\x83\x82\x8b\xab\x9b\x63\xa8\x4f\x88\x82\x5f\xa8\x8f\x0f\x61\xdd\x29\xa9\x39\xa9\x38\x75\x03\x02\x00\x00\xff\xff\x6f\x0d\x26\xb9\xab\x00\x00\x00")

func _000002_create_dates_columns_up_sql() ([]byte, error) {
	return bindata_read(
		__000002_create_dates_columns_up_sql,
		"000002_create_dates_columns.up.sql",
	)
}

var __000003_create_accounts_table_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4c\x4e\xce\x2f\xcd\x2b\x29\xb6\x06\x04\x00\x00\xff\xff\x12\x35\x0c\x5f\x1e\x00\x00\x00")

func _000003_create_accounts_table_down_sql() ([]byte, error) {
	return bindata_read(
		__000003_create_accounts_table_down_sql,
		"000003_create_accounts_table.down.sql",
	)
}

var __000003_create_accounts_table_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8e\xd1\x4a\x87\x30\x1c\x85\xef\xff\x4f\x71\x2e\x15\x7c\x83\xae\x36\xfb\x49\xa3\xb5\xd5\xdc\x42\x2f\x87\xfe\x0a\xa1\x34\x36\xad\xd7\x0f\x1c\x04\x79\xfd\x7d\x9c\xef\xb4\x8e\x84\x27\x78\x21\x35\x41\x75\x30\xd6\x83\x06\xd5\xfb\x1e\x71\x9a\xb6\x63\xdd\x33\xaa\x1b\x00\x2c\x33\x5e\x85\x6b\x1f\x84\xc3\xb3\x53\x4f\xc2\x8d\x78\xa4\xb1\x39\x59\xe6\xf4\xcd\x09\x9e\x06\x7f\x2e\x98\xa0\x35\x82\x51\x2f\x81\x8a\x70\x64\x4e\x6b\xfc\xe4\xff\x4a\x61\x5f\x31\xe7\x9f\x2d\xcd\x90\xda\xca\x0b\x4b\xfc\xbe\xe4\x9d\x13\xcf\x90\xd6\x6a\x12\x06\xf7\xd4\x89\xa0\x3d\xde\xe2\x47\xe6\x8b\x5e\x92\xa8\xca\x9f\xe6\x2f\x5b\xdf\xea\xbb\xdf\x00\x00\x00\xff\xff\xb0\xb9\x7a\x23\xea\x00\x00\x00")

func _000003_create_accounts_table_up_sql() ([]byte, error) {
	return bindata_read(
		__000003_create_accounts_table_up_sql,
		"000003_create_accounts_table.up.sql",
	)
}

var __000004_create_sync_columns_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x4e\x4d\x2e\x4a\x2d\x29\xe6\x72\x09\xf2\x0f\x50\x70\xf6\xf7\x09\xf5\xf5\x53\x28\xae\xcc\x4b\x4e\x4d\xb1\xe6\xe2\x22\xa4\x30\x25\x35\x27\xb5\x04\xa4\x12\x10\x00\x00\xff\xff\x70\x52\x0a\xf2\x52\x00\x00\x00")

func _000004_create_sync_columns_down_sql() ([]byte, error) {
	return bindata_read(
		__000004_create_sync_columns_down_sql,
		"000004_create_sync_columns.down.sql",
	)
}

var __000004_create_sync_columns_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\xcc\xb1\x0d\x02\x31\x0c\x05\xd0\x3e\x53\xfc\x3d\xa8\x1c\x62\x2a\x63\x4b\xc8\x19\x00\x25\xa6\x8a\x28\x70\x1a\xb6\xbf\x11\xee\x16\x78\x24\xce\x2f\x38\x55\x61\x64\x8c\x5f\xec\x2c\xd4\x1a\xee\x26\xfd\xa9\xc8\xff\x77\xc4\x44\x35\x13\x26\x45\xe3\x07\x75\x71\x7c\xde\x2b\x03\x6a\x0e\xed\x22\xb7\x52\x4e\x98\x19\x2b\xf6\x05\xe7\x08\x00\x00\xff\xff\x49\x18\xb5\xd0\x8e\x00\x00\x00")

func _000004_create_sync_columns_up_sql() ([]byte, error) {
	return bindata_read(
		__000004_create_sync_columns_up_sql,
		"000004_create_sync_columns.up.sql",
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
	"000001_create_secrets_table.down.sql":  _000001_create_secrets_table_down_sql,
	"000001_create_secrets_table.up.sql":    _000001_create_secrets_table_up_sql,
	"000002_create_dates_columns.down.sql":  _000002_create_dates_columns_down_sql,
	"000002_create_dates_columns.up.sql":    _000002_create_dates_columns_up_sql,
	"000003_create_accounts_table.down.sql": _000003_create_accounts_table_down_sql,
	"000003_create_accounts_table.up.sql":   _000003_create_accounts_table_up_sql,
	"000004_create_sync_columns.down.sql":   _000004_create_sync_columns_down_sql,
	"000004_create_sync_columns.up.sql":     _000004_create_sync_columns_up_sql,
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
	"000001_create_secrets_table.down.sql":  &_bintree_t{_000001_create_secrets_table_down_sql, map[string]*_bintree_t{}},
	"000001_create_secrets_table.up.sql":    &_bintree_t{_000001_create_secrets_table_up_sql, map[string]*_bintree_t{}},
	"000002_create_dates_columns.down.sql":  &_bintree_t{_000002_create_dates_columns_down_sql, map[string]*_bintree_t{}},
	"000002_create_dates_columns.up.sql":    &_bintree_t{_000002_create_dates_columns_up_sql, map[string]*_bintree_t{}},
	"000003_create_accounts_table.down.sql": &_bintree_t{_000003_create_accounts_table_down_sql, map[string]*_bintree_t{}},
	"000003_create_accounts_table.up.sql":   &_bintree_t{_000003_create_accounts_table_up_sql, map[string]*_bintree_t{}},
	"000004_create_sync_columns.down.sql":   &_bintree_t{_000004_create_sync_columns_down_sql, map[string]*_bintree_t{}},
	"000004_create_sync_columns.up.sql":     &_bintree_t{_000004_create_sync_columns_up_sql, map[string]*_bintree_t{}},
}}
