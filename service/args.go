package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func (p *params) getArgs() error {
	helpMsg := "Params:\n" +
		" -s, -src, --src		source folders\n" +
		" -d, -dst, --dst		destination folders\n" +
		" -m, -mask, --mask		file mask, eg *2021-0[5-7]-*.log (default *)\n" +
		" -p, -parallel, --parallel	count of parallel mv-operations\n" +
		" -h, -help, --help		this help message\n"

	args := os.Args[1:]

	if len(args) == 0 {
		log.Println(helpMsg)
		return fmt.Errorf("No parameters. Boredom :'/\n")
	}

	p.Parallel = 1
	p.FileMask = "*"

	mode := 0
	for _, arg := range args {
		switch arg {
		case "--src", "-src", "-s":
			mode = 1
		case "--dst", "-dst", "-d":
			mode = 2
		case "--mask", "-mask", "-m":
			mode = 4
		case "--help", "-help", "-h":
			log.Println(helpMsg)
			return nil
		case "--parallel", "-parallel", "-p":
			mode = 3
		case "-norand":
			p.NoRand = true
		default:
			switch mode {
			case 1:
				p.Src = append(p.Src, arg)
			case 2:
				p.Dst = append(p.Dst, arg)
			case 3:
				p.Parallel, _ = strconv.Atoi(arg)
				mode = 0
			case 4:
				p.FileMask = arg
				mode = 0
			}
		}
	}
	return nil
}
func (p *params) checkArgs() error {
	src := make([]string, 0)
	for _, v := range p.Src {
		if p.checkDouble(src, v) {
			continue
		}
		info, err := os.Stat(v)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("%v is not a dir", v)
		}
		path, err := filepath.Abs(v)
		if err != nil {
			return err
		}
		src = append(src, path)
	}
	p.Src = src
	log.Println("src: ", p.Src)

	dst := make([]string, 0)
	for _, v := range p.Dst {
		if p.checkDouble(dst, v) {
			continue
		}
		info, err := os.Stat(v)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("%v is not a dir", v)
		}
		path, err := filepath.Abs(v)
		if err != nil {
			return err
		}
		dst = append(dst, path)
	}
	p.Dst = dst
	log.Println("dst: ", p.Dst)
	return nil
}

func (p *params) checkDouble(arr []string, s string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
