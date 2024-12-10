package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

//@function: RegisterRule
//@description: 注册自定义规则方案建议在路由初始化层即注册
//@param: key string, rule Rules
//@return: err error

func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

//@function: NotEmpty
//@description: 非空 不能为其对应类型的0值
//@return: string

func NotEmpty() string {
	return "notEmpty"
}

// @function: RegexpMatch
// @description: 正则校验 校验输入项是否满足正则表达式
// @param:  rule string
// @return: string

func RegexpMatch(rule string) string {
	return "regexp=" + rule
}

//@function: Lt
//@description: 小于入参(<) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Lt(mark string) string {
	return "lt=" + mark
}

//@function: Le
//@description: 小于等于入参(<=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Le(mark string) string {
	return "le=" + mark
}

//@function: Eq
//@description: 等于入参(==) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Eq(mark string) string {
	return "eq=" + mark
}

//@function: Ne
//@description: 不等于入参(!=)  如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Ne(mark string) string {
	return "ne=" + mark
}

//@function: Ge
//@description: 大于等于入参(>=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Ge(mark string) string {
	return "ge=" + mark
}

//@function: Gt
//@description: 大于入参(>) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Gt(mark string) string {
	return "gt=" + mark
}

//
//@function: Verify
//@description: 校验方法
//@param: st interface{}, roleMap Rules(入参实例，规则map)
//@return: err error

func Verify(st interface{}, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"lt": true,
		"le": true,
		"eq": true,
		"ne": true,
		"ge": true,
		"gt": true,
	}

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind() // 获取到st对应的类别
	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	num := val.NumField()
	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		if tagVal.Type.Kind() == reflect.Struct {
			if err = Verify(val.Interface(), roleMap); err != nil {
				return err
			}
		}
		if len(roleMap[tagVal.Name]) > 0 {
			for _, v := range roleMap[tagVal.Name] {
				switch {
				case v == "notEmpty":
					if isBlank(val) {
						return errors.New(tagVal.Name + "值不能为空")
					}
				case strings.Split(v, "=")[0] == "regexp":
					if !regexpMatch(strings.Split(v, "=")[1], val.String()) {
						return errors.New(tagVal.Name + "格式校验不通过")
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(tagVal.Name + "长度或值不在合法范围," + v)
					}
				}
			}
		}
	}
	return nil
}

//@function: compareVerify
//@description: 长度和数字的校验方法 根据类型自动校验
//@param: value reflect.Value, VerifyStr string
//@return: bool

func compareVerify(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String:
		return compare(len([]rune(value.String())), VerifyStr)
	case reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

//@function: isBlank
//@description: 非空校验
//@param: value reflect.Value
//@return: bool

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

//@function: compare
//@description: 比较函数
//@param: value interface{}, VerifyStr string
//@return: bool

func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	default:
		return false
	}
}

func regexpMatch(rule, matchStr string) bool {
	return regexp.MustCompile(rule).MatchString(matchStr)
}

func CheckName(str string) bool {
	if len(str) == 0 {
		return false
	}

	strRune := []rune(str)
	length := len(strRune)
	if len(strRune) > 20 {
		return false
	}

	if strRune[0] == '_' {
		return false
	}

	for i := 0; i < length; i++ {
		if unicode.IsLetter(strRune[i]) || unicode.IsDigit(strRune[i]) || unicode.Is(unicode.Han, strRune[i]) || strRune[i] == '_' {
			continue
		}
		return false
	}

	return true
}

func VolidatorRequest(req interface{}) (string, error) {
	validate := validator.New()
	validate.RegisterValidation("nameValidator", func(fl validator.FieldLevel) bool {
		name := fl.Field().Interface().(string)
		name = TrimSpace(name)
		return CheckName(name)
	})

	var errMsg string
	if err := validate.Struct(req); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += TranslateError(err) + ";"
		}
		return errMsg, err
	}
	return "", nil
}

// 自定义错误信息翻译函数
func TranslateError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s 字段是必填的", err.Field())
	case "gt":
		return fmt.Sprintf("%s 字段的值必须大于 %s", err.Field(), err.Param())
	case "min":
		return fmt.Sprintf("%s 字段的最小长度为 %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s 字段的最大长度为 %s", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("%s 字段必须是有效的电子邮件地址", err.Field())
	case "gte":
		return fmt.Sprintf("%s 字段的值必须大于或等于 %s", err.Field(), err.Param())
	case "lte":
		return fmt.Sprintf("%s 字段的值必须小于或等于 %s", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s 字段的取值范围是 %s", err.Field(), err.Param())
	case "nameValidator":
		return fmt.Sprintf("%s 字段的命名规则是 %s", err.Field(), "名称字符串非空、长度不超过20个字符、并且只能包含字母、数字、汉字或下划线且首字符不能为下划线")
	default:
		return fmt.Sprintf("%s 字段校验失败", err.Field())
	}
}

func TranslateErrorWithDesc(err validator.FieldError, data interface{}) string {
	field := err.StructField()
	description := getFieldDescription(data, field)

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s(%s) 字段是必填的", description, err.Field())
	case "gt":
		return fmt.Sprintf("%s(%s) 字段的值必须大于 %s", description, err.Field(), err.Param())
	case "min":
		return fmt.Sprintf("%s(%s) 字段的最小长度为 %s", description, err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s(%s) 字段的最大长度为 %s", description, err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("%s(%s) 字段必须是有效的电子邮件地址", description, err.Field())
	case "gte":
		return fmt.Sprintf("%s(%s) 字段的值必须大于或等于 %s", description, err.Field(), err.Param())
	case "lte":
		return fmt.Sprintf("%s(%s) 字段的值必须小于或等于 %s", description, err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s(%s) 字段的取值范围是 %s", description, err.Field(), err.Param())
	case "nameValidator":
		return fmt.Sprintf("%s(%s) 字段的命名规则是 %s", description, err.Field(), "名称字符串非空、长度不超过20个字符、并且只能包含字母、数字、汉字或下划线且首字符不能为下划线")
	case "k8sName":
		return fmt.Sprintf("%s(%s) 字段的命名规则是 %s", description, err.Field(), "名称必须以字母开头，必须以字母或者数字或者结尾，长度不超过64个字符，只能包含小写字母、数字、连字符 -且连字符 -不能连续使用")
	case "k8sLabel":
		return fmt.Sprintf("%s(%s) 字段的命名规则是 %s", description, err.Field(), "键值必填、唯一、键长度1-56字符，值长度1-63字符，以字母或数字开头和结尾，中间可以包含字母、数字、减号（-）、下划线（_）、小数点（.）")
	default:
		return fmt.Sprintf("%s 字段校验失败", err.Field())
	}
}

// 获取字段描述
func getFieldDescription(data interface{}, fieldName string) string {
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return fieldName
	}

	tag := field.Tag.Get("description")
	if tag == "" {
		return fieldName
	}

	return tag
}
