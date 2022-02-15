package main

import "fmt"

var santaWake chan struct{} = make(chan struct{})

var showInLock chan struct{} = make(chan struct{})
var helpRDLock chan struct{} = make(chan struct{})
var showOutLock chan struct{} = make(chan struct{})

var elves int = 0
var elvesHelped int = 0

var elfGroup int = 3

var hitchLock chan struct{} = make(chan struct{})
var toyLock chan struct{} = make(chan struct{})
var unhitchLock chan struct{} = make(chan struct{})

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

func SantaSend(c chan struct{}, no int) {
	for i := 0; i < no; i++ {
		c <- struct{}{}
	}
}

func Santa() {
	for true {
		<-santaWake

		if deer >= deerGroup { // Deer priority
			fmt.Println("Helping deer")

			SantaSend(hitchLock, deerGroup)

			for deerHitched < deerGroup {
				// Wait for all deer to get hitched
			}
			deerHitched = 0

			SantaSend(toyLock, deerGroup)

			// Give toys

			for deerHitched < deerGroup {
				// Wait for all deer to be unhitched
			}

			deerHitched = 0
			deer -= deerGroup
			SantaSend(unhitchLock, deerGroup)

		} else { // Elf case
			SantaSend(showInLock, elfGroup)

			for elvesHelped < elfGroup {
				// Wait for elves to enter study
			}

			elvesHelped = 0
			SantaSend(helpRDLock, elfGroup)

			// Help elves

			for elvesHelped < elfGroup {
				// Wait for elves to leave study
			}

			elvesHelped = 0
			elves -= elfGroup
			SantaSend(showOutLock, elfGroup)
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

		<-hitchLock // Wait for santa to wake up

		deerHitched++
		<-toyLock // Wait for santa to hitch

		// Deliver presents

		deerHitched++
		<-unhitchLock // Wait for santa to unhitch
	}
}

func Elf() {
	for true {
		// Working

		elves++
		if elves == elfGroup {
			santaWake <- struct{}{}
		}

		<-showInLock // Wait for santa to wake up

		elvesHelped++

		<-helpRDLock // Wait for santa to show into study

		// Receive help

		elvesHelped++
		<-showOutLock // Wait for santa to show out of study
	}
}
