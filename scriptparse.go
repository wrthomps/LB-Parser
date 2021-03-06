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

var SPEAKER_MAP map[string] string

var CALENDAR_MAP = map[string] string {
	"May 13th (Sun)": "[img]http://lpix.org/1071014/BS_DTA0513.png[/img]",
	"May 14th (Mon)": "[img]http://lpix.org/1071015/BS_DTA0514.png[/img]",
	"May 15th (Tue)": "[img]http://lpix.org/1071016/BS_DTA0515.png[/img]",
	"May 16th (Wed)": "[img]http://lpix.org/1071017/BS_DTA0516.png[/img]",
	"May 17th (Thu)": "[img]http://lpix.org/1071018/BS_DTA0517.png[/img]",
	"May 18th (Fri)": "[img]http://lpix.org/1071019/BS_DTA0518.png[/img]",
	"May 19th (Sat)": "[img]http://lpix.org/1071020/BS_DTA0519.png[/img]",
	"May 20th (Sun)": "[img]http://lpix.org/1071021/BS_DTA0520.png[/img]",
	"May 21st (Mon)": "[img]http://lpix.org/1071022/BS_DTA0521.png[/img]",
	"May 22nd (Tue)": "[img]http://lpix.org/1071023/BS_DTA0522.png[/img]",
	"May 23rd (Wed)": "[img]http://lpix.org/1071024/BS_DTA0523.png[/img]",
	"May 24th (Thu)": "[img]http://lpix.org/1071025/BS_DTA0524.png[/img]",
	"May 25th (Fri)": "[img]http://lpix.org/1071026/BS_DTA0525.png[/img]",
	"May 26th (Sat)": "[img]http://lpix.org/1071027/BS_DTA0526.png[/img]",
	"May 27th (Sun)": "[img]http://lpix.org/1071028/BS_DTA0527.png[/img]",
	"May 28th (Mon)": "[img]http://lpix.org/1071029/BS_DTA0528.png[/img]",
}

// Command-line flags. Right now only a single argument but there may
// be some expansion in the future
var (
	scriptNumber string
)

// Constructs the speaker map
func constructSpeakers(speakersByte []byte) map[string]string {
	speakerString := string(speakersByte)
	speakers := strings.Split(speakerString, "\n")
	m := make(map[string]string)
	for i := range speakers {
		m[strings.Split(speakers[i], ",")[0]] = strings.TrimSpace(strings.Split(speakers[i], ",")[1])
	}

	return m
}

// Trims off the beginning line number <xxxx>
func trimNumber(line string) string {
	return (strings.SplitN(line, " ", 2))[1]
}

// Splits the line into a speaker, if one exists, and a message
// If there is no speaker, the first return value will be the empty string
func splitMessage(line string) (string, string) {
	var speaker, message string
	message = line

	// A speaker is indicated by \{charName} at the start of the line
	if line[0] == '\\' {
		// If there is a speaker, the first } on the line will be the matching
		// brace. No other control sequences can be in the speaker's name
		speakerEndIndex := strings.Index(line, "}")
		speaker = line[2:speakerEndIndex]

		// Skip past the } to start the message
		message = line[speakerEndIndex+1:]
	} 

	// If there is no speaker, return ("", line)
	return speaker, message
}

// After splitting out the speaker and message, consult the speaker map to
// encode the line into part of the final update
func bbEncodeLine(speaker, prevSpeaker, message string) string {

	// If there's no speaker, return just the message with maybe a blank line
	if speaker == "" {
		if narrationSpeechSwap(speaker, prevSpeaker) {
			return "\n" + message
		} else {
			return message
		}
	}

	encodedSpeaker := SPEAKER_MAP[speaker]

	// Represent faceless speakers with a bold name. If following a faced-
	// speaker, add a blank line before the text for spacing purposes.
	// If a faced-speaker follows a faceless speaker, also add a blank line
	if encodedSpeaker == "" {
		if SPEAKER_MAP[prevSpeaker] != "" {
			encodedSpeaker = "\n"
		}
		encodedSpeaker += "[b]" + speaker + "[/b]: "
	} else if prevSpeaker != "" && SPEAKER_MAP[prevSpeaker] == "" {
		encodedSpeaker = "\n" + encodedSpeaker
	}

	finalEncode := encodedSpeaker + message

	// Add a blank line if switching between narration and speech
	if narrationSpeechSwap(speaker, prevSpeaker) {
		finalEncode = "\n" + finalEncode
	}
	return finalEncode
}

// Replace dates with their respective calendar images
func encodeDate(message string) string {
	return CALENDAR_MAP[strings.TrimSpace(message)] + "\n"
}

// Returns true iff the previous and current lines differ in regard to
// whether they are speech or narration
func narrationSpeechSwap(speaker, prevSpeaker string) bool {
	return (prevSpeaker == "" && speaker != "") || (prevSpeaker != "" && speaker == "") 
}

// Returns true iff the input string is a date indicator
func isDate(message string) bool {
	return CALENDAR_MAP[strings.TrimSpace(message)] != ""
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
	cleanedLine := controlRegex.ReplaceAllString(line, "")

	// Special case: Remove \p sequences
	cleanedLine = strings.Replace(cleanedLine, "\\p", "", -1)

	return cleanedLine
}

// Parses the command-line flags and reconstructs the speaker map from the
// spoiler-protected encoded serialization
func init() {
	flag.StringVar(&scriptNumber, "script", "0513", "The script number to be parsed, e.g. 0513")
	flag.Parse()

	encodedSpeakers, err := os.Open("speakers_encoded.txt")
	if err != nil {
		log.Fatal("Unable to open file: speakers_encoded.txt")
	}

	decodedSpeakers := decodeSpeakers(encodedSpeakers)
	SPEAKER_MAP = constructSpeakers(decodedSpeakers[:len(decodedSpeakers)-1])
}

// Writes a single line of output to the encoded script file
func writeOutputLine(trimmedFileBuf *bufio.Writer, finalEncode string) {
	if _, writeErr := trimmedFileBuf.WriteString(finalEncode); writeErr != nil {
		log.Println(writeErr.Error())
	}
	trimmedFileBuf.Flush()
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

		// Save the current speech fields as previous speech fields for the next
		// iteration
		prevSpeaker = speaker

		// For dates, replace them with the calendar images
		if isDate(message) {
			finalEncode = encodeDate(finalEncode)
		}

		writeOutputLine(trimmedFileBuf, finalEncode)
	}
}
