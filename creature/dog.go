package creature

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
	Life()
}

type Pet struct {
	Name   string `json:"name"`
	Hungry int    `json:"hungry"`
	Energy int    `json:"energy"`
	Alive  bool   `json:"alive"`
	mu     sync.Mutex
}

func NewPet(name string) *Pet {
	return &Pet{
		Name:   name,
		Hungry: 0,
		Energy: 100,
		Alive:  true,
	}
}

func (pet *Pet) Save() error {
	pet.mu.Lock()
	defer pet.mu.Unlock()

	data, err := json.MarshalIndent(pet, "", "  ") // indent两个字符，符合json格式化
	if err != nil {
		return err
	}

	return os.WriteFile("pet_data.json", data, 0o644)
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

func (pet *Pet) Life() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			pet.mu.Lock()
			pet.Hungry += 5
			pet.Energy += 5
			if pet.Hungry >= 100 {
				pet.Alive = false
			}
			if !pet.Alive {
				fmt.Println("\n宠物已经死了")
				os.Exit(0)
			}
			pet.mu.Unlock()
		}
	}()
}
