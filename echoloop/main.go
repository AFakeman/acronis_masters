package main

import (
    "bufio"
    "fmt"
    "github.com/pkg/errors"
    "log"
    "net"
    "os"
    "os/signal"
)

const SockAddr = "/tmp/echo.sock"
const BindAddressInUse = "bind: address already in use"

func receiveParams(conn net.Conn, updates chan string) {
    defer conn.Close()
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        updates <- scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Printf("Could not read all the arguments: %s", err.Error())
	} else {
        log.Println("Received new parameters")
    }
}

func printLoop(args []string, updates chan string) {
    for {
        select {
        case update := <-updates:
            args = append(args, update)
        default:
            for _, arg := range(args) {
                fmt.Println(arg)
            }
        }
    }
}

func sendParams(args []string, address string) (error) {
    conn, err := net.Dial("unix", address)
    if err != nil {
        return errors.Wrap(err, "Could not connect to existing socket")
    }

    defer conn.Close()
    writer := bufio.NewWriter(conn)
    for _, arg := range(args) {
        _, err = writer.WriteString(arg)
        if err == nil {
            _, err = writer.WriteString("\n")
        }
        if err != nil {
            return errors.Wrap(err, "Could not write argument")
        }
    }
    err = writer.Flush()
    if err != nil {
        return errors.Wrap(err, "Could not send the data over to the server")
    }
    return nil
}

func runServer(l net.Listener, updates chan string) {
    defer l.Close()
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Printf("Could not accept a new connection: %s", err.Error())
        }

        go receiveParams(conn, updates)
    }
}

func main() {
    l, err := net.Listen("unix", SockAddr)
    if err != nil {
        if err.(*net.OpError).Err.Error() == BindAddressInUse {
            log.Println("Found existing socket")
            err := sendParams(os.Args[1:], SockAddr)
            if err != nil {
                log.Fatal(err)
            } else {
                log.Println("Send params successfully")
            }
            return
        } else {
            log.Fatal("Could not listen to the socket", err.Error())
        }
    }
    defer os.RemoveAll(SockAddr)
    defer log.Println("hello")

    log.Println("Started the server")

    updates := make(chan string, 16)

    args := make([]string, len(os.Args) - 1)
    copy(args, os.Args[1:])
    go printLoop(args, updates)

    go runServer(l, updates)

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
}
