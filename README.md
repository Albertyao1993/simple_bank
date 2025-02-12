# Simple Bank System 

## init project

```golang
go mod init github.com/Albertyao1993/simple_bank
go mod tidy
```

---
在数据库中 注意Trasanction 的数据保护，以及 dead lock 问题，
- SELECT * FROM AAA for UPDATE  