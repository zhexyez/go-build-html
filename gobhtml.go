package main

import (
	"fmt"
	//"time"
	//"io"
	"bufio"
	"os"
	"os/exec"
	"strconv"
)

// A program to quickly create an HTML template
//
// Please, feel free to contribute

const Version string = "dev/1"
const colorReset = "\033[0m"
const colorRed   = "\033[0;31m"
const colorGreen = "\033[0;32m"
const colorYellow= "\033[0;33m"

// A function to clear the screen
func WipeScreen() {
	wipe := exec.Command("clear")
	wipe.Stdout = os.Stdout
	wipe.Run()
}

var (
	aio 		bool	// if true, the code must be in single file; otherwise HTML to .html, JS to .js, CSS TO .css, but in one directory and 1 file per type
	filename	string
	usrinteraction	string
)

// Base type of Block's attributes
type Attribute struct {
	name 	string
	value 	string
}

// Global counter to distinguish Blocks
var blockID int

// Base type of Block
type Block struct {
	ID		int
	definition 	string
	attributes 	[]Attribute
	content 	string
}

// Global counter to distinguish Blocks
var linkID int

// Buffer for created Links
var linksbuf []Link

// Base type to make hirearchy
type Link struct {
	ID int
	parent *Block
	child  *Block
}

// Buffer for created Blocks
var blocksbuf []Block

// Adds item to the buffer
func AddItem (newblock Block) {
	blocksbuf = append(blocksbuf, newblock)
}

// Function to interact with user
func interaction(exit chan<- bool) {
	var noexit bool = true
	for noexit {
		var input string
		fmt.Println("|------------------------------------------‐-----------------------------------------------|")
		fmt.Println(colorYellow, "     A<add new item>      ", colorReset, " /",colorYellow,"  L<link items> ",colorReset," /",colorYellow," C<list items>",colorReset," /",colorYellow,"  V<list links>", colorReset)
		fmt.Println(colorYellow,"R<remove item with links>",colorReset," /",colorYellow,"  W<wipes screen>",colorReset,"/ ",colorYellow,"    D<done> ",colorReset,"  /",colorYellow,"  E<exit no save>",colorReset)
		fmt.Println("|------------------------------------------------------------------------------------------|")
		fmt.Scanln(&input)
		switch input {
			case "A", "a":
				var newblock Block
				var itemproperty string
				for cnt := 0; cnt < 3; cnt++ {
					switch cnt {
						case 0:
							fmt.Println(colorYellow,"Enter element type (div, p, etc.)",colorReset)
							fmt.Scanln(&itemproperty)
							newblock.definition = itemproperty
							blockID++
							newblock.ID = blockID
						case 1:
							ondone := true
							var condition string
							for ondone {
								fmt.Println(colorYellow,"Input E to exit / I to input attributes",colorReset)
								fmt.Scanln(&condition)
								switch condition {
									case "E", "e":
										ondone = false
										break
									case "I","i":
										var newattr Attribute
										fmt.Println(colorYellow,"Attribute name",colorReset)
										fmt.Scanln(&itemproperty)
										newattr.name = itemproperty
										fmt.Println(colorYellow,"Attribute value",colorReset)
										scanner := bufio.NewScanner(os.Stdin)
										if scanner.Scan() {
											newattr.value = scanner.Text()
										}
										newblock.attributes = append(newblock.attributes, newattr)
									default:
										fmt.Println(colorRed,"Bad symbol(s). Try again",colorReset)
								}
							}
						case 2:
							fmt.Println(colorYellow,"Enter content",colorReset)
							scanner := bufio.NewScanner(os.Stdin)
							if scanner.Scan() {
								newblock.content = scanner.Text()
							}
					}
				}
				AddItem(newblock)
				fmt.Println(colorGreen,"New item created with ID", newblock.ID, "and added to buffer",colorReset)
			case "C", "c":
				if len(blocksbuf) == 0 {
					fmt.Println(colorRed,"There is no items",colorReset)
					break
				}
				for currentblock := range blocksbuf {
					block := blocksbuf[currentblock]
					fmt.Println()
					fmt.Println(colorGreen,"ID: ", block.ID)
					fmt.Println("type: ", block.definition)
					fmt.Println("attributes: ", block.attributes)
					fmt.Println("content: ", block.content,colorReset)
				}
			case "V","v":
				if len(linksbuf) == 0 {
					fmt.Println(colorRed,"There is no links",colorReset)
					break
				}
				for currentlink := range linksbuf {
					link := linksbuf[currentlink]
					fmt.Println()
					fmt.Println(colorGreen,"ID: ", link.ID)
					fmt.Println("parent ID: ", link.parent.ID)
					fmt.Println("child ID: ", link.child.ID,colorReset)
				}
			case "L","l":
				if len(blocksbuf) <= 1 {
					fmt.Println(colorRed,"At least 2 items must be present to make a link",colorReset)
					break
				}
				var newlink Link
				var strparent string
				var strchild string
				fmt.Println(colorYellow,"Parent ID of type int",colorReset)
				fmt.Scanln(&strparent)
				intparent, errparent := strconv.Atoi(strparent)
				fmt.Println(colorYellow,"Child ID of type int",colorReset)
				fmt.Scanln(&strchild)
				intchild, errchild := strconv.Atoi(strchild)
				if errparent != nil || errchild != nil {
					fmt.Println(colorRed,"string to int conversion failed. Try again",colorReset)
					break
				}
				if intparent == intchild {
					fmt.Println(colorRed,"Parent and child IDs cannot be the same",colorReset)
					break
				}
				var repeated bool
				if len(linksbuf) >= 1 {
					for i := range linksbuf {
						if linksbuf[i].parent.ID == intparent && linksbuf[i].child.ID == intchild { repeated = true }
						if linksbuf[i].parent.ID == intchild && linksbuf[i].child.ID == intparent { repeated = true }
					}
				}
				if repeated {
					fmt.Println(colorRed,"Such link already exist",colorReset)
					break
				}
				var parentexist bool
				var childexist bool
				for i := range blocksbuf {
					current := blocksbuf[i]
					if current.ID == intparent {
						parentexist = true
						newlink.parent = &current
					} else if current.ID == intchild {
						childexist = true
						newlink.child = &current
					}
				}
				if parentexist && childexist {
					linkID++
					newlink.ID = linkID
					linksbuf = append(linksbuf, newlink)
					fmt.Println(colorGreen,"New link created with ID", newlink.ID, "and added to buffer",colorReset)
				}
				if parentexist == false {
					fmt.Println(colorRed,"Parent ID does not exist",colorReset)
				}
				if childexist == false {
					fmt.Println(colorRed,"Child ID does not exist",colorReset)
				}
			case "R", "r":
				if len(blocksbuf) == 0 {
					fmt.Println(colorRed,"There is nothing to remove",colorReset)
					break
				}
				found := false
				var strid string
				fmt.Println(colorYellow,"Enter ID to remove",colorReset)
				fmt.Scanln(&strid)
				intid, err := strconv.Atoi(strid)
				if err != nil {
					fmt.Println(colorRed,"string to int conversion error. Try again",colorReset)
					break
				}
				if intid == 0 {
					fmt.Println(colorRed,"ID 0 is not applicable",colorReset)
					break
				}
				for i := range blocksbuf {
					if blocksbuf[i].ID == intid {
						found = true
						if len(linksbuf) >= 1 {
							for k := range linksbuf {
								if linksbuf[k].parent.ID == intid || linksbuf[k].child.ID == intid {
									currentlinkID := linksbuf[k].ID
									linksbuf[k] = linksbuf[len(linksbuf)-1]
									linksbuf = linksbuf[:len(linksbuf)-1]
									fmt.Println(colorGreen,"Removed link with ID", currentlinkID, "for item with ID", intid,colorReset)
								}
							}
						}
						blocksbuf[i] = blocksbuf[len(blocksbuf)-1]
						blocksbuf = blocksbuf[:len(blocksbuf)-1]
						fmt.Println(colorGreen,"Removed item with ID", intid,colorReset)
						break
					}
				}
				if found == false {
					fmt.Println(colorRed,"No item with ID", intid, "found",colorReset)
				}
			case "W", "w":
				WipeScreen()
			case "E", "e":
				noexit = false
				exit <- true
				break
			default:
				fmt.Println(colorRed,"Bad symbol(s). Try again",colorReset)
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(colorRed, "No filename specified. Please, provide [filename] as the first command-line argument", colorReset)
		return
	}
	for i := range os.Args {
		switch i {
			case 1:
				filename = os.Args[1] + ".html"
			case 2:
				if os.Args[2] != "aio" {
					fmt.Println(colorRed, "Unknown keyword ", os.Args[2], ". Must be aio for all-in-one file or leave it blank for separate files", colorReset)
					return
				}
				aio = true
		}
	}
	WipeScreen()
	fmt.Println(colorYellow, "Welcome to the HTML template constructor version", colorGreen, Version, colorReset)
	fmt.Println(colorYellow, "Use the following menu to make an objects. When you are done - press D. It will generate template and place it in the specified directory", colorReset)
	fmt.Println()
	ch_interaction := make(chan bool)
	go interaction(ch_interaction)
	interaction_complete := <- ch_interaction
	if interaction_complete == true {
		fmt.Println(colorYellow, "Program has peacefully exited", colorReset)
		return
	}
}
