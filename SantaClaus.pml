bool lock =  false;
bool reindeerLock = false;
bool elveLock = false;
int reindeer = 0;
int elves = 0;

active proctype Santa ()
{
	Remainder:
		skip;

	Entry:
		(lock = true)
        if
          :: (reindeer == 9) -> 
                reindeerLock = true;
                reindeer = 0;
                GivingGifts:
                    skip;
          :: else -> 
                reindeerLock = true;
                elves = 0;
                HelpingElves:
                    skip;
        fi;
	Dog:
		// Dog in Yard
		skip;
	
	Exit:
		
		goto Remainder;
}

active proctype Elves ()
{
	Remainder:
		skip;

	Entry:

	Study:
        // Santa helping them
		skip;
	
	Exit:
		
		goto Remainder;
}

active proctype Reindeer ()
{
	Remainder:
		skip;

	Entry:

	Harnessed:
		// Santa riding them and delivering presents
		skip;
	
	Exit:
		
		goto Remainder;
}

ltl ReindeerSafety {always (Reindeer@Harnessed -> always (!Elves@Study && Santa@GivingGifts && !Santa@HelpingElves))}
ltl ReindeerSafety {always (Elves@Study -> always (!Reindeer@Harnessed && Santa@GivingGifts && !Santa@HelpingElves))}
ltl Nothing {always (!Santa@GivingGifts && !Santa@HelpingElves) -> always (!Reindeer@Harnessed && !Elves@Study)}