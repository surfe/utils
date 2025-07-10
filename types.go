package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go"
	"github.com/surfe/logger/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PtrTo[T any](v T) *T {
	return &v
}

func PBool(x bool) *bool {
	return &x
}

func PString(x string) *string {
	return &x
}

func PInt(x int) *int {
	return &x
}

func PTime(x time.Time) *time.Time {
	return &x
}

func Bool(x *bool) bool {
	if x == nil {
		return false
	}
	return *x
}

func String(x *string) string {
	if x == nil {
		return ""
	}
	return *x
}

func Int64(x *int64) int64 {
	if x == nil {
		return 0
	}
	return *x
}

func Int(x *int) int {
	if x == nil {
		return 0
	}
	return *x
}

func GetSafeType[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}

	return *ptr
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Atof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func IsEmptyType[T comparable](t T) bool {
	var empty T
	return t == empty
}

func GetPointerOrNil[T comparable](t T) *T {
	var zero T
	if t == zero {
		return nil
	}

	return &t
}

func UnsafeObjectIDHex(s string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}

func SafeString(x interface{}) string {
	if x == nil {
		return ""
	}

	switch reflect.TypeOf(x).Kind() {
	case reflect.Ptr:
		if reflect.ValueOf(x).IsNil() {
			return ""
		}
		return fmt.Sprint(reflect.Indirect(reflect.ValueOf(x)))
	case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		if reflect.ValueOf(x).IsNil() {
			return ""
		}
	default:
		return fmt.Sprint(x)
	}

	return fmt.Sprint(x)
}

func GenerateCode(n int) string {
	var letterRunes = []rune("0123456789")
	var maxVal = big.NewInt(int64(len(letterRunes)))

	b := make([]rune, n)

	for i := range b {
		bn, _ := rand.Int(rand.Reader, maxVal)
		b[i] = letterRunes[bn.Int64()]
	}
	return string(b)
}

func ToStripeError(err error) (e stripe.Error) {
	_ = json.Unmarshal([]byte(err.Error()), &e)
	return
}

func Transfer(from, to interface{}) {
	x, _ := json.Marshal(from)
	err := json.Unmarshal(x, to)
	if err != nil {
		logger.Log(context.Background()).Err(err).Error("Could not unmarshal")
	}
}

func FindLooseStringInSlice(slice []string, item string) int {
	return slices.IndexFunc(slice, func(s string) bool {
		return EqualsLooseString(s, item)
	})
}

// EqualsLooseString compares two strings after some clean-up, case in-sensitively
func EqualsLooseString(s, t string) bool {
	return strings.EqualFold(LooseString(s), LooseString(t))
}

func GetType(value any) string {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		return strings.ToUpper(t.Elem().Name())
	}

	return strings.ToUpper(t.Name())
}

func GetFloat(value any) (float64, error) {
	return strconv.ParseFloat(fmt.Sprint(value), 64)
}

func GetStringSlice(f interface{}) []string {
	switch f := f.(type) {
	case string:
		return []string{f}
	case []string:
		return f
	case []interface{}:
		var res []string
		for _, val := range f {
			res = append(res, fmt.Sprint(val))
		}
		return res
	default:
		return []string{SafeString(f)}
	}
}

func UUIDFromString(str string) uuid.UUID {
	str = strings.ReplaceAll(str, "-", "")
	buf, _ := hex.DecodeString(str)
	uuid := [16]byte{}
	copy(uuid[:], buf)
	return uuid
}
