package test 

import "testing"

func TestSqrt(t *testing.T) {
    v := Sqrt(16)
    if v != 4 {
        t.Errorf("Sqrt(16) failed, Got %v, expectd 4.", v)
    }
}
