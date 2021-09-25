# Getting Started

## Setting up the project
First we will create a folder for our project. Using the command
```ps1
mkdir obboy-blockchain
cd obboy-blockchain
```

We will start the project with go modules initialization. Lets open the folder in vscode with 
```ps1
code .
```
Open up the built in terminal with **ctrl+shift+`**. Now to initialize go modules
```ps1
go mod init github.com/<username>/<project name>
// go mod init github.com/kalam360/obboy-blockchain
```

We will develop the cli first. create a folder name `cmd/obboy` and create the `main.go` file in it. Our `main.go` file will now have the basic command for our blockchain. We will use the `cobra` package for building the cli. First you need to setup the `gopath` in the system. Now lets get the `cobra` package with 
```ps1
go get github.com/spf13/cobra
```
Inside the `main.go` file. Write the following lines. I will discuss the lines later. 
```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
    // create the cobra command 
	var obboyCmd = &cobra.Command {
		Use: "obboy",
		Short: "Obboy Blockchain Cli",
        // it will run the following function when called
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

    // This will execute the command in the terminal and catch any error and print it. 
	err := obboyCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```
Lets now build the cli and install in the `GOPATH/bin`.
```ps1
go install ./cmd/...
```


## Blockchain is a database
### Genesis File
### State
### Transactions

