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
		(hitchedReindeer == 2)
		hitchedReindeer = 0;
		reindeerLock2 = true;

		// Give toys

		(hitchedReindeer == 2)
		hitchedReindeer = 0;
		reindeerLock = false;
		reindeerLock2 = false;
		reindeer = reindeer - 2;
		
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
		hitchedReindeer++;
		(reindeerLock2)

		// Let Santa riding them and delivering presents

		hitchedReindeer++;
		(!reindeerLock)
		goto Holiday;

}

ltl SantaSleeping {always ((reindeer < 2 && elves < 1) -> Santa@Sleep)}

// ltl SantaSleeping {always (reindeer < 2 && elves < 1) -> always (Santa@Sleep)}

ltl SantaGivingToys {always (reindeer == 9 -> eventually (Santa@GiveToys)) }
ltl SantaHelpingElves {always (elves == 2 -> eventually (Santa@HelpElves))}
ltl ElfHelp {always (Santa@HelpElves -> always (Reindeer@HomeFromHoliday || Reindeer@Holiday))}
ltl ReindeerHarnessed {always (Santa@GiveToys) -> always (Elves@LeaveWorkbench || Elves@LeaveWorkbench)}

// ltl SantaHelpingElves {always (Santa@HelpElves) -> always (!Reindeer@GetHitched && elves == 3)}
// ltl ReindeerHoliday {always (!Reindeer@GetHitched -> always (!Santa@PrepareSleigh))}
// ltl ReindeerHitched {always (Reindeer@GetHitched -> always (!Elves@GetHelp && !Santa@PrepareSleigh && !Santa@HelpElves && Santa@PostPrepareSleigh && !Santa@PostHelpElves))}
// ltl ElfWorking {always (!Elves@GetHelp -> always (!Santa@ShowOut))}
// ltl ElfGettingHelp {always (Elves@GetHelp -> always (Santa@ShowOut))}

// spin -search SantaClaus.pml
// spin -search -ltl ReindeerSafety SantaClaus.pml
// spin -search -ltl ElfWorking SantaClaus.pml
// spin -search -ltl Sleeping SantaClaus.pml
// spin -p -g -t SantaClaus.pml