package watcher

import (
	"context"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	ctx context.Context
	*fsnotify.Watcher
}

func New(ctx context.Context, paths ...string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			return nil, err
		}
	}

	return &Watcher{
		ctx:     ctx,
		Watcher: watcher,
	}, nil
}

func (watcher *Watcher) Start() error {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("events closed")
			}
			log.Printf("%s %s\n", event.Name, event.Op)
		case err, ok := <-watcher.Errors:
			if !ok {
				return fmt.Errorf("errors closed")
			}
			log.Println("error:", err)
		case <-watcher.ctx.Done():
			return nil
		}
	}
}

func (watcher *Watcher) Stop() error {
	return watcher.Close()
}
