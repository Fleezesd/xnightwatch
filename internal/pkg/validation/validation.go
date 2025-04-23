package validation

import (
	"context"
	"reflect"
	"strings"

	"github.com/fleezesd/xnightwatch/internal/pkg/middleware/validate"
	"github.com/google/wire"
	"k8s.io/klog/v2"
)

// validator implement the validate.IValidator interface.
type validator struct {
	registry map[string]reflect.Value
}

// ProviderSet is validator providers.
var ProviderSet = wire.NewSet(New, wire.Bind(new(validate.IValidator), new(*validator)))

// New create and initialize the custom validator
func New(cv any) *validator {
	return &validator{
		registry: GetValidateFuns(cv),
	}
}

func (vd *validator) Validate(ctx context.Context, req any) error {
	m, ok := vd.registry[reflect.TypeOf(req).Elem().Name()]
	if !ok {
		return nil
	}

	val := m.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
	if !val[0].IsNil() {
		return val[0].Interface().(error)
	}
	return nil
}

func GetValidateFuns(cv any) map[string]reflect.Value {
	funcs := make(map[string]reflect.Value)
	typeOf := reflect.TypeOf(cv)
	valueOf := reflect.ValueOf(cv)
	for i := 0; i < typeOf.NumMethod(); i++ {
		m := typeOf.Method(i)
		val := valueOf.MethodByName(m.Name)

		// make sure method is valid
		if !val.IsValid() {
			continue
		}

		// make sure method has "validate" prefix
		if !strings.HasPrefix(m.Name, "Validate") {
			continue
		}

		// method must have 2 input parameters and 1 return value.
		typ := val.Type()
		if typ.NumIn() != 2 || typ.NumOut() != 1 {
			continue
		}

		// first params must be context.Context
		if typ.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
			continue
		}

		// second params must be pointer
		if typ.In(1).Kind() != reflect.Pointer {
			continue
		}

		// 获取第二个参数（指针类型）指向的结构体名称
		vName := typ.In(1).Elem().Name()
		// 方法名必须是 "Validate" + 结构体名, 例如：ValidateUser 方法必须用于验证 User 结构体
		if m.Name != ("Validate" + vName) {
			continue
		}

		// 返回值类型必须是 error 接口
		if typ.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		klog.V(4).InfoS("Register validator", "validator", vName)
		funcs[vName] = val
	}
	return funcs
}
