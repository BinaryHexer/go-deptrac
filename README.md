# Package Dependency Checker for Golang

> :construction: This repository is currently WIP 

go-deptrac is a static analysis tool that helps to enforce rules for dependencies between software layers in your golang projects.
It is heavily inspired from [deptrac].

## Getting Started

### The Depfile

The depfile is the config file for this tool, in this file you define (mainly) three things:
                                              
- The location of your source code.
- The layers of your application.
- The allowed dependencies between your layers.

Here is an example depfile:

```yaml
# depfile.yaml
paths:
  - ./src
exclude_files:
  - .*test.*
layers:
  - name: Controller
    collectors:
      - type: directory
        regex: .*Controller.*
  - name: Repository
    collectors:
      - type: directory
        regex: .*Repository.*
  - name: Service
    collectors:
      - type: directory
        regex: .*Service.*
ruleset:
  Controller:
    - Service
  Service:
    - Repository
  Repository: ~
```

#### Explanation
In the first section, paths, you declare where deptrac should look for your code. As this is an array of directories, you can 
specify multiple locations.

With the exclude_files section, you can specify one or more regular expressions for files that should be excluded, the most 
common being probably anything containing the "test" word in the path.

We defined three layers in the example: Controller, Repository and Service. Deptrac is using so called collectors to group 
packages into layers (in this case by the name of the directory).

The ruleset section defines, how these layers may or may not depend on other layers. In the example, every package of the 
Controller layer may depend on packages that reside in the Service layer, and packages in the Service layer may depend on packages 
in the Repository layer.

Packages in the Repository layer may NOT depend on any packages in other layers. The ruleset acts as a whitelist, therefore the 
Repository layer rules can be omitted, however explicitly stating that the layer may not depend on other layers is more 
declarative.

If a package in the Repository layer uses a package in the Service layer, deptrac will recognize the dependency and raises a violation 
for this case. The same counts if a Service layer package uses a Controller layer package.

## Installation

```
go get -u gopkg.in/BinaryHexer/go-deptrac.v1
```

go-deptrac was only tested on Linux and should also work on OS X. Probably it doesn't work well on Windows.

## Running

To run the tool, simply pass the path of the depfile:

```shell script
go-deptrac ./examples/simple-mvc/depfile.yaml
```

## Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License. Please read the [LICENSE](./LICENSE) for the complete details.

[references]: #
[deptrac]: https://github.com/sensiolabs-de/deptrac
