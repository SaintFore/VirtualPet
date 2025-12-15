package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Pet struct {
	Name   string
	Hungry int
	Energy int
	Alive  bool
	mu     sync.Mutex
}

func (pet *Pet) PrintStatus() {
	fmt.Println("==========")
	fmt.Printf("%s的饥饿度为%d,能量为%d\n", pet.Name, pet.Hungry, pet.Energy)
	fmt.Println("==========")
}

func (pet *Pet) Feed() {
	pet.mu.Lock()
	defer pet.mu.Unlock()
	if pet.Hungry -= 20; pet.Hungry < 0 {
		pet.Hungry = 0
	}
	fmt.Printf("%v正在进食\n", pet.Name)
	fmt.Println("当前状态")
	pet.PrintStatus()
}

func (pet *Pet) Play() {
	pet.mu.Lock()
	defer pet.mu.Unlock()
	if pet.Energy < 10 {
		fmt.Println("太累了")
	} else if pet.Hungry+10 > 100 {
		fmt.Println("太饿了")
	} else {
		pet.Energy -= 10
		pet.Hungry += 10
		fmt.Printf("%v玩得很开心\n", pet.Name)
	}
}

func (pet *Pet) Die() {
	pet.mu.Lock()
	defer pet.mu.Unlock()
	pet.Alive = false
}

func main() {
	myPet := Pet{Name: "brian", Hungry: 90, Energy: 10, Alive: true}
	go func() {
		for {
			time.Sleep(time.Second * 5)
			myPet.mu.Lock()
			myPet.Hungry += 5
			myPet.Energy += 5
			myPet.mu.Unlock()
			if myPet.Hungry >= 100 {
				myPet.Die()
			}
		}
	}()

	go func() {
		for {
			if !myPet.Alive {
				time.Sleep(time.Second * 1)
				fmt.Println("\n宠物已经死了")
				os.Exit(0)
			}
		}
	}()

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
			fmt.Println("再见")
			return
		default:
			fmt.Println("输入参数错误")
		}

	}
}
