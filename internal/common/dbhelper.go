package common

import "fmt"

// Placeholders tạo chuỗi @p1, @p2, ... @pn dùng cho SQL Server (go-mssqldb).
func Placeholders(n int) string {
	s := ""
	for i := 1; i <= n; i++ {
		if i > 1 {
			s += ", "
		}
		s += fmt.Sprintf("@p%d", i)
	}
	return s
}

// ArgsToSlice chuyển các giá trị thành slice (dùng cho QueryRow/Exec với mssql).
func ArgsToSlice(args ...any) []any {
	out := make([]any, len(args))
	for i, a := range args {
		out[i] = a
	}
	return out
}
