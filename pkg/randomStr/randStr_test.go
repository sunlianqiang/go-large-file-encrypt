package randomStr

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	for i := 0; i < 10; i++ {

		fmt.Println(GetRandomStr(32))

		// fmt.Println(getToken(32))
	}
}
