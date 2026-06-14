package main

import (
	"fmt"
	"os"
)

func main() {
	files := []string{
		"wails.json",
		"go.mod",
		"main.go",
		"app.go",
		"frontend/index.html",
	}
	allExist := true
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fmt.Printf("❌ Thiếu file quan trọng: %s\n", f)
			allExist = false
		} else {
			fmt.Printf("✅ Đã tìm thấy: %s\n", f)
		}
	}
	if allExist {
		fmt.Println("🚀 Môi trường Wails Phase 01 đã sẵn sàng hoàn hảo!")
	} else {
		fmt.Println("⚠️ Vui lòng hoàn thành các tác vụ còn thiếu.")
		os.Exit(1)
	}
}
