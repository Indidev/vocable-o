# vocable-o
A terminal based vocable training application written in GO

### Installation and Dependencies ###

Vocable-o depends on termbox-go (https://github.com/nsf/termbox-go)

Install and update with ```go get -u github.com/indidev/vocable-o```

### Features ###

#### Language support / learning ####
Vocable-o supports every language you like, however you have to add all languages yourself (this can be done in the application).
Five pockets are used to maximize your learning efficiency.
Every vocable moves to the next pocket after it was guessed correctly x times, where x = pocketnumber + 1 (x can be modified in a later version).
If a word is guessed false once and the success-counter is greater then 0, the success-counter of that word is set to 0, otherwise the word is moved back to the very first pocket.

#### Character replacements ####
Characters which are not present on your keyboard(or pretty inconvenient) can be substituted by other keys/key-sequences by defining them in the replacements.txt in following form:

old-key-sequence := new-key-sequence

E.g.:
```°A := Å```
to write Å in the application you can now use the combination of °A.

When editing, check that non of old-key-sequences is part of a new-key-sequence, things like ```b := abc``` will turn out, that every time you type a key, the b in abc will be replaced with abc.

### Planned features ###

* ignoring of punctuation characters (DONE)
* in-program modification of character replacements
* improve pocket system (modifiable x)
* add helping support (number of characters, character suggestion,...)
* add reverse language learning (new language -> known language)
* fancy colors?
* suggest me more stuff.
