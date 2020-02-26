# iot_go_playground
IOT, GoLang, MQTT playground

This is a playground for playing with IOT, GoLang, MQTT, or whatever else strikes my fancy.

Initially, this project makes use of the following in order to periodically sample images from one (or more) cameras, at a given frequency, and send the images to an MQTT topic.  A corresponding process subscribes to corresonding topics and reconstitutes the images, for viewing on a centralized host. 

## This project currently utilizes the following
* Linux
* fswebcam for capturing images from USB web cams
* GoLang
* MQTT 

## Future
* It is likely that I'll add more IOT sensors of various types
* More nodes for sensing
* Visualizations
