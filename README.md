Kurtosis Go
===========
A Go client for the Kurtosis testing framework.

What's Kurtosis?
----------------
The Kurtosis testing framework is a system that allows you to write tests against arbitrary distributed systems. These tests could be as simple as, "create a single Elasticsearch node and make a request against it", or as complex as "spin up a network containing a database, a Kafka queue, and your custom services and run end-to-end tests on it". Each test declares the network it wants and the test logic to run, and Kurtosis handles launching the network and running the test logic. The nodes of the network are Docker containers, so anything you have a Docker image for can be used in Kurtosis.

How do I use it?
----------------
