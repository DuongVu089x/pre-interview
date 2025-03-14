package gracefullshutdown

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefullShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Channel to wait for cleanup completion
	done := make(chan bool)

	go func() {
		<-stop
		fmt.Println("Graceful shutdown in progress...")
		time.Sleep(2 * time.Second) // Giả lập cleanup
		fmt.Println("Cleanup completed. Exiting with code 0.")

		done <- true // Signal that cleanup is complete

	}()

	fmt.Println("Running... Press Ctrl+C to stop.")

	time.Sleep(3 * time.Second)
	fmt.Println("Triggering graceful shutdown...")

	// Gửi tín hiệu SIGTERM để chương trình tự xử lý thoát
	syscall.Kill(os.Getpid(), syscall.SIGTERM)

	<-done

	for {
		println("running")
		time.Sleep(1 * time.Second)
	}

	// os.Exit(0) // Thoát với exit code 0 sau khi cleanup xong

}
