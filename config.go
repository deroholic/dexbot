package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
)

type ConfigType struct {
    Token	string `json:"token"`
    Prefix	string `json:"prefix"`
    Owner	string `json:"owner"`
    Channel	string `json:"channel"`
    Daemon	string `json:"daemon"`
}

var config ConfigType

func ReadConfig() error {
    file, err := ioutil.ReadFile("./config.json")
    if err != nil {
        log.Fatal("Can't open config file:", err)
        return err
    }
    err = json.Unmarshal(file, &config)
    if err != nil {
        log.Fatal("Config error:", err)
        return err
    }

    return nil
}
