`prql-go` is an cgo-free alternative to https://github.com/PRQL/prql/blob/main/prql-lib/README.md 

It uses https://github.com/tetratelabs/wazero as a WASM runtime to compile PRQL to SQL


```
$ go get github.com/pims/prql-go
```

```go
import "github.com/pims/prql-go"
```

```go
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
// SELECT
//   name,
//   age
// FROM
//   employees
```