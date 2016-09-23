package bash

import (
	"fmt"
	"fuzzer"
	"fuzzer/src/common"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

type spFunc func(string)

type BashInterface struct {
	commonCmd   []string
	specialCmd  []string
	specialFunc []spFunc
	readline    *readline.Instance
	running     bool
}

func CreateBashApp() *BashInterface {
	bsh := BashInterface{
		commonCmd:  []string{"ls", "cat", "rm"},
		specialCmd: []string{"exit", "cd", "-upload", "-download", "-sys", "-encode", "-decode", "-info"},
	}

	bsh.specialFunc = []spFunc{
		bsh.Exit,
		bsh.SendCd,
		bsh.SendUpload,
		bsh.SendDownload,
		bsh.Sys,
		bsh.Encode,
		bsh.Decode,
		bsh.Info,
	}

	bsh._buildPrompt()

	return &bsh
}

func (b *BashInterface) _buildPrompt() {
	autocompleter := readline.NewPrefixCompleter()
	allCmd := append(b.commonCmd, b.specialCmd...)

	for _, item := range allCmd {
		child := readline.PcItem(item)
		autocompleter.SetChildren(append(autocompleter.GetChildren(), child))
	}

	config := &readline.Config{
		Prompt:          "\033[32m•\033[0m\033[32m»\033[0m [Bash]$ ",
		HistoryFile:     "/tmp/readlinebash.tmp",
		AutoComplete:    autocompleter,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	l, err := readline.NewEx(config)

	if err != nil {
		panic(err)
	}

	b.readline = l
}

func (b *BashInterface) loop() {
	defer b.readline.Close()
	log.SetOutput(b.readline.Stderr())
	CallClear()

	for b.IsRunning() {
		line, err := b.readline.Readline()
		if err == readline.ErrInterrupt || err == io.EOF {
			return
		}

		if len(line) == 0 {
			continue
		}

		CallClear()
		fmt.Println(b.readline.Config.Prompt + line)

		b.Run(line)
	}
}

func (b *BashInterface) Run(l string) {
	if strings.Contains(l, "&&") {
		b.SendRaw(l)
		return
	}

	arr := strings.Fields(l)
	for i, item := range b.specialCmd {
		if item == arr[0] {
			b.specialFunc[i](l)
			return
		}
	}

	b.SendRaw(l)
}

func (b *BashInterface) Start() {
	if !fuzzer.NET.IsSetup() {
		e := fuzzer.SetupErr()
		e.Error()
		return
	}

	b.running = true
	b.loop()
}

func (b *BashInterface) Stop() {
	b.running = false
	fuzzer.CMD.SetContext("")
}

func (b *BashInterface) IsRunning() bool {
	return b.running
}

func (b *BashInterface) SetPrompt(p string) {
	b.readline.SetPrompt(p)
}

func (b *BashInterface) Exit(str string) {
	b.Stop()
}

func (b *BashInterface) Sys(str string) {
	arr := strings.Fields(str)

	if len(arr) < 2 {
		return
	}

	arr = append(arr[1:], arr[len(arr):]...)
	full := strings.Join(arr, " ")

	common.Syscall(full)
}

func (b *BashInterface) Encode(str string) {
	strArr := strings.Fields(str)

	if len(strArr) < 2 {
		return
	}

	sEnc := fuzzer.Encode(strArr[1])
	fmt.Println(sEnc)
}

func (b *BashInterface) Decode(str string) {
	strArr := strings.Fields(str)

	if len(strArr) < 2 {
		return
	}

	sDec := fuzzer.Decode(strArr[1])
	fmt.Println(sDec)
}

func (b *BashInterface) Info(str string) {
	if fuzzer.NET.GetResponse() == nil {
		fmt.Println("You havn't made any request. You must make a request before being able to see information")
		return
	}

	fmt.Println("Request => ")
	b.RequestInfo(str)

	fmt.Println("Response => ")
	b.ResponseInfo(str)
}

func (b *BashInterface) RequestInfo(str string) {
	flag := false
	r := fuzzer.NET.GetRequest()

	if strings.Contains(str, "-url") {
		fmt.Println(r.URL)
		flag = true
	}

	if strings.Contains(str, "-method") {
		fmt.Println(r.Method)
		flag = true
	}

	if strings.Contains(str, "-body") {
		fmt.Println(r.PostForm)
		flag = true
	}

	if strings.Contains(str, "-header") {
		fmt.Println(r.Header)
		flag = true
	}

	if !flag {
		fmt.Println(r)
	}
}

func (b *BashInterface) ResponseInfo(str string) {
	flag := false
	r := fuzzer.NET.GetResponse()

	if strings.Contains(str, "status") {
		fmt.Println(r.Status)
		flag = true
	}

	if strings.Contains(str, "body") {
		body := fuzzer.NET.GetBodyStr(r)
		fmt.Println("body: " + body)

		flag = true
	}

	if strings.Contains(str, "headers") {
		fmt.Println(r.Header)
		flag = true
	}

	if !flag {
		fmt.Println(r)
	}
}
