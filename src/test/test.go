package  main 

import (
        "math"
        "fmt"
        "bytes"
       )

func Sqrt(i int) int {
    v := math.Sqrt(float64(i))
    return int(v)
}

func toBase10(num string) int64 {
    str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    base := make(map[byte] int, 36) 
    for i := 0; i < len(str); i ++ {
        base[str[i]] = i
    }
    var rs int64 = 0
    for i := 0; i < len(num); i ++ {
        idx, _ := base[num[i]]
        pow := len(num) - i -1
        tmp := math.Pow(36,float64(pow))
        fmt.Println(num,",str(i):",num[i],",idx:",idx,",tmp:", tmp, ",pow:", pow)
        rs += int64(idx) * int64(tmp) 
    }
    return rs
}
func toBase36 (num int64) string {
    str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    base := []byte(str)
    var rs []byte;
    var init int64 = 1234567890
    num += init
    mod := num % 36 
    div := (num - mod) / 36
    rs = append(rs, base[mod])
    for {
        if div == 0 {
            break
        }
        mod = div % 36
        div = (div - mod) / 36
        rs = append(rs, base[mod])
    }
    
    for i, j := 0, len(rs) - 1; i < j; i, j = i + 1, j -1 {
        rs[i], rs[j] = rs[j], rs[i]
    }
    return string(rs)
}

func main() {
    var bf bytes.Buffer
    bf.WriteString(fmt.Sprintf("%s", "test"))
    bf.WriteString(fmt.Sprintf("%s", "-test"))

    fmt.Println(bf.String())
}
