### ants-go 
open source, restful, distributed crawler engine  
### why
I wrote a crawler engine named [ants](https://github.com/wcong/ants) in python base on [scrapy](https://github.com/scrapy/scrapy). But sometimes, dynamic language is chaos.
So I start to write it in a compile language.  
### gopath
``` shell
export GOPATH=PATH/TO/ants-go
```
### requirement
``` shell
go get github.com/PuerkitoBio/goquery
go get github.com/go-sql-driver/mysql
```
### install

``` shell
go install src/ants/ants/ants.go
```

### run

``` shell
cd bin
./ants
```

#### cluster in one computer
to test cluster in one computer,you can run it from different port in different terminal

one node,use the default port tcp 8300 http 8200

``` shell
cd bin
./ants
```

the other node set tcp port and http port

``` shell
cd bin
./ants -tcp 9300 -http 9200
```
#### flags
there are some flags you can set,check out the help message

``` shell
./ants -h
./ants -help
```

### Customize spider
1.	go to *src/spiders*
2.	write your spiders follow the example *deap_loop_spider.go* or go to the [spider page](./SPIDER.md)
3.	add you spider to spiderMap,follow the example in *LoadAllSpiders* in *load_all_spider.go*
4.	install again

