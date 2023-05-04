package command

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/damirm/links-warehouse/internal/fetcher"
	"github.com/damirm/links-warehouse/internal/http"
	"github.com/damirm/links-warehouse/internal/parser"
	"github.com/damirm/links-warehouse/internal/postgres"
	"github.com/damirm/links-warehouse/internal/processor"
	"github.com/damirm/links-warehouse/internal/telegram"
	"github.com/damirm/links-warehouse/internal/warehouse"
	"github.com/damirm/links-warehouse/internal/worker"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/yaml.v3"
)

type options struct {
	config            string
	testDataDirectory string
}

type Command struct {
	o *options
}

func NewWarehouseCommand() *Command {
	return &Command{
		o: &options{},
	}
}

func (c *Command) ExportFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.o.config, "config", "configs/config.yaml", "Path to configuration file")
	fs.StringVar(&c.o.testDataDirectory, "test-data-directory", "", "Path to files containing links")
}

type config struct {
	Http      *http.Config
	Telegram  *telegram.Config
	Database  *postgres.Config
	Worker    *worker.Config
	Processor *processor.Config
}

func (c *Command) Run() error {
	conf, err := readConfig(c.o.config)
	if err != nil {
		log.Printf("failed to read file: %s", c.o.config)
		return err
	}

	w := worker.NewWorker(conf.Worker)

	db, err := postgres.Connect(conf.Database)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return err
	}

	s, err := postgres.InitStorage(context.Background(), db, conf.Database)
	if err != nil {
		log.Printf("failed to initialize storage: %v", err)
		return err
	}

	service := warehouse.NewWarehouseService(s)

	if c.o.testDataDirectory != "" {
		if err := processor.PopulateQueueFromDirectory(c.o.testDataDirectory, service); err != nil {
			return err
		}
	}

	f := &fetcher.HttpFetcher{}
	p := &parser.HabrParser{}

	lp := processor.NewLinkProcessor(service, w, f, p, conf.Processor)

	lp.Start()
	w.Start()

	handleSignals(func(s os.Signal) {
		lp.Stop()
		w.Stop()
	}, os.Interrupt, os.Kill)

	lp.Join()
	w.Join()

	return nil
}

func handleSignals(cb func(os.Signal), signals ...os.Signal) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, signals...)

	go func() {
		for {
			s := <-sc
			cb(s)
		}
	}()
}

func readConfig(configPath string) (*config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	log.Printf("loaded config: %s", data)

	expanded := []byte(os.ExpandEnv(string(data)))
	conf := &config{}

	if err := yaml.Unmarshal(expanded, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
