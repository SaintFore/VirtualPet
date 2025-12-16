// Package creature 提供虚拟宠物核心逻辑
package creature

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	MaxHunger  = 100
	MaxEnergy  = 100
	HungerInc  = 1  // 每次自动增加的饥饿
	EnergyInc  = 1  // 每次自动恢复的精力
	FeedVal    = 20 // 喂食减少的饥饿
	PlayCost   = 10 // 玩耍消耗的精力
	PlayHunger = 10 // 玩耍增加的饥饿
)

type LivingBeing interface {
	Feed()
	Play()
	PrintStatus()
	Save() error
	Life()
	GetState() ([]byte, error)
}

type Pet struct {
	Name      string        `json:"name"`
	Hungry    int           `json:"hungry"`
	Energy    int           `json:"energy"`
	Alive     bool          `json:"alive"`
	DeathChan chan struct{} `json:"-"`
	QuitChan  chan struct{} `json:"-"`
	mu        sync.Mutex    `json:"-"`
}

func (pet *Pet) GetState() ([]byte, error) {
	pet.mu.Lock()
	defer pet.mu.Unlock()
	return json.Marshal(pet)
}

func NewPet(name string) *Pet {
	return &Pet{
		Name:      name,
		Hungry:    0,
		Energy:    MaxEnergy,
		Alive:     true,
		DeathChan: make(chan struct{}),
		QuitChan:  make(chan struct{}),
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
	pet.DeathChan = make(chan struct{})
	pet.QuitChan = make(chan struct{})
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
	if pet.Hungry -= FeedVal; pet.Hungry < 0 {
		pet.Hungry = 0
	}
	fmt.Printf("%v正在进食\n", pet.Name)
}

func (pet *Pet) Play() {
	pet.mu.Lock()
	defer pet.mu.Unlock()
	if pet.Energy < PlayCost {
		fmt.Println("太累了")
	} else if pet.Hungry+PlayHunger > MaxHunger {
		fmt.Println("太饿了")
	} else {
		pet.Energy -= PlayCost
		pet.Hungry += PlayHunger
		fmt.Printf("%v玩得很开心\n", pet.Name)
	}
}

func (pet *Pet) Life() {
	pet.mu.Lock()
	pet.QuitChan = make(chan struct{})
	pet.mu.Unlock()
	ticker := time.NewTicker(time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				pet.mu.Lock()
				if pet.Energy += EnergyInc; pet.Energy > MaxEnergy {
					pet.Energy = MaxEnergy
				}
				if pet.Hungry += HungerInc; pet.Hungry >= MaxHunger {
					pet.Alive = false
					close(pet.DeathChan)
				}
				if !pet.Alive {
					fmt.Println("\n宠物已经死了")
					pet.mu.Unlock()
					return
				}
				pet.mu.Unlock()
			case <-pet.QuitChan:
				fmt.Println("时间暂停了")
				return
			}
		}
	}()
}

func (pet *Pet) StopLife() {
	close(pet.QuitChan)
}

func (pet *Pet) StartLife() {
	pet.Life()
}
