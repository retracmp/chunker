package compile

import (
	"acid/chunker/src/compile"
	"acid/chunker/src/helpers"
)

func Start() error {
	t := helpers.NewPerformanceTimer()
	defer t.EndTimer()

	compiler := compile.NewCompiler(`./builds/++Fortnite+Release-1.8-CL-3724489-Windows.acidmanifest`, `./builds/downloaded`, `./build`)
	if err := compiler.Compile(); err != nil {
		return err
	}
	
	return nil
}