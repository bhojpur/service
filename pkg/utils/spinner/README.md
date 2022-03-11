# Bhojpur Service - Spinner Library

The `spinner` is a simple library to add a spinner / progress indicator to any terminal application.

## Installation

```bash
go get github.com/bhojpur/service/pkg/utils/spinner
```

## Features

* Start
* Stop
* Restart
* Reverse direction
* Update the spinner character set
* Update the spinner speed
* Prefix or append text
* Change spinner color, background, and text attributes such as bold / italics
* Get spinner status
* Chain, pipe, redirect output
* Output final string on spinner/indicator completion

## Examples

```Go
package main

import (
	"github.com/bhojpur/service/pkg/utils/spinner"
	"time"
)

func main() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)  // Build our new spinner
	s.Start()                                                    // Start the spinner
	time.Sleep(4 * time.Second)                                  // Run for some time to simulate work
	s.Stop()
}
```

## Update the character set and restart the spinner

```Go
s.UpdateCharSet(spinner.CharSets[1])  // Update spinner to use a different character set
s.Restart()                           // Restart the spinner
time.Sleep(4 * time.Second)
s.Stop()
```

## Update spin speed and restart the spinner

```Go
s.UpdateSpeed(200 * time.Millisecond) // Update the speed the spinner spins at
s.Restart()
time.Sleep(4 * time.Second)
s.Stop()
```

## Reverse the direction of the spinner

```Go
s.Reverse() // Reverse the direction the spinner is spinning
s.Restart()
time.Sleep(4 * time.Second)
s.Stop()
```

## Provide your own spinner

```Go
someSet := []string{"+", "-"}
s := spinner.New(someSet, 100*time.Millisecond)
```

## Prefix or append text to the spinner

```Go
s.Prefix = "prefixed text: " // Prefix text before the spinner
s.Suffix = "  :appended text" // Append text after the spinner
```

## Set or change the color of the spinner. Default color is white. The spinner will need to be restarted to pick up the change

```Go
s.Color("red") // Set the spinner color to red
```

You can specify both the background and foreground color, as well as additional attributes such as `bold` or `underline`.

```Go
s.Color("red", "bold") // Set the spinner color to a bold red
```

To set the background to black, the foreground to a bold red:

```Go
s.Color("bgBlack", "bold", "fgRed")
```

Below is the full color and attribute list:

```Go
// default colors
red
black
green
yellow
blue
magenta
cyan
white

// attributes
reset
bold
faint
italic
underline
blinkslow
blinkrapid
reversevideo
concealed
crossedout

// foreground text
fgBlack
fgRed
fgGreen
fgYellow
fgBlue
fgMagenta
fgCyan
fgWhite

// foreground Hi-Intensity text
fgHiBlack
fgHiRed
fgHiGreen
fgHiYellow
fgHiBlue
fgHiMagenta
fgHiCyan
fgHiWhite

// background text
bgBlack
bgRed
bgGreen
bgYellow
bgBlue
bgMagenta
bgCyan
bgWhite

// background Hi-Intensity text
bgHiBlack
bgHiRed
bgHiGreen
bgHiYellow
bgHiBlue
bgHiMagenta
bgHiCyan
bgHiWhite
```

## Generate a sequence of numbers

```Go
setOfDigits := spinner.GenerateNumberSequence(25)    // Generate a 25 digit string of numbers
s := spinner.New(setOfDigits, 100*time.Millisecond)
```

## Get spinner status

```Go
fmt.Println(s.Active())
```

## Unix pipe and redirect)

Setting the Spinner Writer to Stderr helps show progress to the user, with the enhancement to chain, pipe, or redirect the output.

This is the preferred method of setting a Writer at this time.

```go
s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
s.Suffix = " Encrypting data..."
s.Start()
// Encrypt the data into ciphertext
fmt.Println(os.Stdout, ciphertext)
```

```sh
> myprog encrypt "Secret text" > encrypted.txt
⣯ Encrypting data...
```

```sh
> cat encrypted.txt
1243hjkbas23i9ah27sj39jghv237n2oa93hg83
```

## Final String Output

Add additional output when the spinner/indicator has completed. The "final" output string can be multi-lined and will be written to wherever the `io.Writer` has been configured for.

```Go
s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
s.FinalMSG = "Complete!\nNew line!\nAnother one!\n"
s.Start()                 
time.Sleep(4 * time.Second)
s.Stop()                   
```

Output

```sh
Complete!
New line!
Another one!
```
