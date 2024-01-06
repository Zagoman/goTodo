# Simple todo API application

The purpose of this api is to learn how the http server works in Go, and to practice good API patterns in go.

#### DB structure

```mermaid
    classDiagram;
    Todo: +int ID;
    Todo: +string Task;
    Todo: +timestamp CreatedAt;
```