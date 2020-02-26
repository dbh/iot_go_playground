package main

import (
	"log"
	"flag"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"encoding/base64"
	"os"
)

const (
  WEB_DIR = "/home/pi/projects/camera_image/web"
)

var (
  c mqtt.Client
)


var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
  log.Printf("TOPIC: %s\n", msg.Topic())
  log.Printf("MSG: %s\n", msg.Payload())
}

var fileMsgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
  log.Printf("PICS TOPIC: %s\n", msg.Topic())
  log.Printf("PICS Payload len: %v", len(msg.Payload()))
  decoded, err := base64.StdEncoding.DecodeString(string(msg.Payload()))
  if err != nil {
    log.Fatal(err)
  }
  f, err := os.Create("/var/www/wp_davidbharrison.com/poc/" + "image.jpg")
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  if _, err := f.Write(decoded); err != nil {
    log.Fatal(err)
  }
  if err := f.Sync(); err != nil {
    log.Fatal(err)
  }
  fi, _ := f.Stat()
  log.Printf("file size: %v\n", fi.Size())
}

func main() {
	hostPtr := flag.String("host", "", "host name or empty for any")
	portPtr := flag.Int("port", 16000, "port number for http server")
	refreshPtr := flag.Int64("refresh", 15, "Time in seconds for camera image refresh")
	mqttPtr := flag.String("mqtt_url", "tcp://localhost:1883", "URL for mqtt")
	mqttUserPtr := flag.String("mqtt_username", "", "username for mqtt")
	mqttPasswordPtr := flag.String("mqtt_password", "", "password for mqtt")

	flag.Parse()

	log.Printf("host: %v", *hostPtr);
	log.Printf("port: %v", *portPtr);
	log.Printf("refresh: %v", *refreshPtr);
	log.Printf("mqtt_url: %v", *mqttPtr);

	log.Println("Setting MQTT connection")
	opts := mqtt.NewClientOptions().AddBroker(*mqttPtr)
	if (*mqttUserPtr != "") {
		opts.SetUsername(*mqttUserPtr)
	}
	if (*mqttPasswordPtr != "") {
		opts.SetPassword(*mqttPasswordPtr)
	}
	opts.SetClientID("go-sub-pics")
	opts.SetDefaultPublishHandler(f)
	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	
	token := c.Subscribe("pics", 0, fileMsgHandler)
	log.Println("Waiting around")	
	token.Wait() 
	if token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Println("select to wait forever");
	select{}

	log.Println("Unsubscribing")
	if token := c.Unsubscribe("pics"); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	c.Disconnect(250)	
}
