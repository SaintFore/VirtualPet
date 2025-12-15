package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type LivingBeing interface {
	Feed()
	Play()
	PrintStatus()
	Save() error
}

type Pet struct {
	Name   string `json:"name"`
	Hungry int    `json:"hungry"`
	Energy int    `json:"energy"`
	Alive  bool   `json:"alive"`
	mu     sync.Mutex
}

func (pet *Pet) Save() error {
	pet.mu.Lock()
	defer pet.mu.Unlock()

	data, err := json.MarshalIndent(pet, "", "  ") // indent两个字符，符合json格式化
	if err != nil {
		return err
	}

	return os.WriteFile("pet_data.json", data, 0644)
}

func Load() (*Pet, error) {
	data, err := os.ReadFile("pet_data.json")
	if err != nil {
		return nil, err
	}
	var pet Pet
	if err = json.Unmarshal(data, &pet); err != nil {
		return nil, err
	}
	return &pet, nil
}

func (pet *Pet) PrintStatus() {
	pet.mu.Lock()
	defer pet.mu.Unlock()
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
	myPet, err := Load()
	if err != nil {
		myPet = &Pet{Name: "brian", Hungry: 0, Energy: 100, Alive: true}
	}
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
			time.Sleep(time.Second * 1)
			myPet.mu.Lock()
			if !myPet.Alive {
				fmt.Println("\n宠物已经死了")
				os.Exit(0)
			}
			myPet.mu.Unlock()
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
