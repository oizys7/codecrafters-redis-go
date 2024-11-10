# redis-in-go
这是一个学习项目。它是 codecrafters 提供的 Build your own Redis 挑战的解决方案。
在此之前我并未学过 go 语言，这个项目会让我熟悉 go 语言的语法和使用 go 语言构建一个项目会遇到的问题，并且学习到 redis 客户端是如何工作的。我希望尽可能地实现 redis 客户端的所有功能，并且随时检查并修改我的代码，提升自己的编程能力和组件抽象能力。

目前以及实现的功能有：

- set、get 指令
- 支持 set 指令的 px、ex 过期时间
- hset、hget、hgetall 指令
- rdb 文件的简单解析
