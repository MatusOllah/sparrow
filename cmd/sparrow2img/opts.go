package main

var opts struct {
	PosArgs struct {
		PNGPath string `description:"Path to png image"`
		XMLPath string `description:"Path to xml file"`
	} `positional-args:"yes" required:"yes"`
	Verbose    bool   `short:"v" long:"verbose" description:"Print verbose information"`
	Output     string `short:"o" long:"output" description:"Output directory" default:"out/"`
	SubTexture string `short:"t" long:"sub-texture" description:"Extract a specific frame / subtexture"`
}
