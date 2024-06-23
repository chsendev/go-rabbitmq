package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

type User struct {
	Msg
	Name string `json:"name"`
}

func TestName(t *testing.T) {
	f := func(u *User) {}

	to := reflect.TypeOf(f)
	fmt.Println("Parameters:")
	str := "{\"name\":\"abc\"}"
	for i := 0; i < to.NumIn(); i++ {
		t := to.In(i)
		fmt.Println(t.Name())
		//json.Unmarshal()
		//fmt.Println(t.Elem())
		// 创建一个该类型的新实例
		instance := reflect.New(t)
		err := json.Unmarshal([]byte(str), instance.Interface())
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Parameter %d type: %s, value: %v\n", i, t, instance.Elem())
	}

	abc(func(t Msg) {

	})
}

func abc(f Fun) {

}

type Msg struct {
}

type Fun func(t Msg)

//func TestName3(t *testing.T) {
//	re()
//}

type Handel[T any] interface {
	handel(t T)
}

func re() {
	// 得到一个[]byte数组
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
