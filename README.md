LB-Parser
=========
LB-Parser is a utility for parsing the script files used by the Key visual novel <i>Little
Busters!</i> and encoding them in BBcode for the purpose of publishing as a Let's Play. The
utility is written in Go. It's extremely unlikely anyone would find this useful as anything
other than insight into an undergraduate's coding style, but I don't have any reason not to
release it.

How it works
------------
LB-Parse takes a `.sjs` file as input. <i>Little Busters!</i> utilizes several of these files.
The game's engine itself uses a `.sjs` file along with a `.ke` file to create the game
experience--the former contains dialogue and narration, the latter contains instructions on
structure and the appearance, change, and animation of character sprites and backgrounds.

Usage
-----
Run `go run scriptparse.go` with the proper `.sjs` file in same directory. The output will
be in a file named `script_####.txt` where #### is the same as the 4-digit number in the
input file name. For now that name is hardcoded to be `SEEN0513.sjs`. In the future the
input file will be a command-line argument.

File format - .sjs
------------------
The first line of the file contains a comment relating to what packed `.txt` file it relates
to, for example

	// Resources for SEEN0513.TXT

is the first line of `SEEN0513.sjs`, the script file for the gameplay on May 13th. Following
that is a blank line, then a set of character declarations of the form

	#character 'Riki'

which represent all of the speakers present in the script within that file. Next is another
blank line, followed by a set of script lines. An example such line is

	<0009> \{Masato}"...To the fight."

The four-digit number at the start of the line is a line number of sorts, and increments
every line of the file. If a character speaks the current line, it's represented by
`\{characterName}`. The remainder of the line is the message displayed in the dialogue box.

Other backslash control sequences exist to change the presentation of the text. The command
`\shake{n}` takes an integer <i>n</i> as its argument and shakes the screen based on the
argument. The command `\wait{n}` waits for <i>n</i> milliseconds before drawing any more
of the text. LB-Parser removes all such control sequences before the final output.

Output
------
The output of the program is a BBcode-encoded version of the script suitable for pasting into
a BBcode-compatible forum as part of a Let's Play. Note that the `.sjs` file used as input
does not make any divisions for different choices and so you will have to manually remove
sections that do not correspond with your actual choices.

All major characters, when speaking, are given a 112x110px transparent portrait in front
of their line. Minor characters are represented with their name in bold. Line breaks are
added appropriately to ensure reasonable text spacing.

Future improvements
-------------------
* Remove the hard-coded input file name and make it a command-line argument
* Add support for calendar images by parsing text such as `May 13th (Sun)`
* Figure out why an occasional letter "p" appears where it should not
