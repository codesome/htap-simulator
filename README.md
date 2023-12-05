# HTAP Simulator

![Architecture](assets/img/htap-simulator-architecture.png)

```bash
# All commands run in parallel, in the given order
$ go run htap-brain/*
$ python3 load-generator/write.py
$ python3 load-generator/read.py
```