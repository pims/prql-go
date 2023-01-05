package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pims/prql-go"
)

func main() {
	ctx := context.Background()
	w, err := prql.New(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer w.Close(ctx)

	sql, err := w.Compile(ctx, "from employees | select [name,age]")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(sql)
}
