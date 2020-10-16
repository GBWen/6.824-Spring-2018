package mapreduce

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
)

func doMap(
	jobName string, // the name of the MapReduce job
	mapTask int,    // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(filename string, contents string) []KeyValue,
) {
	//
	// doMap manages one map task: it should read one of the input files
	// (inFile), call the user-defined map function (mapF) for that file's
	// contents, and partition mapF's output into nReduce intermediate files.
	//
	// There is one intermediate file per reduce task. The file name
	// includes both the map task number and the reduce task number. Use
	// the filename generated by reduceName(jobName, mapTask, r)
	// as the intermediate file for reduce task r. Call ihash() (see
	// below) on each key, mod nReduce, to pick r for a key/value pair.
	//
	// mapF() is the map function provided by the application. The first
	// argument should be the input file name, though the map function
	// typically ignores it. The second argument should be the entire
	// input file contents. mapF() returns a slice containing the
	// key/value pairs for reduce; see common.go for the definition of
	// KeyValue.
	//
	// Look at Go's ioutil and os packages for functions to read
	// and write files.
	//
	// Coming up with a scheme for how to format the key/value pairs on
	// disk can be tricky, especially when taking into account that both
	// keys and values could contain newlines, quotes, and any other
	// character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!
	//
	// Your code here (Part I).
	//

	/**
		doMap管理一个map任务：它应该读取一个输入文件（inFile），为该文件的内容调用用户定义的map函数（mapF），并将mapF的输出分区为nReduce中间文件。
	    每个reduce任务有一个中间文件。文件名包括映射任务编号和reduce任务编号。使用reduceName（jobName，mapTask，r）生成的文件名作为reduce task r的中间文件。在每个键mod nReduce上调用ihash（）（见下文），为一个键/值对选择r。
		mapF（）是应用程序提供的map函数。第一个参数应该是输入文件名，尽管map函数通常会忽略它。第二个参数应该是整个输入文件的内容。 mapF（）返回包含reduce的键/值对的切片;请参阅common.go以了解KeyValue的定义。
		查看Go的ioutil和os包以获取读写文件的功能。
		提出如何格式化磁盘上的键/值对的方案可能很棘手，尤其是考虑到键和值都可以包含换行符，引号和您能想到的任何其他字符。
		通常用于将数据序列化为另一端可以正确重建的字节流的一种格式是JSON。您不需要使用JSON，但由于reduce tasks *的输出必须是JSON，因此在这里熟悉它可能会很有用。您可以使用下面的注释代码将数据结构写为文件的JSON字符串。相应的解码功能可以在common_reduce.go中找到。
	*/

	log.Printf("doMap: job name = %s, input file = %s, map task id = %d, nReduce = %d\n",
		jobName, inFile, mapTask, nReduce)

	contents, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal("doMap: read file ", inFile, "error!", err)
	}
	log.Println("doMap: read file", inFile, "success.")

	keyValues := mapF(inFile, string(contents))

	encoders := make([]*json.Encoder, nReduce)
	for i := 0; i < nReduce; i++ {
		fileName := reduceName(jobName, mapTask, i)
		filePtr, err := os.Create(fileName)
		if err != nil {
			log.Fatal("doMap: create file", fileName, "error!", err)
		}
		defer filePtr.Close()
		encoders[i] = json.NewEncoder(filePtr)
	}

	for _, kv := range keyValues {
		key := kv.Key
		reduceIndex := ihash(key) % nReduce
		err := encoders[reduceIndex].Encode(kv)
		if err != nil {
			log.Fatal("doMap:", kv, "write file error!", err)
		}
	}
	log.Println("doMap: write files success.")
}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}
