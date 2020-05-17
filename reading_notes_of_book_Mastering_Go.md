# Mastering Go
- 英文版读书笔记
- 书中 go version 1.9.1 / 2017
## Chapter 1
- go 的一些优点
    - 有垃圾回收器
    - 没有前置处理器
    - 使用静态链接，也就是编译出来的程序迁移到其他机器就可以运行，而不需要外部依赖程序。但这样编译出来的文件也比较大。
    - 支持 Unicode，方便处理多国语言
    - 支持命令行脚本方式运行
- 前置处理器 preprocessor，是指在编译前先对源代码按照一定规则进行修改，这可能会修改了源代码的逻辑。java 有前置处理器。
- go 的一些缺点
    - 没有内置支持面向对象编程
    - 没有 C 运行快
- godoc 工具，类似 javadoc，用于生成、浏览代码文档，还可以 godoc fmt Printf 查看具体的函数文档或者启动一个 web server，在浏览器中查看文档
- 编译命令 go build xxx.go
- 安装程序 go install xxx.go，会在 GOPATH\bin\ 下生成一个 xxx 文件
- go run 命令实际也是先编译生成可执行文件，然后执行，执行完毕之后会把可执行文件删除
- 删除一个下载的 package 之前，需要先执行如下的 clean 命令
    ```shell
    go clean -i -v -x github.com/xxx/xxx/xxx
    rm GOPATH\src\github.com\xxx\xxx\xxx
    ```
- Unix 系统中的标准输入、输入、错误设备
    - Unix 系统中一切都是文件，而文件以文件以一个正整数作为描述符，这样比用文件路径能更方便的访问文件。
    - 每个 Unix 系统中都有 3 个一直打开着的文件
        - /dev/stdin /dev/stdout /dev/stderr
        - 对应的文件描述符是0、1、2
    - MacOS 中，标准输入在 /dev/fd/0
    - Debian 中，标准输入可以同时通过如下路径访问 /dev/fd/0，/dev/pts/0
    - go 中通过 os.Stdin os.Stdout os.Stderr 访问
    - 运行程序时，2>&1 表示将标准错误打印到标准输出，>/dev/null 2>&1 表示将标准输入和标准错误都重定向到 /dev/null 设备
- 程序中获取运行时的命令行参数
  
    - main方法中， args := os.Args，定义一个数组来获取，第一个值是程序名称
- 日志
    - go 中使用 Unix 系统日志服务，一般的日志在 /var/log 文件夹中，当然这个可以配置 
    - MaxOs 中日志服务进程是 syslogd，其他大部分 Linux 系统的日志服务进程是 rsyslogd
    - 日志服务的配置文件是 /etc/syslog.conf 或 /etc/rsyslog.conf
    - 日志级别 debug info notice warning err crit alert emerg
    - 日志工具 logging facility
        - auth authpriv cron daemon kern lpr mail mark news syslog user UUCP local0 local1 ... local7
    - go 中的日志工具在 log 包中，在 Windows 系统中没有实现
    - log.Fatal() 会在记录日志后使得程序退出 exit status 1
    - log.Panic() 会得到更多日志信息，也会使得程序退出，exit status 2
- go 中非常重视错误，专门定义了 error 这个错误类型
    - 创建 error 变量 err := errors.New("error info")
    - 获取 error 中的信息 err.Error()，返回字符串
## Chapter 2
### 编译器
- go tool compile hello.go 将得到 .o 文件
- go tool compile -pack hello.go 将得到 .a 文件
- go tool compile -S hello.go 将会得到更多的信息
### 垃圾回收器
- go 的垃圾回收器跟 java 的有几分相似，都是程序运行过程中多线程并发执行，先标记再清除，但没有年轻代、老年代的概念，也没有压缩
- go 中提供了获取垃圾回收器数据的方法
    ```go
    var men runtime.MenStats
    // 读取新的垃圾回收器数据
    runtime.ReadMenStats(&men)
    men.Alloc
    men.TotalAlloc
    men.HeadAlloc
    men.NumGC
    ```
    - 运行时增加 GODEBUG=gctrace=1 参数可以打印 gc log
    ```shell
    GODEBUG=gctrace=1 go run hello.go
    ```
    - 垃圾回收器使用的是 tricolor mark-and-sweep 算法，这个算法将内存区域分为黑－灰－白三种颜色集合。mark 过程：垃圾回收开始时，所以对象都是白色，垃圾回收器访问所有的 root object，并将其染为灰色。然后垃圾回收器就会选一个灰色的对象，染为黑色，并扫描看其是否有指向白色对象的指针，如果有，这些白色的对象就会被染为灰色。直到灰色的对象都被染为黑色扫描完毕，还是白色的对象就是可以回收的对象。如果灰色的对象变得不可达，那么将会再下一次 gc 时进行回收。
    - 会 stop the world
    - 新的对象，指针变动的对象，都会是黑色的，通过 mutator 的 write barrier 方法实现。不影响已经标记完成的白色对象，这方案使得垃圾回收器能够并发执行，减少 stw 造成的延迟，但是 write barrier 方法执行的时间变成一种代价
    - channels 在不可达时，即使没有 closed，也会被回收
    - 代码中可以通过 runtime.GC() 显式调用 gc，但会阻塞
### unsafe code
- unsafe code 一般是在处理指针，关系到内存安全。
- 需要导入 Unsafe 包
- 如果指针运算错误，读取一个不在范围内的指针内容，go 无法捕获这样的错误，返回值将会无法预测
### go 调用 C 代码
- C 代码不多的话，将 C 代码写到 Go 代码文件中，C 代码要写在注释里，wtf...
- C 代码也可以写到单独的文件了，编译成库，再在 Go 中引入，引入代码还是要写在注释里
- Go 传给 C 函数的参数，要手工释放内存
### C 代码调用 Go
- 将 Go 编译成库文件，会自动生成头文件
- C 代码也将会依赖 Go 生成的文件，比如需要传 int 型给 Go 函数，C 代码里需要定义 Goint 类型变量
### defer 关键字
- defer 关键字使得被修饰的方法延迟到包含 defer 关键字的方法 return 时， 按后进先出顺序执行
- 如下函数将会先输出 after for loop，说明 fmt.Println() 被延迟执行，再输出 2，最后输出 1，说明使用后进先出顺序执行被 defer 修饰的方法。因为 i 是通过参数传给 fmt.Println() 方法的，如果是匿名方法中直接使用 i 而不是使用参数传递的话，相当于这个匿名方法在 for 循环结束之后再取 i 的值，得 3
    ```go
    func dfunc(){
        for i := 1; i < 3; i++{
            defer fmt.Println("i=", i)
        }
        fmt.Println("after for loop")
    }
    ```
- 通常用于文件读写功能中，用于写关闭文件操作，这样不用到方法末尾再写，避免忘记
### panic() and recover()
- 这类似于其他语言的 trow exception try...catch 语句
    ```go
    func a(){
        fmt.Println("in func a()")
        defer func(){
            if r := recover(); r != nil{
                fmt.Println("func a() recover")
            }
        }()
        panic("func a() panic")
        fmt.Println("this line will never be execute")
    }
    ```
- 当然也可以只 panic() 而不 recover()
### strace 工具
- strace 工具是 Linux 系统中的工具，MacOs也没有
- strace 工具用于追踪系统调用和信号
- strace 的 -c 参数，能够打印各个系统调用的好是比例、调用次数，发生错误次数等
    ```shell
    strace ls
    strace -c ls
    ```
### dtrace 工具
- Unix中，类似于 strace 工具
- MacOs 和 FreeBSD 中也有一个版本
- MacOs 中有 dtruss 命令行工具
- 需要 root 权限
### 获取系统信息
- 通过 runtime 包获取‘
    ```go
    runtime.Version
    runtime.Compiler()
    runtime.GOARCH
    runtime.NumCPU()
    runtime.NumGoroutine()
    ```
### 汇编器 assembler
- 通过命令行能看到 go 汇编后的代码
    ```shell
    GOOS=windows GOARCH=amd64 go tool compile -S hello.go
    GOOS=linux GOARCH=386 go build -gcflags -S hello.go
    ```
- GOOS 可选值有 android、darwin、freebsd、openbsd、solaris、linux、windows 等
- GOARCH 可选值有 386、amd、amd64、arm64、arm、s390 等等
### Node Trees
- go tool compile -W hello.go
### go build -x
- 获取多的构建信息
    ```shell
    go build -x hello.go
    ```
### 一般的 go 代码建议
- 对于 error，要么打到日志里处理要么返回调用方
- 使用接口定义行为，而不是数据和数据结构
- 使用 io.Reader 和 io.Writer 接口使得程序具有更好的扩展性
- 尽量不使用指针作为函数参数
- error 是个变量，不是字符串
- 不要在生产环境测试
- 对于不了解的 go 特性，先测试再使用
- 不要怕犯错，勇于探索
## chapter 3 基本数据结构

### array
- 过
### slice
- 常用 slice 少用 array
- slice 也可以是多维的
- copy(dst, src)，长度不同时，按小的来
- 排序 sort.Slice()，new in go 1.8
### map
- 不要使用 float 型做 key，因为浮点型数有 == 比较问题
- 对零值 nil 的 map 执行添加键值会报错，但 len() delete() range 操作不会报错

- 判断一个 key 存不存在，不能简单通过 map["key"] 的 value 返回值判断，要通过第二个返回参数判断，因为当 key 不存在时，返回的 value 是零值，比如 int 型返回 0 时就无法判断。

  ```go
  _, ok = map["key"]
  if ok {
      fmt.Println("key exists")
  }
  ```

### constant

- 常量是在编译器就确定值的

- go 使用 Boolean、string、number 作为常量的类型，能够在处理时有更多的可拓展性

- 定义常量时，如果不指定类型，使用时能够自动类型转换

  ```go
  const s1 = 123
  var v1 float32 = s1 * 12
  ```

- iota 略

### pointer

- 使用 & 取得变量的内存地址赋给指针变量
- 指针变量通过 * 取得内存地址的变量的值
- go 中 string 是指类型，c 中是指针类型

### date time

- epoch time  是指 1970 年 1 月 1日至今经过的秒数

  ```go
  //BasicUse shows use basic of date time
  func BasicUse()  {
  	fmt.Println("epoch time:", time.Now.Unix())
  	t := time.Now()
  	fmt.Println(t, t.Format(time.RFC3339))
  	fmt.Println(t.Weekday, t.Date, t.Month, t.Year)
  
  	// sleep 2 seconds
  	time.Sleep(time.Second * 2)
  
  	t1 := time.Now()
  	fmt.Println("time diff:", t1.Sub(t))
  }
  ```

- 格式化

  | go 格式             | 其他语言格式 | 说明                     |
  | ------------------- | ------------ | ------------------------ |
  | 2006                | yyyy         | 4 位年份表示法           |
  | 01                  | MM           | 2 位月份表示法           |
  | 02                  | dd           | 2 位日表示法             |
  | 03                  | hh           | 12 小时制小时表示法      |
  | 15                  | HH           | 24小时制小时表示法       |
  | 04                  | mm           | 2 位分钟表示法           |
  | 05                  | ss           | 2 位秒表示法             |
  | -0700               |              | 时区                     |
  | Mon                 |              | 3 位缩写字母的星期表示法 |
  | Monday              |              | 全单词的星期表示法       |
  | 1                   |              | 月                       |
  | 2                   |              | 日                       |
  | 3、3pm、03AM        |              | 时                       |
  | 4                   |              | 分                       |
  | 5                   |              | 秒                       |
  | 06                  |              | 2 位年份                 |
  | -07、-0700、Z0700   |              | 时区                     |
  | Z07:00、-07:00、MST |              | 时区                     |

  

  ```go
  func formating() {
  	t := time.Now()
  	cnDateFormatPattern := "2006/01/02 15:04:05 -0700"
  	fmt.Println("CN Date Formating:", t.Format(cnDateFormatPattern))
      // output--> CN Date Formating:2020/05/17 12:49:56 +0800
  }
  ```

  

​	