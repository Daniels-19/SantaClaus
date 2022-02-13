bool santaLock =  false;
bool reindeerLock = false;
bool elfLock = false;
int reindeer = 0;
int hitchedReindeer = 0;
int elves = 0;
int helpedElves = 0;

proctype Santa ()
{
	Remainder:
		skip;

	Entry:
		(santaLock)
        if
          :: (reindeer >= 9) ->
                goto PrepareSleigh;
          :: else -> 
                goto HelpElves;
        fi;

	PrepareSleigh:
		// Prepare sleigh
		goto PostPrepareSleigh;

	PostPrepareSleigh:
		reindeerLock = true;
		(hitchedReindeer == 9)
		atomic {
			reindeerLock = false;
			reindeer = reindeer - 9;
			hitchedReindeer = 0;
		}
		goto Exit;
	
	HelpElves:
		// Helping elves
		goto PostHelpElves;

	PostHelpElves:
		elfLock = true;
		(helpedElves == 3)
		atomic {
			elfLock = false;
			elves = elves - 3;
			helpedElves = 0;
		}
		goto Exit;

	Exit:
		goto Remainder;
}

proctype Elves ()
{
	Remainder:
		skip;

	Entry:
		(elves < 3)
		atomic {
			elves++;
			if
			:: (elves == 3) -> santaLock = true;
			:: (elves > 3) -> goto Entry;
			fi;
		}
		(elfLock)

	GetHelp:
        // Santa helping them
		skip;
	
	Exit:
		helpedElves++;
		goto Remainder;
}

proctype Reindeer ()
{
	Remainder:
		skip;

	Entry:
		(reindeer < 9)
		atomic {
			reindeer++;
			if
			:: (reindeer == 9) -> santaLock = true;
			:: (reindeer > 9) -> goto Entry;
			fi
		}
		(reindeerLock)

	GetHitched:
		// Let Santa riding them and delivering presents
		skip;
	
	Exit:
		hitchedReindeer++;
		goto Remainder;
}

init {
	run Santa();
	//run Reindeer(); run Reindeer(); run Reindeer(); run Reindeer(); run Reindeer(); run Reindeer(); run Reindeer(); run Reindeer(); run Reindeer(); 
	run Elves(); run Elves(); run Elves();
}

ltl ReindeerSafety {always (Reindeer@GetHitched -> always (!Elves@GetHelp && !Santa@PrepareSleigh && !Santa@HelpElves && Santa@PostPrepareSleigh && !Santa@PostHelpElves))}
ltl ElfSafety {always (Elves@GetHelp -> always (!Reindeer@GetHitched && !Santa@HelpElves && !Santa@PostPrepareSleigh && Santa@PostHelpElves))}
ltl Sleeping {always (!Santa@PrepareSleigh && !Santa@HelpElves) -> always (!Reindeer@GetHitched && !Elves@GetHelp)}

// spin -search SantaClaus.pml
// spin -search -ltl ReindeerSafety SantaClaus.pml
// spin -search -ltl ElfSafety SantaClaus.pml
// spin -search -ltl Sleeping SantaClaus.pml
// spin -p -g -t SantaClaus.pml