# Dogpool

Dogpool is a Go package that provides a simple and efficient task scheduling and execution system based on GORM.

## Features

- Task scheduling
- Worker pool for task execution
- MySQL database integration using GORM
- Customizable task execution functions
- Error handling and logging

## Installation

To use Dogpool in your Go project, run:
```
go get github.com/chinmayrelkar/dogpool
```


## Usage

### Setting up the database

Dogpool uses GORM with MySQL. Set up your database connection like this:

```

gormDb, err := gorm.Open(mysql.New(mysql.Config{
    DriverName: "mysql",
    DSN:        os.Getenv("MYSQL_DSN"),
}), &gorm.Config{})
```


### Creating a Worker


```
worker, exit := dogpool.NewWorker(context.TODO(), gormDb)
```


### Scheduling Tasks

To schedule a task, use the `ScheduleTask` method:

```
err := worker.ScheduleTask("task_name", []byte("task_data"), time.Now().Add(5*time.Minute))
```

### Registering Task Handlers

Register task handlers using the `RegisterTaskHandler` method:

```
worker.RegisterTaskHandler("task_name", func(ctx context.Context, taskData []byte) error {
    // Task execution logic here
    return nil
})
```

### Starting the Worker

Start the worker to begin processing tasks:


```
go worker.Start()
```

### Stopping the Worker

To stop the worker gracefully:


```
worker.Stop()
```

## Example

Here's a complete example of how to use Dogpool:

```
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/chinmayrelkar/dogpool"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    // Set up database connection
    gormDb, err := gorm.Open(mysql.New(mysql.Config{
        DriverName: "mysql",
        DSN:        os.Getenv("MYSQL_DSN"),
    }), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Create a new worker
    worker, exit := dogpool.NewWorker(context.TODO(), gormDb)

    // Register task handler
    worker.RegisterTaskHandler("example_task", func(ctx context.Context, taskData []byte) error {
        log.Printf("Executing task with data: %s", string(taskData))
        return nil
    })

    // Start the worker
    go worker.Start()

    // Schedule a task
    err = worker.ScheduleTask("example_task", []byte("Hello, Dogpool!"), time.Now().Add(10*time.Second))
    if err != nil {
        log.Printf("Failed to schedule task: %v", err)
    }

    // Wait for a while to allow task execution
    time.Sleep(15 * time.Second)

    // Stop the worker
    exit()
}
```

## Contributing

Contributions to Dogpool are welcome! Please feel free to submit a Pull Request.

