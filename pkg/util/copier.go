package util

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/wendisx/puzzle/pkg/clog"
)

const (
	_default_ttl      = 30 * time.Minute
	_default_capacity = 1 << 7
)

// key(string): src_to_dest value(*Copier): pot some functions collection
var copierPool *ttlcache.Cache[string, *Copier]

type (
	Copier struct {
		fns []func(dest, src reflect.Value)
	}
)

func init() {
	copierPool = ttlcache.New(
		ttlcache.WithTTL[string, *Copier](_default_ttl),
		ttlcache.WithCapacity[string, *Copier](_default_capacity),
	)
}

func GetCopier(src, dest interface{}) *Copier {
	cKey := fmt.Sprintf("%v_to_%v", src, dest)
	if copierPool.Has(cKey) {
		return copierPool.Get(cKey).Value()
	}
	copier := NewCopier(reflect.TypeOf(src), reflect.TypeOf(dest))
	copierPool.Set(cKey, copier, ttlcache.DefaultTTL)
	return copier
}

func NewCopier(srcType, destType reflect.Type) *Copier {
	var fns []func(dest, src reflect.Value)
	if srcType.Kind() == reflect.Pointer {
		srcType = srcType.Elem()
	}
	if destType.Kind() == reflect.Pointer {
		destType = destType.Elem()
	}
	var buildFns func(currDestType reflect.Type, indexPrefix []int)
	buildFns = func(currDestType reflect.Type, indexPrefix []int) {
		for i := 0; i < currDestType.NumField(); i++ {
			df := currDestType.Field(i)
			fullIndex := append(indexPrefix, i)
			if df.Anonymous && df.Type.Kind() == reflect.Struct {
				buildFns(df.Type, fullIndex)
				continue
			}
			if !df.IsExported() {
				continue
			}
			sf, ok := srcType.FieldByName(df.Name)
			if !ok {
				continue
			}
			idx := fullIndex
			if sf.Type.AssignableTo(df.Type) {
				fns = append(fns, func(dest, src reflect.Value) {
					dest.FieldByIndex(idx).Set(src.FieldByName(df.Name))
				})
			} else if sf.Type.ConvertibleTo(df.Type) {
				fns = append(fns, func(dest, src reflect.Value) {
					dest.FieldByIndex(idx).Set(src.FieldByName(df.Name).Convert(df.Type))
				})
			}
		}
	}
	buildFns(destType, []int{})
	return &Copier{fns: fns}
}

func (c *Copier) Copy(dest, src interface{}) {
	destValue := reflect.ValueOf(dest)
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Pointer || srcValue.Kind() != reflect.Pointer {
		clog.Panic("[copier] - both dest and src must be pointers to struct")
	}

	destValue = destValue.Elem()
	srcValue = srcValue.Elem()

	if destValue.Kind() != reflect.Struct || srcValue.Kind() != reflect.Struct {
		clog.Panic("[copier] -  both dest and src must point to structs")
	}

	for _, fn := range c.fns {
		fn(destValue, srcValue)
	}
}

func CopySlice[S, D any](src []S) []D {
	if len(src) == 0 {
		clog.Warn("try to copy slice, but src's length is zero.")
		return []D{}
	}
	dest := make([]D, len(src))
	cp := GetCopier(src[0], dest[0])
	for i := range src {
		cp.Copy(&dest[i], &src[i])
	}
	return dest
}
