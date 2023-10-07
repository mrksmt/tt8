# Клиенты с ограничением количества выполняемых заданий

## Клиент с ограничителем на базе Redis DB

expected output

```text
since start: 1ms
items processed: 3
retry after: 0s

since start: 104ms
items processed: 3
retry after: 0s

// в момент 125ms один из элементов уже закончил процессинг
// и на следующем шаге вместо 2 элементов в работу можно взять 3

since start: 203ms
items processed: 3 
retry after: 0s

// к моменту 250ms второй элемент закончил процессинг
// на освободившееся место можно взять в работу один элемент

since start: 302ms
items processed: 1
retry after: 0s

// по мере освобождения очереди далее в работу будет передаваться или один элемент или ни одного

since start: 402ms
items processed: 0
retry after: 25ms

since start: 502ms
items processed: 1
retry after: 0s

since start: 602ms
items processed: 0
retry after: 24ms

since start: 704ms
items processed: 1
retry after: 0s

since start: 804ms
items processed: 0
retry after: 25ms

since start: 904ms
items processed: 1
retry after: 0s

since start: 1.003s
items processed: 0
retry after: 26ms

since start: 1.103s
items processed: 1
retry after: 0s

since start: 1.203s
items processed: 0
retry after: 25ms

since start: 1.304s
items processed: 1
retry after: 0s

since start: 1.403s
items processed: 0
retry after: 26ms

since start: 1.502s
items processed: 1
retry after: 0s

since start: 1.603s
items processed: 0
retry after: 24ms

since start: 1.703s
items processed: 1
retry after: 0s

since start: 1.802s
items processed: 0
retry after: 25ms

since start: 1.902s
items processed: 1
retry after: 0s
```

```
since start: 0s
items processed: 3
retry after: 0s

since start: 101ms
items processed: 3
retry after: 0s

since start: 201ms
items processed: 2
retry after: 799ms

since start: 301ms
items processed: 0
retry after: 699ms

since start: 401ms
items processed: 0
retry after: 599ms

since start: 501ms
items processed: 0
retry after: 499ms

since start: 601ms
items processed: 0
retry after: 399ms

since start: 700ms
items processed: 0
retry after: 300ms

since start: 801ms
items processed: 0
retry after: 199ms

since start: 901ms
items processed: 0
retry after: 99ms

since start: 1s
items processed: 3
retry after: 0s

since start: 1.101s
items processed: 3
retry after: 0s

since start: 1.201s
items processed: 2
retry after: 799ms

since start: 1.3s
items processed: 0
retry after: 700ms

since start: 1.4s
items processed: 0
retry after: 600ms

since start: 1.501s
items processed: 0
retry after: 499ms

since start: 1.601s
items processed: 0
retry after: 399ms

since start: 1.701s
items processed: 0
retry after: 299ms

since start: 1.8s
items processed: 0
retry after: 200ms

since start: 1.9s
items processed: 0
retry after: 100ms

since start: 2.001s
items processed: 3
retry after: 0s
```
