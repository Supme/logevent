// Example use:
//   logevent.exe -m "wmi_logevent_bad_login_count" -d "Bad login event"
// Create sheduler task for log with required code
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	filename    string
	metric      string
	description string
)

func init() {
	flag.StringVar(&filename, "f", "C:\\Program Files\\wmi_exporter\\textfile_inputs\\logevent.prom", "Metric file")
	flag.StringVar(&metric, "m", "wmi_logevent_count", "Metric name")
	flag.StringVar(&description, "d", "Log event count", "Metric description")
	flag.Parse()
}

func main() {
	tmpfilename := filename + "$$"

	rand.Seed(time.Now().Unix())
	timeout := time.Tick(10 * time.Second)

	for {
		select {
		case <-timeout:
			errorExit("Wait timeout when remove temp file")
		default:
			_, err := os.Stat(tmpfilename)
			if os.IsNotExist(err) {
				goto EndLoop
			}
			r := rand.Intn(20) + 10
			time.Sleep(time.Duration(r) * time.Millisecond)
		}
	}

EndLoop:

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 666)
	if err != nil {
		errorExit(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	found := false
	for i := range lines {
		if strings.HasPrefix(lines[i], metric) {
			found = true
			s := strings.Split(lines[i], " ")
			var d int
			if len(s) == 2 {
				d, err = strconv.Atoi(s[1])
				if err != nil {
					errorExit("Can not parse to int metric '%s' in line %d", s[1], i+1)
				}
				lines[i] = fmt.Sprintf("%s %d", metric, d+1)
			} else {
				errorExit("Bad metric format in line %d", i+1)
			}
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("# HELP %s %s", metric, description))
		lines = append(lines, fmt.Sprintf("# TYPE %s counter", metric))
		lines = append(lines, fmt.Sprintf("%s 1", metric))
	}

	if err := file.Close(); err != nil {
		errorExit("Close file with error: %s", err)
	}

	tmpfile, err := os.OpenFile(tmpfilename, os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		errorExit(err.Error())
	}
	defer file.Close()

	if err := tmpfile.Sync(); err != nil {
		errorExit("Create temp file error: %s", err.Error())
	}

	for i := range lines {
		if _, err := tmpfile.WriteString(lines[i] + "\r\n"); err != nil {
			errorExit("Write line %d to temp file with error: %s", i+1, err.Error())
		}
	}

	if err := tmpfile.Sync(); err != nil {
		errorExit("Sync temp file with error: %s", err.Error())
	}

	if err := tmpfile.Close(); err != nil {
		errorExit("Close temp file with error: %s", err.Error())
	}

	if err := os.Rename(tmpfilename, filename); err != nil {
		errorExit("Rename temp file to file with error: %s", err.Error())
	}

}

func errorExit(format string, v ...interface{}) {
	log.Printf(format, v...)
	os.Exit(2)
}
