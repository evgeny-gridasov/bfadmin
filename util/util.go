package util

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
)

func Atoi(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func Atof(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
}

func MakeId(mapName string, gameType string, modName string) string {
	return mapName + ":" + gameType + ":"+ modName
}

func MakeNameRefractor1(mapName string, modName string) string {
	var m = modName
	switch modName {
	case "bfvietnam":
		m = "bfv"
	case "dc_final":
		m = "dcf"
	case "desertcombat":
		m = "dc"
	case "bf1942":
		m = "1942"
	case "bf1918":
		m = "1918"
	}
	return "[" + m + "] " + mapName
}

func MakeNameRefractor2(mapName string, gameType string, size string) string {
	gt := strings.TrimPrefix(gameType, "gpm_")
	return mapName + " [" + gt + "-" +size+ "] "
}

func ReadPropertiesFile(file string) map[string]string{
	m := make(map[string]string)
	f, err := os.Open(file)
	defer f.Close()
	if CheckErr(err) {
		return m
	}
	reader := bufio.NewReader(f)
	for {
		readString, err := reader.ReadString('\n')
		readString = strings.TrimSpace(readString)
		fields := strings.SplitN(readString, "=", 2)
		if len(fields) == 2 && !strings.HasPrefix(fields[0], "#") {
			m[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
		}
		if err != nil {
			break
		}
	}
	return m
}

func ParseCommandLine(line string) []string {
	var ret []string
	esc := false
	var arg bytes.Buffer
	for i:=0; i< len(line); i++ {
		c := line[i]
		if c == ' ' && !esc {
			if arg.Len() > 0 {
				ret = append(ret, arg.String())
			}
			arg.Reset()
			continue
		}
		if c == '\\' && !esc {
			esc = true
			continue
		}
		esc = false
		arg.WriteByte(c)
	}
	if arg.Len() > 0 {
		ret = append(ret, arg.String())
	}
	return ret
}

func CheckErr(err error) bool {
	if err != nil {
		log.Printf("%T: %s", err, err)
		return true
	}
	return false
}