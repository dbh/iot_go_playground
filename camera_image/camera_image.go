package main

import (
	"log"
	"os/exec"
	"os"
	"time"
	"net/http"
	"flag"
	"strconv"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"bufio"
	"encoding/base64"
)

const (
  WEB_DIR = "/home/pi/projects/camera_image/web"
)

var (
  c mqtt.Client
)

func getImage(refresh time.Duration, device string) {
	var tmpName = "web/"+device+".tmp"   
	var finalName = "web/"+device+".jpg"   
	for true {
		cmd := exec.Command("fswebcam", "-d", "/dev/"+device, tmpName)
		log.Printf("Running command to get image "+finalName)
		if err := cmd.Run(); err != nil {
			log.Fatal("Command finished with error: %v", err)
		}
		os.Rename(tmpName, finalName)
		log.Printf("Command finished")

		f, _ := os.Open(finalName)
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		encoded := base64.StdEncoding.EncodeToString(content)
		log.Printf("Encoded len %v", len(encoded))

		token := c.Publish("pics", 0, false, encoded)
		token.Wait()

		time.Sleep(refresh)
	}
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
  log.Printf("TOPIC: %s\n", msg.Topic())
  log.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	hostPtr := flag.String("host", "", "host name or empty for any")
	portPtr := flag.Int("port", 16000, "port number for http server")
	refreshPtr := flag.Int64("refresh", 15, "Time in seconds for camera image refresh")
	mqttPtr := flag.String("mqtt_url", "tcp://localhost:1883", "URL for mqtt");
	mqttUserPtr := flag.String("mqtt_username", "", "username for mqtt");
	mqttPasswordPtr := flag.String("mqtt_password", "", "password for mqtt");
	flag.Parse()

	log.Printf("host: %v", *hostPtr);
	log.Printf("port: %v", *portPtr);
	log.Printf("refresh: %v", *refreshPtr);
	log.Printf("mqtt_url: %v", *mqttPtr);
	log.Printf("mqtt_user: %v", *mqttUserPtr);
	log.Printf("mqtt_password: %v", *mqttPasswordPtr);

	log.Println("Setting MQTT connection")
	opts := mqtt.NewClientOptions().AddBroker(*mqttPtr)
	opts.SetClientID("camera_image")
	opts.SetDefaultPublishHandler(f)

	if (*mqttUserPtr != "") {
		opts.SetUsername(*mqttUserPtr)
	}
	if (*mqttPasswordPtr != "") {
		opts.SetPassword(*mqttPasswordPtr)
	}

	//create and start a client using the above ClientOptions
	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go getImage(time.Duration(*refreshPtr) * time.Second, "video1")

	log.Println("Starting http server")
	var hostAndPort = *hostPtr + ":" + strconv.Itoa(*portPtr)
	http.Handle("/", 
		http.StripPrefix("/", 
 		http.FileServer(http.Dir(WEB_DIR))))
	if err := http.ListenAndServe(hostAndPort, nil); err != nil {
        	log.Fatal("ListenAndServe: ", err)
	}

	
}
