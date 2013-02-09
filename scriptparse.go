package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"log"
	"regexp"
	"strings"
)

// Lines we want to turn into a valid script are of the form
//		<xxxx> [\{charName}]message

var SPEAKER_MAP = map[string] string {
	"Haruka":	"[img]http://lpix.org/1069382/haruka.png[/img]",
	"Kengo":	"[img]http://lpix.org/1069381/kengo.png[/img]",
	"Komari":	"[img]http://lpix.org/1069383/komari.png[/img]",
	"Kud":		"[img]http://lpix.org/1069384/kud.png[/img]",
	"Kurugaya":	"[img]http://lpix.org/1069385/kurugaya.png[/img]",
	"Kyousuke":	"[img]http://lpix.org/1069386/kyousuke.png[/img]",
	"Masato":	"[img]http://lpix.org/1069377/masato.png[/img]",
	"Mio":		"[img]http://lpix.org/1069387/mio.png[/img]",
	"Riki":		"[img]http://lpix.org/1069378/riki.png[/img]",
	"Rin":		"[img]http://lpix.org/1069388/rin.png[/img]",
	"Sasami":	"[img]http://lpix.org/1069389/sasami.png[/img]",
	"Voice":	"[img]http://lpix.org/1069379/unknown.png[/img]",
	"default":	"[img]http://lpix.org/1069379/unknown.png[/img]",
}

// Command-line flags. Right now only a single argument but there may
// be some expansion in the future
var (
	scriptNumber string
)


// Trims off the beginning line number <xxxx>
func trimNumber(line string) string {
	return (strings.SplitN(line, " ", 2))[1]
}

// Splits the line into a speaker, if one exists, and a message
// If there is no speaker, the first return value will be the empty string
func splitMessage(line string) (string, string) {
	var speaker, message string
	speakerEndIndex := 0

	// A speaker is indicated by \{charName} at the start of the line
	if line[0] == '\\' {
		// TODO: I know I'm being stupid, there is 100% an obvious and better
		// way to do this that right now I can't see

		// Iterate from the { until finding the matching }. The speaker name is
		// everything in between
		for speakerEndIndex = 2; line[speakerEndIndex] != '}'; speakerEndIndex++ {
		}
		speaker = line[2:speakerEndIndex]

		// Skip past the } to start the message
		message = line[speakerEndIndex+1:]
	} else {
		// The entire line is the message. TODO: Pretty sure I can just totally omit
		// this else clause but I'll test that later.
		message = line[speakerEndIndex:]
	}
	return speaker, message
}

// After splitting out the speaker and message, consult the speaker map to
// encode the line into part of the final update
func bbEncodeLine(speaker, prevSpeaker, message string) string {
	if speaker == "" {
		return message
	}
	encodedSpeaker := SPEAKER_MAP[speaker]

	// Represent faceless speakers with a bold name. If following a faced-
	// speaker, add a line break before the text for spacing purposes.
	// If a faced-speaker follows a faceless speaker, also add a line break
	// TODO: It's probably more elegant to combine this logic with the speaking
	//   vs. narration logic in main() somehow.
	if encodedSpeaker == "" {
		if SPEAKER_MAP[prevSpeaker] != "" {
			encodedSpeaker = "\n"
		}
		encodedSpeaker += "[b]" + speaker + "[/b]: "
	} else if prevSpeaker != "" && SPEAKER_MAP[prevSpeaker] == "" {
		encodedSpeaker = "\n" + encodedSpeaker
	}

	return encodedSpeaker + message
}

// Returns true iff we want to consider this line in the script. We do this
// by checking if there is a line number at the start, because we want to
// consider a line iff it begins with a line number
func isScriptLine(line string) bool {
	return strings.Index(line, "<") == 0
}

// Removes display control sequences from the script message, like \wait{} and
// \shake{}. We do this by matching against a regex and replacing matches with
// an empty string. The argument to this function should be the message returned
// by splitMessage--sending the raw line beforehand will strip out the speaker
func removeExtraneousControls (line string) string {
	controlRegex := regexp.MustCompile("\\\\.*{{1}.*}{1}")
	return controlRegex.ReplaceAllString(line, "")
}

// Run before main is called due to how Go works. Parses all the command line flags
func init() {
	flag.StringVar(&scriptNumber, "script", "0513", "The script number to be parsed, e.g. 0513")
	flag.Parse()
}

func main() {
	// Open our raw input file and create the output file that will contain
	// the trimmed script
	rawFile, err := os.Open("SEEN" + scriptNumber + ".sjs")
	if err != nil {
		log.Fatal("Unable to open file: SEEN" + scriptNumber + ".sjs")
	}
	defer rawFile.Close()

	trimmedFile, err := os.Create("script_" + scriptNumber + ".txt")
	if err != nil {
		log.Fatal("Unable to create file: script_" + scriptNumber + ".txt")
	}
	defer trimmedFile.Close()
	
	fmt.Println("Parsing file " + "SEEN" + scriptNumber + ".sjs...")

	rawFileBuf := bufio.NewReader(rawFile)
	trimmedFileBuf := bufio.NewWriter(trimmedFile)

	// Read, process, and write one line at a time
	line := ""

	// Keep track of whether the previous line was speech, and if so,
	// who was speaking. Use this to space out text with conditional
	// line breaks
	prevWasSpeaker := false
	prevSpeaker := ""

	for err == nil {
		line, err = rawFileBuf.ReadString('\n')

		// Only process valid lines that are also script lines
		if err != nil || !isScriptLine(line) {
			continue
		}

		// Do all the processing steps on the line to produce an actual script
		// with BBCode pictures and spacing and all
		trimmedLine := trimNumber(line)
		speaker, message := splitMessage(trimmedLine)
		message = removeExtraneousControls(message)
		finalEncode := bbEncodeLine(speaker, prevSpeaker, message)

		// Is the current line narration or speech?
		currentIsSpeaker := speaker != ""

		// If we switched between speaking and narrating, add a line break before
		// the final message
		if currentIsSpeaker != prevWasSpeaker {
			finalEncode = "\n" + finalEncode
		}

		// Save the current speech fields as previous speech fields for the next
		// iteration
		prevWasSpeaker = currentIsSpeaker
		prevSpeaker = speaker

		if _, writeErr := trimmedFileBuf.WriteString(finalEncode); writeErr != nil {
			log.Println(writeErr.Error())
		}
		trimmedFileBuf.Flush()
	}
}
