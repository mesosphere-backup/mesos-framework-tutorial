# mesos-framework-tutorial in Go

### 1) [Deploy a Mesos Cluster](https://github.com/mesosphere/playa-mesos)

### 2) Setup the Environment
#### Install Git and Mercurial
```sh
sudo apt-get install -y git
sudo apt-get install -y mercurial
```

#### Install Go
Within the VM created above [install Go](https://golang.org/doc/install) and [setup its required workspace](https://golang.org/doc/code.html).

### 3) Create a new framework

#### Get the template code

```sh
$ mkdir -p $GOPATH/src/github.com/mesosphere
$ cd $GOPATH/src/github.com/mesosphere
$ git clone git@github.com:mesosphere/mesos-framework-tutorial.git
$ go get ./...
$ cd mesos-framework-tutorial/
$ git checkout -b tutorial origin/tutorial
```

At the point you should have all the tutorial code and be in the 'tutorial' branch.  This tutorial steps through commits on that branch adding framework functionality as we go.

#### Run the code
At any point from here on, you should be able to compile and run the code.  Both the scheduler and the executor must be compiled as follows:

```sh
$ cd $GOPATH/src/github.com/mesosphere/mesos-framework-tutorial
$ go build -o example_scheduler
$ cd $GOPATH/src/github.com/mesosphere/mesos-framework-tutorial/executor
$ go build -o example_executor
```

The example can then be run at any commit in the tutorial branch with:

```sh
$ cd $GOPATH/src/github.com/mesosphere/mesos-framework-tutorial
$ ./example_scheduler --master=127.0.0.1:5050 --executor="/home/vagrant/code/go/src/github.com/mesosphere/mesos-framework-tutorial/executor/example_executor" --logtostderr=true
```

#### Framework template

The absolute minimum requirements for a framework consist of 3 components: scheduler, executor, and file server.

The scheduler receives resource offers from mesos and makes decisions about what tasks should consume which resources.

The executor knows how to run the tasks that the schedulder launches.  A more detailed explanation is [available here](http://mesos.apache.org/documentation/latest/mesos-architecture/).

The file server is necessary in order to provide Mesos a location, from which it can retrieve the executor binary.  This third component is usually not explicitly called out as a requirement of a framework.  It is not technically a component of the framework but is a necessity for the end-to-end functioning of a framework.  Of course, multiple frameworks can share the same server.

In the [first commit](https://github.com/mesosphere/mesos-framework-tutorial/commit/aae4f846a6dd7e5e0fba2d737dc82718ddde9e2b) the three components are in their respective directories.

'main.go' initializes the configuration of all three components and packages them into a configuration object.  This object is passed to the MesosSchedulerDriver, which is then started.

'scheduler/example_scheduler.go' implements the required scheduler interface and logs all calls from Mesos.

'executor/example_executor.go' compiles to an executable binary which is capable of hosting tasks.  It implements the executor interface and for the most part just logs calls from Mesos.  The exception is the LaunchTask method which makes status updates regarding tasks, but does not actually do any work.

If you compile and run the example code at this point you will see that the scheduler receives one resource offer from Mesos and then appears to block.  By not accepting the resource offer the scheduler has implicitly rejected the offer.  No tasks are launched.  A configurable timeout will eventually occur and the resource will again be offered to the scheduler.  The output should like this:

```sh
...
I0713 19:03:42.775536   25174 scheduler.go:446] Framework registered with ID=20150713-1...
I0713 19:03:42.775962   25174 example_scheduler.go:48] Scheduler Registered with Master...
I0713 19:03:42.776181   25174 utils.go:32] Received Offer <20150713-...> with cpus=2 mem=1000
```

#### Launch Tasks

In order to do something which is at least moderately interesting, let's start accepting a few offers from Mesos and launch some tasks.  If we look at the second commit (PUT COMMIT SHA1 here) we see that we now iterate across the offers provided by Mesos and launch tasks until we run out of resources.

The executor launches the tasks, and reports status to Mesos indicating that the tasks are finished.  This frees the resources and they are offered to the scheduler again.  This loop continues endlessly.  As long as the scheduler process doesn't crash, a long running distributed service has now been completed.  'example_executor.go' indicates where real work should be done in it's 'LaunchTask' method.

When the code is run we should see output which indicates that tasks are running:

```sh
I0713 19:06:52.967857   25228 utils.go:32] Received Offer <20150713-...> with cpus=2 mem=1000
I0713 19:06:52.967939   25228 example_scheduler.go:90] Prepared task: go-task-1 with offer 20150713-... for launch
I0713 19:06:52.967973   25228 example_scheduler.go:90] Prepared task: go-task-2 with offer 20150713-... for launch
I0713 19:06:52.968075   25228 example_scheduler.go:96] Launching  2 tasks for offer 20150713-...
I0713 19:06:54.173174   25228 example_scheduler.go:103] Status update: task 1  is in state  TASK_RUNNING
I0713 19:06:54.174417   25228 example_scheduler.go:103] Status update: task 2  is in state  TASK_RUNNING
I0713 19:06:54.176197   25228 example_scheduler.go:103] Status update: task 1  is in state  TASK_FINISHED
I0713 19:06:54.178064   25228 example_scheduler.go:103] Status update: task 2  is in state  TASK_FINISHED
...
```

#### Image Processing Example
In order to show data and metadata flowing back and forth between the scheduler and the executor, we can implement a simple batch image processing framework.  In this case we'll invert the images, as show below.

![Original image](https://raw.githubusercontent.com/mesosphere/mesos-framework-tutorial/tutorial/original.jpg?token=AAinR_TyrX7bO_bT7H4QJRMtj5Be-jYAks5VrZPSwA%3D%3D)

![Inverted image](https://raw.githubusercontent.com/mesosphere/mesos-framework-tutorial/tutorial/inverted.jpg?token=AAinR7nLI_1CmA9ImzaiARIniQt0K1lyks5VrZQRwA%3D%3D)

The framework should take [a list of image URLs](https://github.com/mesosphere/mesos-framework-tutorial/blob/tutorial/images), assign each URL to one task, and collect the results of the image processing.  The changes necessary to accomplish this can be seen in the [third commit of this tutorial](https://github.com/mesosphere/mesos-framework-tutorial/commit/8c91479951b1a5bd6467e548f1ebd861f34ba547).

In order for each task to know which image it should process we need to [encapsulate the URL in the task here](https://github.com/mesosphere/mesos-framework-tutorial/blob/tutorial/scheduler/example_scheduler.go#L104).  The executor then [reads this information on the other side](https://github.com/mesosphere/mesos-framework-tutorial/blob/tutorial/executor/example_executor.go#L65).  The executor then [processes the image and uploads it](https://github.com/mesosphere/mesos-framework-tutorial/blob/tutorial/executor/example_executor.go#L72-87) to the same HTTP server that previously served the executor binary.

The HTTP server was modified to allow for image uploads, by [registering an upload handler function](https://github.com/mesosphere/mesos-framework-tutorial/blob/tutorial/executor/example_executor.go#L72-87).  It saves images to the directory from which the Vagrant VM was launched.  This [directory mapping was added](https://github.com/mesosphere/mesos-framework-tutorial/blob/tutorial/Vagrantfile#L63) in the VagrantFile.


