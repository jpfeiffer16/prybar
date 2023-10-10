package main

// USING_CGO

import (
	"log"
	"syscall"

	"github.com/chzyer/readline"
	"github.com/creack/pty"
	"golang.org/x/term"

	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	// "regexp"
	// "strconv"
	// "strings"
)

type Hoon struct {
	ps1 string
	ps2 string
}

var u   *exec.Cmd
var tty *os.File

func (p Hoon) Open() {
	// println("Open")
	// Create arbitrary command.
        c := exec.Command(GetPath(), "zod")
	// u = exec.Command(GetPath(), "zod", "-t")

        // Start the command with a pty.
        ptmx, err := pty.Start(c)
        // if err != nil {
        //         return err
        // }
        // Make sure to close the pty at the end.
        defer func() { _ = ptmx.Close() }() // Best effort.

        // Handle pty size.
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGWINCH)
        go func() {
                for range ch {
                        if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
                                log.Printf("error resizing pty: %s", err)
                        }
                }
        }()
        ch <- syscall.SIGWINCH // Initial resize.
        defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

        // Set stdin in raw mode.
        oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
        if err != nil {
                panic(err)
        }
        defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

        // Copy stdin to the pty and the pty to stdout.
        // NOTE: The goroutine will keep reading until the next keystroke before returning.
        go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
        _, _ = io.Copy(os.Stdout, ptmx)

	// u = exec.Command(GetPath(), "zod", "-t")
	// f, err := pty.Start(u)
	// if err != nil {
	// 	panic(err)
	// }
	// tty = f;

	// io.Copy(os.Stdout, f)


	// c := exec.Command("grep", "--color=auto", "bar")
	// f, err := pty.Start(c)
	// if err != nil {
	// 	panic(err)
	// }
    //
	// go func() {
	// 	f.Write([]byte("foo\n"))
	// 	f.Write([]byte("bar\n"))
	// 	f.Write([]byte("baz\n"))
	// 	f.Write([]byte{4}) // EOT
	// }()
	// io.Copy(os.Stdout, f)

	// stdin, err := urb.StdinPipe()
	// io.WriteString(stdin, command+"\n")
	// stdin.Close()
	// out, err := u.Output()
	// u.StderrPipe.fmt.Printf("u.StderrPipe: %v\n", u.StderrPipe)
	// if err != nil {
	// 	panic(err)
	// }
	// println(out)

	// u.Stderr = os.Stderr
	// u.Stdout = os.Stdout
	// u.Run()
}

func (p Hoon) SetPrompts(ps1, ps2 string) {
	p.ps1 = ps1
	p.ps2 = ps2
}

func GetPath() string {
	pat, err := exec.LookPath("urbit")
	if err != nil {
		panic(err)
	}
	return pat
}

func (p Hoon) Version() string {
	urb := exec.Command(GetPath(), "--version")
	versions, err := urb.Output()
	if err != nil {
		panic(err)
	}

	return string(versions)
}

func (p Hoon) EvalFile(file string, args []string) int {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	out, _, err := RunCommand(string(fileContents))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

	return 0
}

func RunCommand(command string) (string, bool, error) {
	io.Copy(os.Stdout, tty)
	io.Copy(tty, os.Stdin)
	io.WriteString(tty, command+"\n")
	// out, err := urb.Output()
	// if err != nil {
	// 	return "", false, err
	// }

	// stringOut := string(out)

	// lines := strings.Split(stringOut, "\n")

	// needsMoreInput := true
	// if len(lines) > 1 && len(lines[0]) == 21 {
	// 	needsMoreInput = needsMoreInput && (lines[0][11:16] == "/eval")
    //
	// 	rxp := regexp.MustCompile(`\{(\d+) (\d+)\}`)
	// 	lineMatches := rxp.FindSubmatch([]byte(lines[1]))
	// 	if len(lineMatches) == 3 {
	// 		l, _ := strconv.Atoi(string(lineMatches[1]))
	// 		c, _ := strconv.Atoi(string(lineMatches[2]))
    //
	// 		realL := len(strings.Split(command, "\n"))
    //
	// 		needsMoreInput = needsMoreInput && (realL+1 == l) && (c == 1)
    //
	// 	} else {
	// 		needsMoreInput = false
	// 	}
    //
	// } else {
	// 	needsMoreInput = false
	// }
    //
	// return stringOut, needsMoreInput, nil

	return "test", false, nil
}

func (p Hoon) EvalExpression(code string) string {
	out, _, err := RunCommand(code)

	if err != nil {
		panic(err)
	}

	return string(out)
}

func (p Hoon) REPL() {
	for {
		line, err := readline.Line("--> ")
		if err != nil {
			break
		}
		readline.AddHistory(line)
		out, needMoreInput, err := RunCommand(line)

		for needMoreInput {
			newLine, err := readline.Line("... ")
			if err != nil {
				break
			}
			readline.AddHistory(newLine)
			line = line + "\n" + newLine
			out, needMoreInput, err = RunCommand(line)
		}

		if err != nil {
			panic(err)
		}

		strOut := string(out)
		fmt.Println(strOut)
	}
}

func (p Hoon) Close() {
}

var Instance = Hoon{}
