# Top-K Heavykeeper

Native Go implementation of Top-K Heavykeeper algorithm, efficient, high precision and fast algorithm for finding Top-K elephant flows.

Based on work of Junzhi Gong, Tong Yang, Haowei Zhang, and Hao Li: [HeavyKeeper: An Accurate Algorithm
for Finding Top-k Elephant Flows](https://www.usenix.org/system/files/conference/atc18/atc18-gong.pdf)

## Install

```
go get github.com/migotom/heavykeeper
```

## Usage

Sample usage of heavykeeper:

```go
  // keep track on top 25 elephant flows
  k := uint32(25)

  // array width, higher value means more precise results and higher memory consumption
  width := uint32(2048)

  // amount of arrays
  depth := uint32(5)

  // probability decay
  decay := 0.9

  heavykeeper := heavykeeper.New(k, width, depth, decay)
  heavykeeper.Add("some_key_1")

  // ... adding more keys ...

  heavykeeper.Add("some_key_xx")

  // at any time we can lookup into Top-K table
  if frequentKey := heavykeeper.Query("some_key_100"); frequentKey {
    fmt.Println("one of top 25 keys")
  }

  // print list Top K values
  for _, e := range heavykeeper.List() {
		fmt.Println(e.Item, e.Count)
	}
```

## Credits

Application was developed by Tomasz Kolaj and is licensed under Apache License Version 2.0. Please reports bugs at https://github.com/migotom/heavykeeper/issues.
