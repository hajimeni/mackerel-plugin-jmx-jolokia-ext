package mpjmxjolokia

import (
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"net/http"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
	"strings"
)

var logger = logging.GetLogger("metrics.plugin.jmx-jolokia")

// JmxJolokiaPlugin mackerel plugin for Jolokia
type JmxJolokiaPlugin struct {
	Target   string
	Tempfile string
	GcBeanNames     []string
	MemoryBeanNames []string
}

// JmxJolokiaResponse response for Jolokia
type JmxJolokiaResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     map[string]interface{}
	Error     string
}

type JmxJolokiaValueResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     uint64
	Error     string

}

var graphdef = map[string]mp.Graphs{
	"jmx.jolokia.memory.heap_memory_usage": {
		Label: "Jmx HeapMemoryUsage",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "HeapMemoryInit", Label: "init", Diff: false, Type: "uint64"},
			{Name: "HeapMemoryCommitted", Label: "committed", Diff: false, Type: "uint64"},
			{Name: "HeapMemoryMax", Label: "max", Diff: false, Type: "uint64"},
			{Name: "HeapMemoryUsed", Label: "used", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.memory.non_heap_memory_usage": {
		Label: "Jmx NonHeapMemoryUsage",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "NonHeapMemoryInit", Label: "init", Diff: false, Type: "uint64"},
			{Name: "NonHeapMemoryCommitted", Label: "committed", Diff: false, Type: "uint64"},
			{Name: "NonHeapMemoryMax", Label: "max", Diff: false, Type: "uint64"},
			{Name: "NonHeapMemoryUsed", Label: "used", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.class_load": {
		Label: "Jmx ClassLoading",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "LoadedClassCount", Label: "loaded", Diff: false, Type: "uint64"},
			{Name: "UnloadedClassCount", Label: "unloaded", Diff: false, Type: "uint64"},
			{Name: "TotalLoadedClassCount", Label: "total", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.thread": {
		Label: "Jmx Threading",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ThreadCount", Label: "thread", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.ops.cpu_load": {
		Label: "Jmx CpuLoad",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "ProcessCpuLoad", Label: "process", Diff: false, Type: "float64", Scale: 100},
			{Name: "SystemCpuLoad", Label: "system", Diff: false, Type: "float64", Scale: 100},
		},
	},
}


var rep = regexp.MustCompile(`java.lang.name=(.+),type=(.+)`)

// FetchMetrics interface for mackerelplugin
func (j JmxJolokiaPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	if err := j.fetchMemory(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchClassLoad(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchThread(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchOperatingSystem(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchMemoryPool(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchGarbageCollector(stat); err != nil {
		logger.Warningf(err.Error())
	}

	return stat, nil
}


func (j JmxJolokiaPlugin) fetchGarbageCollector(stat map[string]interface{}) error {
	for _, v := range j.GcBeanNames {
		group := rep.FindStringSubmatch(v)
		beanName := group[1]
		nameBase := strings.Replace(beanName, " ", "", -1)

		cresp, err := j.executeGetRequestValue(strings.Replace(v, " ", "%20", -1) + "/CollectionCount")
		fmt.Printf("%v\n", cresp)
		if err != nil {
			return err
		}
		stat["GCCount" + nameBase] = cresp.Value

		tresp, err := j.executeGetRequestValue(strings.Replace(v, " ", "%20", -1) + "/CollectionTime")
		fmt.Printf("%v\n", tresp)
		if err != nil {
			return err
		}
		stat["GCTime" + nameBase] = tresp.Value

	}
	return nil

}

func (j JmxJolokiaPlugin) fetchMemoryPool(stat map[string]interface{}) error {
	for _, v := range j.MemoryBeanNames {
		resp, err := j.executeGetRequest(strings.Replace(v, " ", "%20", -1) + "/Usage")
		if err != nil {
			return err
		}

		group := rep.FindStringSubmatch(v)
		beanName := group[1]
		nameBase := strings.Replace(beanName, " ", "", -1 )

		mem := resp.Value
		stat[nameBase + "Init"] = mem["init"]
		stat[nameBase + "Committed"] = mem["committed"]
		stat[nameBase + "Max"] = mem["max"]
		stat[nameBase + "Used"] = mem["used"]

	}
	return nil
}

func (j JmxJolokiaPlugin) fetchMemory(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=Memory")
	if err != nil {
		return err
	}
	heap := resp.Value["HeapMemoryUsage"].(map[string]interface{})
	stat["HeapMemoryInit"] = heap["init"]
	stat["HeapMemoryCommitted"] = heap["committed"]
	stat["HeapMemoryMax"] = heap["max"]
	stat["HeapMemoryUsed"] = heap["used"]

	nonHeap := resp.Value["NonHeapMemoryUsage"].(map[string]interface{})
	stat["NonHeapMemoryInit"] = nonHeap["init"]
	stat["NonHeapMemoryCommitted"] = nonHeap["committed"]
	stat["NonHeapMemoryMax"] = nonHeap["max"]
	stat["NonHeapMemoryUsed"] = nonHeap["used"]

	return nil
}

func (j JmxJolokiaPlugin) fetchClassLoad(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=ClassLoading")
	if err != nil {
		return err
	}
	stat["LoadedClassCount"] = resp.Value["LoadedClassCount"]
	stat["UnloadedClassCount"] = resp.Value["UnloadedClassCount"]
	stat["TotalLoadedClassCount"] = resp.Value["TotalLoadedClassCount"]

	return nil
}

func (j JmxJolokiaPlugin) fetchThread(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=Threading")
	if err != nil {
		return err
	}
	stat["ThreadCount"] = resp.Value["ThreadCount"]

	return nil
}

func (j JmxJolokiaPlugin) fetchOperatingSystem(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=OperatingSystem")
	if err != nil {
		return err
	}
	stat["ProcessCpuLoad"] = resp.Value["ProcessCpuLoad"]
	stat["SystemCpuLoad"] = resp.Value["SystemCpuLoad"]

	return nil
}

func (j JmxJolokiaPlugin) executeGetRequest(mbean string) (*JmxJolokiaResponse, error) {
	resp, err := http.Get(j.Target + mbean)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respJ JmxJolokiaResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		return nil, err
	}
	return &respJ, nil
}

func (j JmxJolokiaPlugin) executeGetRequestValue(mbean string) (*JmxJolokiaValueResponse, error) {
	resp, err := http.Get(j.Target + mbean)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respJ JmxJolokiaValueResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		return nil, err
	}
	return &respJ, nil
}

// GraphDefinition interface for mackerelplugin
func (j JmxJolokiaPlugin) GraphDefinition() map[string]mp.Graphs {

	gcCountMetrics := []mp.Metrics{}
	gcTimeMetrics := []mp.Metrics{}
	gcTimePercentageMetrics := []mp.Metrics{}
	for _, v := range j.GcBeanNames {
		group := rep.FindStringSubmatch(v)
		beanName := group[1]
		nameBase := strings.Replace(beanName, " ", "", -1)

		gcCountMetrics = append(gcCountMetrics, mp.Metrics{ Name: "GCCount" + nameBase, Label: beanName, Diff: true, })
		gcTimeMetrics = append(gcTimeMetrics, mp.Metrics{ Name: "GCTime" + nameBase, Label: beanName, Diff: true, })
		gcTimePercentageMetrics = append(gcTimePercentageMetrics, mp.Metrics{ Name: "GCTime" + nameBase, Label: beanName, Diff: true, Scale: (100.0 / 60), })
	}

	var extGraphdef = map[string]mp.Graphs{
		"jmx.jolokia.gc.count":  {
			Label: "Jmx GC counts",
			Unit:  "bytes",
			Metrics: gcCountMetrics,
		},
		"jmx.jolokia.gc.time":  {
			Label: "Jmx GC time(sec)",
			Unit:  "float",
			Metrics: gcTimeMetrics,
		},
		"jmx.jolokia.gc.time_percentage":  {
			Label: "Jmx GC time percentage",
			Unit:  "percentage",
			Metrics: gcTimePercentageMetrics,
		},
	}

	for _, v := range j.MemoryBeanNames {
		group := rep.FindStringSubmatch(v)
		beanName := group[1]
		nameBase := strings.Replace(beanName, " ", "", -1 )
		extGraphdef["jmx.jolokia.memory." + strings.ToLower(strings.Replace(beanName, " ", "_", -1))] = mp.Graphs{
			Label: "JMX " + beanName,
			Unit: "bytes",
			Metrics: []mp.Metrics{
				{Name: nameBase + "Init", Label: "init", Diff: false, Type: "uint64"},
				{Name: nameBase + "Committed", Label: "committed", Diff: false, Type: "uint64"},
				{Name: nameBase + "Max", Label: "max", Diff: false, Type: "uint64"},
				{Name: nameBase + "Used", Label: "used", Diff: false, Type: "uint64"},
			},
		}
	}
	for k, v := range graphdef {
		extGraphdef[k] = v
	}
	return extGraphdef
}

func (j JmxJolokiaPlugin) describeExtGraphDefs(jolokiaTarget string) ([]string, []string) {
	resp, err := http.Get(jolokiaTarget)
	if err != nil {
		logger.Warningf("Jolokia search Error. use default")
		return nil, nil
	}
	defer resp.Body.Close()
	var respJ JmxJolokiaSearchResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		logger.Warningf("Jolokia search Response Error. use default")
		return nil, nil
	}

	/*
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
	 */

	gcBeanNames := []string{}
	memoryBeanNames := []string{}
	for _, v := range respJ.Value {
		group := rep.FindStringSubmatch(v)
		if len(group) == 3 {
			beanType := group[2]

			switch beanType {
			case "GarbageCollector":
				gcBeanNames = append(gcBeanNames, v)
			case "MemoryPool":
				memoryBeanNames = append(memoryBeanNames, v)
			}
		}
 	}
	return gcBeanNames, memoryBeanNames
}

type JmxJolokiaSearchResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     []string
	Error     string
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "8778", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var jmxJolokia JmxJolokiaPlugin
	jmxJolokia.Target = fmt.Sprintf("http://%s:%s/jolokia/read/", *optHost, *optPort)

	// Check GC Type
	var searchTarget = fmt.Sprintf("http://%s:%s/jolokia/search/java.lang:*", *optHost, *optPort)

	gc, memory := jmxJolokia.describeExtGraphDefs(searchTarget)
	jmxJolokia.GcBeanNames = gc
	jmxJolokia.MemoryBeanNames = memory

	helper := mp.NewMackerelPlugin(jmxJolokia)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-jmx-jolokia-ext-%s-%s", *optHost, *optPort))
	}
	helper.Run()
}
