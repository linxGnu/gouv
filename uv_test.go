package binduv

import (
	"fmt"
	"testing"
)

func TestGetUVVer(t *testing.T) {
	fmt.Println(GetLibuvVersion())
	fmt.Println(GetLibuvVersionString())
	fmt.Println(UvErrName(1))
}
