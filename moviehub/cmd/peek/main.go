package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// peek เปิดอ่านตั๋วโดยไม่ใช้กุญแจลับเลยสักตัว
func main() {
	if len(os.Args) < 2 {
		fmt.Println("วิธีใช้: go run ./cmd/peek $TOKEN")
		os.Exit(1)
	}
	parts := strings.Split(os.Args[1], ".")
	if len(parts) != 3 {
		fmt.Println("นี่ไม่ใช่ตั๋ว JWT (ต้องมี 3 ท่อน) แต่ได้", len(parts), "ท่อน")
		os.Exit(1)
	}
	fmt.Println("ตั๋วมี", len(parts), "ท่อน")
	for i, name := range []string{"ส่วนหัว", "เนื้อความ"} {
		raw, err := base64.RawURLEncoding.DecodeString(parts[i])
		if err != nil {
			fmt.Println(name, "อ่านไม่ออก", err)
			continue
		}
		fmt.Printf("%s %s\n", name, raw)
	}
	fmt.Println("ลายเซ็น", parts[2], "ยาว", len(parts[2]), "ตัว")
}
