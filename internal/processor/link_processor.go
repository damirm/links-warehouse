package processor

import (
	"context"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/damirm/links-warehouse/internal/fetcher"
	"github.com/damirm/links-warehouse/internal/parser"
	"github.com/damirm/links-warehouse/internal/warehouse"
	"github.com/damirm/links-warehouse/internal/worker"
)

type Config struct {
	PickInterval time.Duration `yaml:"pick-interval"`
}

type LinkProcessor struct {
	service *warehouse.WarehouseService
	worker  *worker.Worker
	fetcher fetcher.Fetcher
	parser  parser.Parser
	config  *Config

	start sync.Once
	stop  sync.Once
	quit  chan struct{}

	processed uint64
}

func NewLinkProcessor(s *warehouse.WarehouseService, w *worker.Worker, f fetcher.Fetcher, p parser.Parser, c *Config) *LinkProcessor {
	quit := make(chan struct{})
	return &LinkProcessor{
		service: s,
		worker:  w,
		fetcher: f,
		parser:  p,
		config:  c,
		quit:    quit,
	}
}

func (p *LinkProcessor) Start() {
	p.start.Do(func() {
		p.watch()
	})
}

func (p *LinkProcessor) Stop() {
	p.stop.Do(func() {
		close(p.quit)
	})
}

func (p *LinkProcessor) Join() {
	<-p.quit
}

func (p *LinkProcessor) ProcessedJobs() uint64 {
	return atomic.LoadUint64(&p.processed)
}

func (p *LinkProcessor) watch() {
	urls := make(chan *url.URL)
	ticker := time.NewTicker(p.config.PickInterval)

	log.Printf("start watching links queue, pick new links every %d seconds", p.config.PickInterval/time.Second)

	go func() {
		for {
			select {
			case <-p.quit:
				log.Printf("stop watching links queue")
				close(urls)
				ticker.Stop()
				return
			case <-ticker.C:
				batch, err := p.pickBatch(10)
				for _, u := range batch {
					log.Printf("picked url: %s %v", u, err)
					urls <- u
				}
			}
		}
	}()

	go func() {
		for u := range urls {
			func(u *url.URL) {
				p.worker.Run(func() error {
					log.Printf("start processing %s", u)
					defer log.Printf("finished processing %s", u)
					defer atomic.AddUint64(&p.processed, 1)
					return p.process(u)
				})
			}(u)
		}
	}()
}

func (p *LinkProcessor) pickBatch(size int) (res []*url.URL, err error) {
	for len(res) < size {
		u, err := p.service.NextProcessingLink()
		if err != nil {
			log.Printf("failed to pick url: %v", err)
			return res, err
		}

		if u != nil {
			res = append(res, u)
		} else {
			return res, nil
		}
	}

	return res, nil
}

func (p *LinkProcessor) process(u *url.URL) error {
	ctx := context.Background()

	log.Printf("fetching url: %s", u)
	body, err := p.fetcher.Fetch(ctx, u)
	if err != nil {
		log.Printf("failed to fetch url: %v", err)
		return err
	}

	log.Printf("parsing url: %s", u)
	link, err := p.parser.Parse(u, body)
	if err != nil {
		log.Printf("failed to parse url: %v", err)
		return err
	}

	log.Printf("saving url: %s", u)
	return p.service.FinishLinkProcessing(ctx, link)
}
