
# GoMLExecutor - High-performance deep-learning model runner 

This repo is designed to be compatible with Go 1.17+ versions, aiming to provide a deep-learning model prediction HTTP service. Its main focus is improving end to end execution by leveraging a caching system. As of now, the only deep-learning library supported is TensorFlow.


The main goal is to make the client manage any key generation and handle model changes seamlessly. The library can drastically speed up execution and provides both TensorFlow and cache-level performance metrics via HTTP REST API.


Apart from having functions for the client, this project also provides libraries for the associated web services. It supports multiple TensorFlow model integrations at the URI level, with `GET`, `POST` method support provided through HTTP 2.0 in full duplex mode. The service reloads any model changes automatically.


### Getting Started

To start a HTTP service with a model, create a `config.yaml` file, start the example server, and invoke a prediction. Details about this are available in the full README file in the repo.


### Caching

Major performance benefits come via trading compute with space. At the moment, this supports categorial features with fixed vocabulary and numeric features (v0.8.0 onwards). Further details are available in the README file.


### Configuration

The server has a vast range of configuration options, mentioned in detail in the README file on the repo.


### Utilization

Code snippets for Server and Client usage are provided in the README file in the repo.


### Server Endpoints

Details about the various server endpoints and their functionalities are provided in the README file.


### License and Contributions

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`. Contributions are welcome.


For further information, credits, acknowledgements, and other relevant details, please refer to the README in the repository.