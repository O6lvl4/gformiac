// gformiac is an Infrastructure-as-Code tool for managing Google Forms from
// declarative YAML definitions. It supports plan (dry-run), apply, and import
// workflows backed by the Google Forms API.
package main

import "github.com/O6lvl4/gformiac/cmd"

func main() {
	cmd.Execute()
}
