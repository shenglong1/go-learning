# go-learning

1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

应该返回Wrap(sql.ErrNoRows), 来表示row.Scan没有数据，否则调用方无法判断"not found" 或是 "数据行存在但值为空";

考虑到sql.ErrNoRows是正常返回的一种情况(调用方可以处理)，并不需要stack信息，可以使用pkg.errors.WithMessage或自定义的Wrapper来wrap；

方便调用者使用sentinel error的形式来捕获；

代码见example_err_no_rows.go


2. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

模拟两种退出方式：ctrl+c 信号退出，某个gr报错退出；

代码见 cmd/errgroup/main.go


第八周作业：缓存

1、使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

2、写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息  , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。
