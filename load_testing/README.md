# Overview
This folder contains the scripts which can be used to run the load tests using [k6](https://k6.io/). The load tests are really basic running 150 virtual concurrent users each making each request every second with the test running for 60 seconds. Random domain names are used. As this is all runnig on one machine, if testing locally you can run into performance bottlenecks and machine specific issues like file descriptor limits.

# Running the Tests
To run the tests you need node.js installed. Once node.js is installed, you can run the tests using:

```bash
npm install
npm test
```

**NOTE:** The tests assume you have the service running locally. You can follow the instructions in the main readme for how to run the service.

# Notes
* You may run into file descriptor issues if running everything locally. Started seeing some failures due to file descriptor limits at about 150 VCU's
