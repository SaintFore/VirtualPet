package main

import (
	"fmt"

	"cyberpet/creature"
)

func main() {
	myPet, err := creature.Load()
	if err != nil {
		myPet = creature.NewPet("brian")
	}
	myPet.Life()

	for {
		var cmd string
		fmt.Print(">>")
		fmt.Scanln(&cmd)
		switch cmd {
		case "feed":
			myPet.Feed()
		case "play":
			myPet.Play()
		case "status":
			myPet.PrintStatus()
		case "exit", "quit":
			if err := myPet.Save(); err != nil {
				fmt.Println("存储失败 ", err)
			} else {
				fmt.Println("存储成功")
			}
			fmt.Println("再见")
			return
		default:
			fmt.Println("输入参数错误")
		}
	}
}
