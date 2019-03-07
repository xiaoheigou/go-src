content := "This Is LimitReader Example"
reader := strings.NewReader(content)
limitReader := &io.LimitedReader{R: reader, N: 8}
for limitReader.N > 0 {
    tmp := make([]byte, 2)
    limitReader.Read(tmp)
    fmt.Printf("%s", tmp)
}
// out: this is
//1.11. PipeReader 和 PipeWriter 类型