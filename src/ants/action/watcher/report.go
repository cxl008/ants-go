package watcher

import (
	"ants/action"
	"ants/crawler"
	"ants/node"
	"log"
	"time"
)

/*
* report request result of scrapy
*	*	container list of report
*	*	dead loop for report result to master
**/

const (
	REPORT_STATUS_RUNNING = iota
	REPORT_STATUS_PAUSE
	REPORT_STATUS_STOP
	REPORT_STATUS_STOPED
)

type Reporter struct {
	Status      int
	ResultQuene *crawler.ResultQuene
	Node        *node.Node
	rpcClient   action.RpcClientAnts
	distributer action.Watcher
}

func NewReporter(node *node.Node, rpcClient action.RpcClientAnts, resultQuene *crawler.ResultQuene, distributer action.Watcher) *Reporter {
	return &Reporter{REPORT_STATUS_STOPED, resultQuene, node, rpcClient, distributer}
}

func (this *Reporter) Start() {
	if this.Status == REPORT_STATUS_RUNNING {
		return
	}
	for {
		if this.IsStop() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	this.Status = REPORT_STATUS_RUNNING
	go this.Run()
}

func (this *Reporter) IsPause() bool {
	return this.Status == REPORT_STATUS_PAUSE
}

func (this *Reporter) Pause() {
	if this.Status == REPORT_STATUS_RUNNING {
		this.Status = REPORT_STATUS_PAUSE
	}
}

func (this *Reporter) Unpause() {
	if this.Status == REPORT_STATUS_PAUSE {
		this.Status = REPORT_STATUS_RUNNING
	}
}

// set it stop,if deap loop see this sign,make it stoped
func (this *Reporter) Stop() {
	this.Status = REPORT_STATUS_STOP
}

// if it is stoped
func (this *Reporter) IsStop() bool {
	return this.Status == REPORT_STATUS_STOPED
}

// pop result quene
// if scraped new request  set node name local node name
// send it to master
func (this *Reporter) Run() {
	log.Println("start reporter")
	for {
		if this.IsPause() {
			time.Sleep(1 * time.Second)
			continue
		}
		if this.Status != REPORT_STATUS_RUNNING {
			this.Status = REPORT_STATUS_STOPED
			break
		}
		result := this.ResultQuene.Pop()
		if result == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		nodeName := this.Node.NodeInfo.Name
		if result.ScrapedRequests != nil {
			for _, request := range result.ScrapedRequests {
				request.NodeName = nodeName
			}
		}
		log.Println(result.Request.SpiderName, ":report request to master:", result.Request.GoRequest.URL.String())
		if this.Node.IsMasterNode() {
			this.Node.ReportToMaster(result)
			this.JudgeAndStopNode()
		} else {
			this.rpcClient.ReportResult(this.Node.GetMasterName(), result)
		}
	}
	log.Println("stop reporter")
}

// stop action is start by local report action so put it here
func (this *Reporter) JudgeAndStopNode() {
	if !this.Node.IsStop() {
		return
	}
	for _, nodeInfo := range this.Node.GetAllNode() {
		if this.Node.IsMe(nodeInfo.Name) {
			this.Node.StopCrawl()
			this.Stop()
			this.distributer.Stop()
		} else {
			this.rpcClient.StopNode(nodeInfo.Name)
		}
	}
}
