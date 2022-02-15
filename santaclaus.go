package main

import "fmt"

var santaWake chan struct{} = make(chan struct{})
var santaCS chan struct{} = make(chan struct{})
var santaExit chan struct{} = make(chan struct{})

var showInLock chan struct{} = make(chan struct{})
var helpRDLock chan struct{} = make(chan struct{})
var showOutLock chan struct{} = make(chan struct{})

var elves int = 0
var elvesShownIn int = 0
var elvesShownOut int = 0

var elfGroup int = 3

var hitchLock chan struct{} = make(chan struct{})
var toyLock chan struct{} = make(chan struct{})
var unhitchLock chan struct{} = make(chan struct{})

var deer int = 0
var deerHitched int = 0
var deerUnhitched int = 0

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

			<-santaCS // Wait for all deer to get hitched

			SantaSend(toyLock, deerGroup)

			// Give toys

			<-santaExit // Wait for all deer to be unhitched

			deer -= deerGroup
			deerHitched = 0
			deerUnhitched = 0
			SantaSend(unhitchLock, deerGroup)

		} else { // Elf case
			SantaSend(showInLock, elfGroup)

			<-santaCS // Wait for elves to enter study

			SantaSend(helpRDLock, elfGroup)

			// Help elves

			<-santaExit // Wait for elves to leave study

			elves -= elfGroup
			elvesShownIn = 0
			elvesShownOut = 0
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
		if deerHitched == deerGroup {
			santaCS <- struct{}{}
		}
		<-toyLock // Wait for santa to hitch

		// Deliver presents

		deerUnhitched++
		if deerUnhitched == deerGroup {
			santaExit <- struct{}{}
		}
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

		elvesShownIn++
		if elvesShownIn == elfGroup {
			santaCS <- struct{}{}
		}
		<-helpRDLock // Wait for santa to show into study

		// Receive help

		elvesShownOut++
		if elvesShownOut == elfGroup {
			santaExit <- struct{}{}
		}
		<-showOutLock // Wait for santa to show out of study
	}
}
