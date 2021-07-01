package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	// Rate limits
	sm := make(chan struct{}, 3) //Semaphore, available only 3 goroutines
	sm2 := make(chan struct{}, 3) //Semaphore, available only 3 goroutines
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	sm <- struct{}{}
	sm2 <- struct{}{}
	go handleSignals(cancel, &sm)
	go func() {
		defer func() { <-sm2 }()
		defer wg.Done()
		if err := startServer(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}

func handleSignals(cancel context.CancelFunc,  sm *chan struct{}) {
	defer func() { <-*sm }()
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt) //Signal for Interrupt for Graceful Shutdown
	for {
		sig := <-sigCh
		switch sig { //checking type of sig
		case os.Interrupt:
			cancel()
			return
		}
	}
}

func startServer(ctx context.Context) error {
	laddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", laddr) // ListenTCP is needed to set deadline
	if err != nil {
		return err
	}

	defer l.Close()

	for {
		select {
		case <-ctx.Done():
			log.Fatalln("Server stopped")
		default:
			if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
				return err
			}
			con, err := l.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}
				return err
			}
			handle(con)
			log.Println("new Client connected")
		}
	}
}

func handle(con net.Conn) {
	// Reading input
	io.WriteString(con, fmt.Sprint("Enter a number: "))
	data, _, err := bufio.NewReader(con).ReadRune()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Converting input
	num, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	//time.Sleep(time.Second*5) // It is needed to check is it really Graceful Shutdown or not :)
	io.WriteString(con, fmt.Sprintf("Square of %d is %d\n", num, num * num)) //stdout
}


