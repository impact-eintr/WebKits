package main

import (
	"log"

	"est/hst"

	lg "github.com/impact-eintr/WebKits/loadgenerator"
	"github.com/impact-eintr/WebKits/loadgenerator/gen"
)

func main() {
	// 初始化载荷发生器。
	pset := hst.PSet

	log.Printf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...",
		pset.TimeoutNS, pset.LPS, pset.DurationNS)

	gen, err := gen.NewGenerator(*pset)
	if err != nil {
		log.Fatalf("Load generator initialization failing: %s\n", err)
	}

	// 开始！
	log.Println("Start load generator...")
	gen.Start()

	printDetail := false
	// 显示结果。
	countMap := make(map[lg.RetCode]int)
	for r := range pset.ResultCh {
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			log.Printf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n",
				r.ID, r.Code, r.Msg, r.Elapse)
		}
	}

	var total int
	log.Println("RetCode Count:")
	for k, v := range countMap {
		codePlain := lg.GetRetCodePlain(k)
		log.Printf("  Code plain: %s (%d), Count: %d.\n",
			codePlain, k, v)
		total += v
	}

	log.Printf("Total: %d.\n", total)
	successCount := countMap[lg.SUCCESS]
	tps := float64(successCount) / float64(pset.DurationNS/1e9)
	log.Printf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)

}
