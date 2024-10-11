package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"time"
)

const partecipantsFile = "partecipants.csv"

func writeToFile(csvObj string) {}

func toJSON() {}

func writeCSVHeader(record []string, w *csv.Writer, field1, field2 string) {
	record[0] = field1
	record[1] = field2
	fmt.Printf("%v", record)
	w.Write(record[:])
	w.Flush()
}

func main() {
	flag.NewFlagSet("gen", flag.ExitOnError)
	flag.NewFlagSet("sim", flag.ExitOnError)
	if len(os.Args) < 2 {
		fmt.Println("subcommand expected")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "gen":
		fmt.Println("generating csv")
		var err error
		genFile, err := os.Create(partecipantsFile)
		if err != nil {
			panic("error creating file")
		}

		defer genFile.Close()
		var name string
		var read int = 1
		id := 0
		var record [2]string
		// var record = make([]string, 2)
		w := csv.NewWriter(genFile)
		writeCSVHeader(record[:], w, "ID", "NAME")
		for read > 0 {
			fmt.Println()
			fmt.Println("partecipant name")
			read, err = fmt.Scanf("%s\n", &name)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic("error scanning name")
			}
			fmt.Printf("name is %s\n", name)
			record[0] = strconv.Itoa(id)
			record[1] = name
			fmt.Printf("%v", record)
			w.Write(record[:])
			w.Flush()
			id++
		}

	case "sim":
		fmt.Println("Reading file " + partecipantsFile)
		rFile, err := os.Open(partecipantsFile)
		if err != nil {
			panic("could not read file")
		}
		defer rFile.Close()
		r := csv.NewReader(rFile)
		var drafted []int
		var seen bool = false

		_, err = r.Read() //read the header
		if err != nil {
			if err == io.EOF {
				return
			}
			panic("could not read record")
		}
		for {
			current, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic("could not read record")
			}
			currId, err := strconv.Atoi(current[0])
			if err != nil {
				log.Println("error parsing int")
				break
			}
			currName := current[1]
			for i := range drafted {
				if currId == drafted[i] {
					seen = true
				}
				if seen {
					break
				}
			}
			drafted = append(drafted, currId)
			fmt.Println("Partecipant: " + currName)
			fmt.Printf("%v\n", drafted)
		}
		fmt.Printf("fight!\n")
		var undefeated []int = slices.Clone(drafted)

		// time to fight
		for {
			time.Sleep(1 * time.Second)
			fmt.Printf("\nundefeated: %v\n", undefeated)
			if len(undefeated) > 2 {
				_, loser := simFight(undefeated)

				// delete the element from the slice
				undefeated = append(undefeated[:loser], undefeated[loser+1:]...) //unpack the second slice to avoid creating an intermediate slice for the second part of the original slice
				// fmt.Printf("undefeated: %v\n", undefeated)
				fmt.Printf("len: %d\n", len(undefeated))
			} else if len(undefeated) == 2 {
				fmt.Printf("\nfinal fight:\n")
				winner, _ := simFight(undefeated)
				fmt.Printf("Tournament Winner: %d\n", undefeated[winner])
				break
			} else {
				break
			}
		}

	default:
		fmt.Println("Expected subcommands gen or sim")
	}
}

func simFight(undefeated []int) (winner int, loser int) {
	//indexes
	p1 := rand.Intn(len(undefeated))
	p2 := rand.Intn(len(undefeated))
	for {
		if p1 == p2 {
			p2 = rand.Intn(len(undefeated))
		} else {
			break
		}
	}

	winner = p1
	loser = p2
	if rand.Int() == 1 {
		winner = p2
		loser = p1
	}

	fmt.Printf("%d VS %d\n", undefeated[p1], undefeated[p2])
	fmt.Printf("Winner: %d\n", undefeated[winner])
	fmt.Printf("Loser: %d\n", undefeated[loser])
	return winner, loser
}
