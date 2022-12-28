package main

import "atlassian_backup/processor"

func main() {
	processor := processor.New()

	processor.Process()

}
