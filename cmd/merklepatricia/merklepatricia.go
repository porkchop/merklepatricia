package main

import (
	"github.com/codegangsta/cli"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/server"
	"os"

	application "github.com/porkchop/merklepatricia/app"
)

func main() {
	app := cli.NewApp()
	app.Name = "cli"
	app.Usage = "cli [command] [args...]"
	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the MerklePatricia server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address",
					Value: "unix://data.sock",
					Usage: "MerklePatricia server listen address",
				},
			},
			Action: func(c *cli.Context) {
				cmdServer(app, c)
			},
		},
	}
	app.Run(os.Args)

}

//--------------------------------------------------------------------------------

func cmdServer(app *cli.App, c *cli.Context) {
	addr := c.String("address")
	mApp := application.NewMerklePatriciaApp()

	// Start the listener
	s, err := server.NewServer(addr, mApp)
	if err != nil {
		Exit(err.Error())
	}

	// Wait forever
	TrapSignal(func() {
		// Cleanup
		s.Stop()
	})
}
