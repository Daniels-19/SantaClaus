package main

import "fmt"

var santaWake chan struct{} = make(chan struct{})

var elfLock chan struct{} = make(chan struct{})
var elfLock2 chan struct{} = make(chan struct{})
var elfLock3 chan struct{} = make(chan struct{})

var elves int = 0
var elvesHelped int = 0

var elfGroup int = 3

var deerLock chan struct{} = make(chan struct{})
var deerLock2 chan struct{} = make(chan struct{})
var deerLock3 chan struct{} = make(chan struct{})

var deer int = 0
var deerHitched int = 0

var deerGroup int = 9

func main() {
	go Santa()

	go Deer()
	go Deer()
	go Deer()
	go Deer()
	go Deer()
	go Deer()
	go Deer()
	go Deer()
	go Deer()

	go Elf()
	go Elf()
	go Elf()
	go Elf()
	go Elf()
	go Elf()
	go Elf()
	go Elf()
	go Elf()
	go Elf()
}

func Santa() {
	for true {
		<-santaWake

		if deer >= deerGroup { // Give toys
			fmt.Println("Helping deer")

			deerLock <- struct{}{}

			for deerHitched < deerGroup {
				// Wait for all deer to get hitched
			}
			deerHitched = 0

			deerLock2 <- struct{}{}

			// Give toys

			for deerHitched < deerGroup {
				// Wait for all deer to be unhitched
			}

			deerHitched = 0
			deer -= deerGroup
			deerLock3 <- struct{}{}

		} else { // Help elves
			elfLock <- struct{}{}
			for elvesHelped < elfGroup {
				// Wait for elves to enter study
			}
			elvesHelped = 0
			elfLock2 <- struct{}{}

			// Help elves

			for elvesHelped < elfGroup {
				// Wait for elves to leave study
			}

			elvesHelped = 0
			elves -= elfGroup
			elfLock3 <- struct{}{}
		}
	}
}

func Deer() {
	for true {
		// On holiday

		deer++
		if deer == deerGroup {
			santaWake <- struct{}{}
		}

		<-deerLock // Wait for santa to wake up

		deerHitched++
		<-deerLock2 // Wait for santa to hitch

		// Deliver presents

		deerHitched++
		<-deerLock3 // Wait for santa to unhitch
	}
}

func Elf() {
	for true {
		// Working

		elves++
		if elves == elfGroup {
			santaWake <- struct{}{}
		}

		<-elfLock // Wait for santa to wake up

		elvesHelped++

		<-elfLock2 // Wait for santa to show into study

		// Receive help

		elvesHelped++
		<-elfLock3 // Wait for santa to show out of study
	}
}

/*

var Bools map[string]map[string]bool = map[string]map[string]bool{
	"Elf": {
		"Awake": false,
		"Lock1": false,
		"Lock2": false,
	},
	"Deer": {
		"Awake": false,
		"Lock1": false,
		"Lock2": false,
	},
}

var Ints map[string]map[string]int = map[string]map[string]int{
	"Elf": {
		"Number1": 0,
		"Number2": 0,
		"Group":   2,
	},
	"Deer": {
		"Number1": 0,
		"Number2": 0,
		"Group":   2,
	},
}

func Worker(key string) {
	for true {
		Ints[key]["Number1"]++

		if Ints[key]["Number1"] == Ints[key]["Group"] {
			Bools[key]["Awake"] = true
		}
		for !Bools[key]["Lock1"] {
			// Wait for santa to wake up
		}
		Bools[key]["Awake"] = false
		Ints[key]["Number2"]++
		for !Bools[key]["Lock2"] {
			// Wait for santa to enter critical section
		}

		// Critical section

		Ints[key]["Number2"]++
		for !Bools[key]["Lock1"] {
			// Wait for santa to exit critical section
		}
	}
}
*/
