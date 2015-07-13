# mesos-framework-tutorial

## OS X

### 1) [Deploy a Mesos Cluster](https://github.com/mesosphere/playa-mesos)

### 2) Setup the Environment
#### Install Git and Mercurial
```sh
sudo apt-get install -y git
sudo apt-get install -y mercurial
```

#### Install Go
Within the VM created above [install Go](https://golang.org/doc/install) and [setup its required workspace](https://golang.org/doc/code.html).

#### [Install mesos-go](https://github.com/mesosphere/mesos-go)
Verify that you are able to run the example framework.

### 3) Create a new framework

This tutorial is an example which was cloned from the mesos-go example and modified.

```sh
$ mkdir -p $GOPATH/src/github.com/mesosphere
$ cd $GOPATH/src/github.com/mesosphere
$ git clone git@github.com:mesosphere/mesos-framework-tutorial.git
```
