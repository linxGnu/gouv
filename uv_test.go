package gouv

import (
	"fmt"
	"testing"
)

func TestGetUVVer(t *testing.T) {
	fmt.Println(GetLibuvVersion())
	fmt.Println(GetLibuvVersionString())
	fmt.Println(UvErrName(1))
	fmt.Println(GetFreeMemory())
	fmt.Println(GetTotalMemory())
}
