package main

import (
	"fmt"
	"net/http"
	"os"

	"cyberpet/creature"
)

// var myPet creature.LivingBeing
// var (
//
//	myPet creature.Pet
//	err   error
//
// )
var myPet *creature.Pet

func main() {
	var err error
	myPet, err = creature.Load()
	if err != nil {
		myPet = creature.NewPet("brian")
	}
	myPet.Life()

	go func() {
		<-myPet.DeathChan
		if err := myPet.Save(); err != nil {
			fmt.Println("存储失败", err)
			os.Exit(1)
		} else {
			fmt.Println("存储成功")
			os.Exit(0)
		}
	}()

	fmt.Println("Server is running on http://localhost:18080")
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/feed", onlyPost(feedHandler))
	http.HandleFunc("/play", onlyPost(playHandler))
	http.HandleFunc("/stop", onlyPost(stopHandler))
	http.HandleFunc("/start", onlyPost(startHandler))
	go func() {
		fs := http.FileServer(http.Dir("./static/"))
		http.Handle("/", fs)
		httpErr := http.ListenAndServe(":18080", nil)
		if httpErr != nil {
			fmt.Println("Server failed:", httpErr)
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

func statusHandler(w http.ResponseWriter, r *http.Request) {
	data, err := myPet.GetState()
	if err != nil {
		fmt.Println("获取数据失败")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	myPet.Feed()
	fmt.Fprintf(w, "Feed success")
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	myPet.Play()
	fmt.Fprintf(w, "Play success")
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	myPet.StopLife()
	fmt.Fprintf(w, "Stop success")
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	myPet.StartLife()
	fmt.Fprintf(w, "Start success")
}

func onlyPost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
