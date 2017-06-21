mackerel-plugin-jmx-jolokia-ext[![Build Status](https://travis-ci.org/hajimeni/mackerel-plugin-jmx-jolokia-ext.svg?branch=master)](https://travis-ci.org/hajimeni/mackerel-plugin-jmx-jolokia-ext)
===========================

Jolokia (https://jolokia.org/) custom metrics plugin for mackerel.io agent.
This plugin extends [mackerel-plugin-jmx-jolokia](https://github.com/mackerelio/mackerel-agent-plugins/tree/master/mackerel-plugin-jmx-jolokia)

## How to install

1. Go to the [Releases Page](/hajimeni/mackerel-plugin-jmx-jolokia-ext/releases).
1. Downloa the binary for your OS.
1. Rename it to `mackerel-plugin-jmx-jolokia-ext`.


## Synopsis

```shell
mackerel-plugin-jmx-jolokia-ext [-host=<host>] [-port=<port>] [-tempfile=<tempfile>]
```

## Metrics

This plugin has all of `mackerel-plugin-jmx-jolokia` metrics.
Extend metrics, see follow.

### MemoryPool metrics

Each memory pool metrics has 4 graph (init, commited, max, used)

- Common memory pool metrics

|Name|
|----|
|Metaspace|
|Code Cache|
|Compressed Class Space|

- Each GC metrics

|Parallel GC|Serial GC|CMS|G1GC|
|----|----|----|----|
|PS Old Gen|Tenured Gen|CMS Old Gen|G1 Old Gen|
|PS Eden Space|Eden Space|Par Eden Space|G1 Eden Space|
|PS Survivor Space|Survivor Space|Par Survivor Space|G1 Survivor Space|

## GC metrics

Each GC metrics has 3 graph (Count, Time, TimePercentage)

|Parallel GC|Serial GC|CMS|G1GC|
|----|----|----|----|
|PS Scavenge|Copy|ParNew|G1 Young Generation|
|PS MarkSweep|MarkSweepCompact|ConcurrentMarkSweep|G1 Old Generation|

## Example of mackerel-agent.conf

```
[plugin.metrics.jolokia]
command = "/path/to/mackerel-plugin-jmx-jolokia-ext"
```

## Example of Metrics output

```
jmx.jolokia.gc.time.GCTimeCopy	0.000000	1497958098
jmx.jolokia.gc.time.GCTimeMarkSweepCompact	0.000000	1497958098
jmx.jolokia.gc.time_percentage.GCTimeCopy	0.000000	1497958098
jmx.jolokia.gc.time_percentage.GCTimeMarkSweepCompact	0.000000	1497958098
jmx.jolokia.memory.code_cache.CodeCacheInit	2555904.000000	1497958098
jmx.jolokia.memory.code_cache.CodeCacheCommitted	25034752.000000	1497958098
jmx.jolokia.memory.code_cache.CodeCacheMax	251658240.000000	1497958098
jmx.jolokia.memory.code_cache.CodeCacheUsed	24533120.000000	1497958098
jmx.jolokia.memory.non_heap_memory_usage.NonHeapMemoryInit	2555904.000000	1497958098
jmx.jolokia.memory.non_heap_memory_usage.NonHeapMemoryCommitted	110936064.000000	1497958098
jmx.jolokia.memory.non_heap_memory_usage.NonHeapMemoryMax	-1.000000	1497958098
jmx.jolokia.memory.non_heap_memory_usage.NonHeapMemoryUsed	108716808.000000	1497958098
jmx.jolokia.gc.count.GCCountCopy	0.000000	1497958098
jmx.jolokia.gc.count.GCCountMarkSweepCompact	0.000000	1497958098
jmx.jolokia.memory.metaspace.MetaspaceInit	0.000000	1497958098
jmx.jolokia.memory.metaspace.MetaspaceCommitted	72966144.000000	1497958098
jmx.jolokia.memory.metaspace.MetaspaceMax	-1.000000	1497958098
jmx.jolokia.memory.metaspace.MetaspaceUsed	71598216.000000	1497958098
jmx.jolokia.memory.eden_space.EdenSpaceInit	143130624.000000	1497958098
jmx.jolokia.memory.eden_space.EdenSpaceCommitted	143327232.000000	1497958098
jmx.jolokia.memory.eden_space.EdenSpaceMax	286326784.000000	1497958098
jmx.jolokia.memory.eden_space.EdenSpaceUsed	117828800.000000	1497958098
jmx.jolokia.memory.tenured_gen.TenuredGenInit	357957632.000000	1497958098
jmx.jolokia.memory.tenured_gen.TenuredGenCommitted	357957632.000000	1497958098
jmx.jolokia.memory.tenured_gen.TenuredGenMax	715849728.000000	1497958098
jmx.jolokia.memory.tenured_gen.TenuredGenUsed	67955480.000000	1497958098
jmx.jolokia.memory.heap_memory_usage.HeapMemoryInit	536870912.000000	1497958098
jmx.jolokia.memory.heap_memory_usage.HeapMemoryCommitted	519176192.000000	1497958098
jmx.jolokia.memory.heap_memory_usage.HeapMemoryMax	1037959168.000000	1497958098
jmx.jolokia.memory.heap_memory_usage.HeapMemoryUsed	185533360.000000	1497958098
jmx.jolokia.class_load.LoadedClassCount	12720.000000	1497958098
jmx.jolokia.class_load.UnloadedClassCount	33.000000	1497958098
jmx.jolokia.class_load.TotalLoadedClassCount	12753.000000	1497958098
jmx.jolokia.ops.cpu_load.ProcessCpuLoad	0.576829	1497958098
jmx.jolokia.ops.cpu_load.SystemCpuLoad	28.787879	1497958098
jmx.jolokia.memory.survivor_space.SurvivorSpaceInit	17891328.000000	1497958098
jmx.jolokia.memory.survivor_space.SurvivorSpaceCommitted	17891328.000000	1497958098
jmx.jolokia.memory.survivor_space.SurvivorSpaceMax	35782656.000000	1497958098
jmx.jolokia.memory.survivor_space.SurvivorSpaceUsed	340560.000000	1497958098
jmx.jolokia.memory.compressed_class_space.CompressedClassSpaceInit	0.000000	1497958098
jmx.jolokia.memory.compressed_class_space.CompressedClassSpaceCommitted	12935168.000000	1497958098
jmx.jolokia.memory.compressed_class_space.CompressedClassSpaceMax	1073741824.000000	1497958098
jmx.jolokia.memory.compressed_class_space.CompressedClassSpaceUsed	12624568.000000	1497958098
jmx.jolokia.thread.ThreadCount	38.000000	1497958098
```

## Example of jolokia search

```
curl -s http://127.0.0.1:8778/jolokia/search/java.lang:*| jq. 
{
  "request": {
    "mbean": "java.lang:*",
    "type": "search"
  },
  "value": [
    "java.lang:name=Metaspace,type=MemoryPool",
    "java.lang:name=Eden Space,type=MemoryPool",
    "java.lang:name=Survivor Space,type=MemoryPool",
    "java.lang:name=Copy,type=GarbageCollector",
    "java.lang:type=Runtime",
    "java.lang:type=OperatingSystem",
    "java.lang:type=Threading",
    "java.lang:name=MarkSweepCompact,type=GarbageCollector",
    "java.lang:name=Code Cache,type=MemoryPool",
    "java.lang:type=Compilation",
    "java.lang:name=Tenured Gen,type=MemoryPool",
    "java.lang:name=CodeCacheManager,type=MemoryManager",
    "java.lang:name=Compressed Class Space,type=MemoryPool",
    "java.lang:type=Memory",
    "java.lang:type=ClassLoading",
    "java.lang:name=Metaspace Manager,type=MemoryManager"
  ],
  "timestamp": 1498007188,
  "status": 200
}
```