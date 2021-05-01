# blinkt-controlled-mq
Reads MQTT message to configure Blinkt piHat in Raspberry Pi.

# Foreword
```
First of all, this is just a proof of concept and should not have any other use than studying or testing
```

Being a PoC, this is highly specific for the technologies and platforms I wanted to test, meaning that, while this could be ported to other hosts than Netlify, using different LEDs than Blinkt or different MQTT than HiveMQ, my goal was to test all of them together, so keep in mind that the configuration and some design decisions may be specific for those choices.

Someday, of course, I will try to expand this app, and any help is always much appreciated!

# Overview
Using a public webpage (with serveless functions) hosted on Netlify, you can configure the LEDs on a Blinkt piHat and publish you choices to a public HiveMQ server. On the other end, a Golang app running on a Pi Zero with a Blinkt piHat installed is subscribed to that same public HiveMQ server. When a new message is sent to the correct topic, the app will update the Blinkt configuration.

This repository is only the **Golang application described above**. If you are interested in the webpage, check my other repository: https://github.com/otaviokr/web-remote-control-bedroom

If you are looking for a simpler version, with the webpage running in the same Raspberry Pi where Blinkt is installed, check: https://github.com/otaviokr/blinkt-web-ui

# How to run
```
Keep in mind you will need to run this program with root privileges due Raspberry Pi GPIO limitations
```

The easiest to get it running is to clone this repo and run the main source file.

```
git clone https://github.com/otaviokr/blinkt-controlled-mq.git
cd blinkt-controlled-mq
sudo -E go run main.go
```

# How it works

You will be presented a page with 8 black circles, a color picker, a slider and a button:

- Each circle represents one of the LEDs in Blinkt. To change its color, select the color and click on the circle. The webpage has no feedback, so the color of the circle initially may differ from its current real state;
- Use the color picker to select the next color you want to configure the LEDs with. Black turns the LED off;
- The slider defines the brightness level of all LEDs. You cannot assign brightness individually to the LEDs, the last value before submitting the configuration will be the one used on Blinkt. I recommend using levels around 10 because the LEDs are very bright!
- After you've configured all the LEDs you want and set the desired brightness level, press the button to submit the values to Blinkt. The array should be updated shortly.

# Dependencies

All interface with GPIO is handled by [periph.io](https://periph.io/). Logging is done with [Sirupsen's Logrus](https://github.com/sirupsen/logrus).

Look'n'Feel is powered by [Semantic-UI](https://semantic-ui.com/) and the [Superhero](https://github.com/semantic-ui-forest/forest-themes/blob/master/dist/bootswatch/v3/semantic.superhero.min.css) theme. They are not required, of course, but nobody wants to use an ugly page...
