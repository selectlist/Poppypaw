package state

import "fmt"

func Close() {
	fmt.Println("Closing connections...")

	err := Discord.Close()

	if err != nil {
		fmt.Println("Error closing Discord session:", err)
	}

	Database.Close()
	err = Redis.Close()

	if err != nil {
		fmt.Println("Error closing Redis client:", err)
	}
}