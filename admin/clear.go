package admin

import (
	"fmt"
	"os"
	pkg "regist/pkg"
)

func Clear() {
	directory := "regist/templates/excel"
	
	err := os.RemoveAll(directory)
	pkg.ForError(err)

	fmt.Println("All files deleted from derectory")
	
} 