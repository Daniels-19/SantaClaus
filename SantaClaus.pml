bool elfAwake =  false;
bool elfLock = false;
bool elfLock2 = false;

int elves = 0;
int helpedElves = 0;

bool reindeerAwake =  false;
bool reindeerLock = false;
bool reindeerLock2 = false;
int reindeer = 0;
int hitchedReindeer = 0;

active proctype Santa ()
{
	Sleep:
		(elfAwake || reindeerAwake)
		goto WakeUp;

	WakeUp:
        if
		:: (reindeer >= 2) ->
			goto GiveToys;
		:: else -> 
			goto HelpElves;
        fi;

	GiveToys:
		reindeerLock = true;
		(hitchedReindeer == 2) // Here we wait for the reindeer to be harnessed
		hitchedReindeer = 0;
		reindeerLock2 = true;
		
		// Give toys

		(hitchedReindeer == 2) //Here we wait for the reindeer to be unharnessed
		hitchedReindeer = 0;
		reindeerLock = false;
		reindeerLock2 = false;
		reindeer = 0;
		
		goto Sleep;
	
	HelpElves:
		elfLock = true;
		(helpedElves == 1) // Wait for all elves to enter
		helpedElves = 0;
		elfLock2 = true;

		// Helping elves

		(helpedElves == 1) // Wait for all elves to want to leave
		helpedElves = 0;
		elfLock = false;
		elfLock2 = false;
		elves = elves - 1;
		
		goto Sleep;

}

active [2] proctype Elves ()
{
	Work:
		skip;

	LeaveWorkbench:
		(elves < 1)
		atomic {
			elves++;
			if
			:: (elves == 1) -> elfAwake = true;
			:: (elves > 1) -> 
				elves--;
				goto LeaveWorkbench;
			fi;
		}
		(elfLock)

	GetHelp:
		elfAwake = false;
		helpedElves++;
		(elfLock2)

        // Santa helping them

		helpedElves++;
		(!elfLock)
		goto Work;
}

active [2] proctype Reindeer ()
{
	Holiday:
		skip;

	HomeFromHoliday:
		(reindeer < 2)
		atomic {
			reindeer++;
			if
			:: (reindeer == 2) -> reindeerAwake = true;
			:: (reindeer > 2) -> 
				reindeer--;
				goto HomeFromHoliday;
			fi
		}
		(reindeerLock)

	GetHitched:
		reindeerAwake = false;
		hitchedReindeer++;		// Here the reindeer are getting harnessed
		(reindeerLock2)

		// Let Santa riding them and delivering presents
		
		hitchedReindeer++;		// Here the reindeer are getting unharnessed
		(!reindeerLock)
		goto Holiday;

}

//We made only one section for Harnessing, giving toys and unharnessing

ltl SantaSleeping {always ((reindeer < 2 && elves < 1) -> Santa@Sleep)}
ltl SantaSleeping2 { always ( (Santa@Sleep) -> eventually (reindeer == 2 || elves == 1))}
ltl SantaGivingToys { always (Santa@GiveToys -> (reindeer == 2 ) ) }
ltl SantaHelpingElves { always (Santa@HelpElves -> (elves == 1 ) ) }

// spin -search SantaClaus.pml
// spin -search -ltl ReindeerSafety SantaClaus.pml
// spin -search -ltl ElfWorking SantaClaus.pml
// spin -search -ltl Sleeping SantaClaus.pml
// spin -p -g -t SantaClaus.pml