---
author: "wendisx"
last-updated: "2026/01/15 am"
todo:
  - Improve the concept document
  - Improve the program design
  - Write the Api manual
---

# Puzzle Concept & Reference 

- [Pre-Concept](#pre_concept)
- [Database](#database)
- [Log](#log)
- [Router](#router)
- [Config](#config)
- [Cli](#cli) 
- [Integration](#integration)
- [Reference](#reference)

---

## <a id="pre_concept">Pre-Concept</a>

So why should I write this simple framework and what is the motivation? When I was doing web applications, I found that it was not easy to use
There is no difference between different web frameworks. In fact, there is nothing worthy of in-depth study in the framework layer.
Most of the time, you just need to use them quickly to complete your own ideas. They are just tools, not so original tools. These tools
Sometimes it's too fragmented for my taste, which results in me having to do it again every time I rebuild web applications.
These things, it's boring. Why not make it templated to reduce unnecessary duplication of behavior? This is how I wrote `puzzle`
The reason is that this is a personal tool package collection. If you think this can help you, please fork the repository to customize it.
righteous.

`puzzle` encapsulates basic web development components and supports the integration of many similar tools (core goals). Note: `puzzle` does not
It is not a web framework, but is used to shield the differences between frameworks (to the greatest extent). The core is to provide a unified interface, so for
Some frameworks may need to provide their own implementation.

## <a id="database">Database</a>

Database mainly encapsulates **sql** and **nosql** common APIs and unifies the **repo** layer calling form. In order to maximize flexibility,
Database only abstracts data sources, not data operations. Data operations are passed from the **repo** layer call. Specifically:
Constructed a local database instance that can be directly operated, but there is no need to pay specific attention to the database interaction details, repo
The layer only needs to pass the database operation and the actual type it wants to get, and the instance can return the specified data collection. The
At present, all the differences between sql and nosql are explicitly distinguished, so that they can be built independently, but different sql and nosql are provided.
Provides a certain degree of unified calling of the API.

You need to **explicitly add the database type** you need, which means you need to know a lot of details, which is beneficial.

The database driver can be found [here](https://go.dev/wiki/SQLDrivers).

Integrated Database:
- [Mysql](https://www.mysql.com)
  - [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [Sqlite](https://sqlite.org)
  - [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

## <a id="log">Log</a>

Log provides the basic implementation of structured and unstructured logs, based on the standard `slog` package. The design concept is basically the same as
Slog is consistent, a simple c/s-like interactive log processing. Currently, plain text (unstructured colored) logs are implemented--
clog, combined with the internal palette package for log coloring. Generally, you only need to simply call the exported standard log
API is enough. If you need a dedicated logger, you can build a new logger, which is useful when you need to integrate a log analysis system.
Very helpful.

In the future, structured logs will be implemented, mainly used for log collection and analysis.

## <a id="router">Router</a>

Router is simply designed according to the concept of `network` sending network packets, and only provides the necessary abstractions and
The default implementation of `Echo` (actually all modules will be implemented based on this framework). Router does not contain automatic behavior,
This means that like Database, you need to explicitly handle and define everything, and eventually you can integrate according to these, but
In fact, `puzzle` will not do this for you. This will be explained later at the Cli level. 

Router is a network graph structure, rather than a simple recursive tree.

## <a id="config">Config</a>

Config is mainly used to load configuration files and environment variables. The DataDict data structure is introduced in `puzzle`, using
Used to store all possible configurations and values, based on the encapsulation of the `ttlcache` library. Any environment variables and related configuration
many times only needs to stay in memory for a short time, which means that DataDict will clear these within a certain period of time, because
Encapsulated from ttlcache, this is easy to do, and many operations are supported. DataDict can be used as a non-special need
to replace `map`, which is very common in non-large web applications. Config does not explicitly and automatically override
any built-in configuration items, which means that if you need to implement the behavior of environment variables overriding the default configuration, you need to implement it explicitly, but this is not a troublesome matter.

Config's current filepath loading mechanism **relies on the executed ospath**, which means that if `cwd` is different from the actual relative path Consistency will make the program unable to find the specified file. In fact, Cli also has this problem.

## <a id="cli">Cli</a>

The CLI parses a specified command file to build a set of commands into a data dictionary. 
This dictionary is then used to construct a Cobra command invocation structure, achieving 
command integration. The CLI has default internal commands, and all commands are ultimately 
integrated into a single root command (the application), thus achieving CLI integration. This 
is for future research and preliminary design of multi-service launchers and reusable 
templates.

## <a id="integration">Integration</a>

Integrate third-party API interfaces. 

| third name | website | progress | usage | 
|:----------:|:-------:|:--------:|:-----:|
| echoSwagger | https://github.com/swaggo/echo-swagger | default | `server.EchoServer.MountSwagRoute()` |

## <a id="reference">Reference</a>

### PUZZLE API

The reference is generated using `go doc`. For more details about `go doc`, please see [here](https://go.dev/doc/comment). To browse the document locally, follow the steps below:

```bash
git clone git@github.com:wendisx/puzzle.git --depth=1
cd ./puzzle
go doc -http
```

### REST API 

You can use the [Swagger tool](https://github.com/swaggo/swag) to generate integrated Open API documentation and then mount 
the specified routes on the server. This requires the integrated web framework to handle 
it automatically, which is no different from the default integrated echo Swagger operation.
