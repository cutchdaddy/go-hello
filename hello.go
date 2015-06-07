 package main

 import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"sync"
)

//initialize global translation tables


type UserA struct {
	Gender string
	Name struct {
		Title, First, Last string
	}
	Location struct {
		Street, City, State, Zip string
	}
	Email string
	Username string
	Registered string
	Dob string
	Phone string
	Cell string
	Picture struct {
		Large, Medium, Thumbnail string
	}
	Ssn string
}

type ServiceA struct {
	Count int
	Users []UserA

}

type UserB struct {
	Address string
	DateOfBirth string
	Email string
	FullName string
	Gender string
	Phone string
	Username string
}

type ServiceB struct {
	Users []UserB
}

type Move struct {
	Input []string
	Output string
}

func (m *Move) Copy(wg *sync.WaitGroup) {
	m.Output = m.Input[0]
	wg.Done()
}

func (m *Move) DateCopy(wg *sync.WaitGroup) {
	fmt.Println(m.Input[0])
	const formatA = "2006-01-02 15:04:05 -0700"
	const formatB = "Monday January 2, 2006"
	t, _ := time.Parse(formatA, m.Input[0])
	m.Output = t.Format(formatB);
	wg.Done()
}

func (m *Move) Merge(wg *sync.WaitGroup) {
	for i := 0; i < len(m.Input); i++ {
		m.Output = m.Output + m.Input[i]
	}
	wg.Done()
}

func (m *Move) Translate(wg *sync.WaitGroup) {
	genderTable := make(map[string]string)
	genderTable["male"] = "M"
	genderTable["female"] = "F"
	for src, targ := range genderTable {
		if m.Input[0] == src {
			m.Output = targ
		}
	}
	wg.Done()
}

func main()  {

	var wg sync.WaitGroup

	file, e := ioutil.ReadFile("backend.json")
	if e != nil {
 		fmt.Printf("File error: %v\n", e)
 		os.Exit(1)
 	}

 	var data ServiceA

 	if e = json.Unmarshal(file, &data); e != nil {
 		panic(e)
 	}

 	importUsers := data.Users
 	exportUsers := make([]UserB, data.Count)

 	for i := 0; i < len(exportUsers); i++ {

 		currUser := importUsers[i]
 		currLocation := currUser.Location
 		moveAddress := Move{Input: []string{currLocation.Street, "\n", currLocation.City, ", ", currLocation.State, " ", currLocation.Zip}}
 		moveDateOfBirth := Move{Input: []string{currUser.Dob}}
 		moveEmail := Move{Input: []string{currUser.Email}}
 		currFullName := currUser.Name
 		moveFullName := Move{Input: []string{currFullName.Title, " ", currFullName.First, " ", currFullName.Last}}
 		moveGender := Move{Input: []string{currUser.Gender}}
 		movePhone := Move{Input: []string{currUser.Phone}}
 		moveUsername := Move{Input: []string{currUser.Username}}

 		wg.Add(7)

 		go moveAddress.Merge(&wg)
 		go moveDateOfBirth.DateCopy(&wg)
 		go moveEmail.Copy(&wg)
 		go moveFullName.Merge(&wg)
 		go moveGender.Translate(&wg)
 		go movePhone.Copy(&wg)
 		go moveUsername.Copy(&wg)

 		wg.Wait()

 		exportUsers[i].Address = moveAddress.Output
 		exportUsers[i].DateOfBirth = moveDateOfBirth.Output
 		exportUsers[i].Email = moveEmail.Output
 		exportUsers[i].FullName = moveFullName.Output
 		exportUsers[i].Gender = moveGender.Output
 		exportUsers[i].Phone = movePhone.Output
 		exportUsers[i].Username = moveUsername.Output
 	}
 
 	export := ServiceB{Users: exportUsers}
 	exportJson, _ := json.Marshal(export)
 	os.Stdout.Write(exportJson)

 	if e = ioutil.WriteFile("serviceB.json", exportJson, 0644); e != nil {
 		fmt.Printf("File error: %v\n", e)
 		os.Exit(1)
 	}

}