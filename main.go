package main

import (
    "log"
)

func main() {
        err := ReadConfig()
        if err != nil {
                log.Fatal(err)
                return
        }

        DexInit()
        BotRun()
}
