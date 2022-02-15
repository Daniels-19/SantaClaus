package main

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

var elfMutex chan struct{} = make(chan struct{})

var hitchLock chan struct{} = make(chan struct{})
var toyLock chan struct{} = make(chan struct{})
var unhitchLock chan struct{} = make(chan struct{})

var deer int = 0
var deerHitched int = 0
var deerUnhitched int = 0

var deerGroup int = 9

var deerMutex chan struct{} = make(chan struct{})

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

func SantaReceive(c chan struct{}, no int) {
	for i := 0; i < no; i++ {
		<-c
	}
}

func Santa() {
	deerMutex <- struct{}{}
	elfMutex <- struct{}{}

	for true {
		<-santaWake

		if deer >= deerGroup { // Deer priority

			SantaSend(hitchLock, deerGroup) // Allow deer to be hitched

			SantaReceive(santaCS, deerGroup) // Wait for all deer to get hitched

			SantaSend(toyLock, deerGroup) // Confirm all deer have entered

			// Give toys

			SantaReceive(santaExit, deerGroup) // Wait for all deer to be unhitched

			deer -= deerGroup
			deerHitched = 0
			deerUnhitched = 0
			deerMutex <- struct{}{}

		} else { // Elf case

			SantaSend(showInLock, elfGroup) // Allow elves into study

			SantaReceive(santaCS, elfGroup) // Wait for elves to enter study

			SantaSend(helpRDLock, elfGroup) // Confirm all elves have entered

			// Help elves

			SantaReceive(santaExit, elfGroup) // Wait for elves to leave study

			elves -= elfGroup
			elvesShownIn = 0
			elvesShownOut = 0
			elfMutex <- struct{}{}
		}
	}
}

func Deer() {
	for true {
		// On holiday

		<-deerMutex
		deer++
		if deer == deerGroup {
			santaWake <- struct{}{}
		} else {
			deerMutex <- struct{}{}
		}

		<-hitchLock // Wait for santa to wake up

		// Get hitched

		santaCS <- struct{}{} // Signal to santa that I am hitched

		<-toyLock // Wait for all deer to be hitched

		// Deliver presents

		// Get unhitched

		santaExit <- struct{}{} // Signal to santa that I am unhitched
	}
}

func Elf() {
	for true {
		// Working

		<-elfMutex
		elves++
		if elves == elfGroup {
			santaWake <- struct{}{}
		} else {
			elfMutex <- struct{}{}
		}

		<-showInLock // Wait for santa to wake up and open study

		// Enter study

		santaCS <- struct{}{} // Signal to santa that I am in the study

		<-helpRDLock // Wait for santa other elves to enter study

		// Receive help

		// Leave study

		santaExit <- struct{}{} // Signal to santa that I left the study
	}
}
