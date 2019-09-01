package main

import "fmt"

func ApplyConfig(config map[string]string) error {
	fmt.Println("[plugin] record.ApplyConfig is called")
	return nil
}
