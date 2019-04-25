package main

import (
  "os";
  "syscall";
  "os/signal";

  "github.com/evalsocket/kafka-connector"
)

func main() {
  config := kafkconnector.NewConfig()
  err := config.ParseFlags()
  if err != nil {
    kafkconnector.Logger.Fatal(err)
  }
  kafkconnector.Logger.Printf("Version: %s", kafkconnector.Version)

  consumer, err := kafkconnector.NewConsumer(config)
  if err != nil {
    kafkconnector.Logger.Fatal(err)
  }
  loader, err := kafkconnector.NewLoader(config)
  if err != nil {
    kafkconnector.Logger.Fatal(err)
  }

  consumer.Batches(func(events [][]byte) error {
    data, err := kafkconnector.ParseJson(events)
    if err != nil {
      return err
    }
    rowsAffected, err := loader.Upsert(data)
    if err != nil {
      return err
    }
    kafkconnector.Logger.Printf("Affected rows: %d", rowsAffected)
    return nil
  })

  wait := make(chan os.Signal, 1)
  signal.Notify(wait, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
  <- wait

  err = consumer.Close()
  if err != nil {
    kafkconnector.Logger.Printf("Failed to close consumer: %v", err)
  }
}