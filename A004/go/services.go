package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Services struct {
	Services []Service
	mu       sync.RWMutex
}

type Service struct {
	Name    string        `json:"name"`
	URL     string        `json:"url"`
	Up      bool          `json:"-"`
	Latency time.Duration `json:"-"`
}

func ReadServices() (*Services, error) {
	b, err := os.ReadFile("../services.json")
	if err != nil {
		return nil, fmt.Errorf("reading services file: %w", err)
	}

	var config struct {
		Services []Service `json:"services"`
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("parsing services file: %w", err)
	}

	return &Services{
		Services: config.Services,
		mu:       sync.RWMutex{},
	}, nil
}

func (s *Services) Query(ctx context.Context) {
	// Initial query
	err := s.innerQuery(ctx)
	if err != nil {
		log.Printf("error querying services: %v", err)
	}

	t := time.NewTicker(time.Minute)
	for range t.C {
		err := s.innerQuery(ctx)
		if err != nil {
			log.Printf("error querying services: %v", err)
		}
	}
}

func (s *Services) innerQuery(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Services {
		service := s.Services[i]

		log.Printf("querying service %d/%s", i, service.Name)

		startTime := time.Now()
		req, err := http.NewRequestWithContext(ctx, "GET", service.URL, nil)
		if err != nil {
			return fmt.Errorf("querying service %d/%s: %w", i, service.Name, err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("querying service %d/%s: %w", i, service.Name, err)
		}

		isUp := resp.StatusCode >= 200 && resp.StatusCode <= 299
		latency := time.Since(startTime).Round(time.Millisecond)

		log.Printf("queried service %v. up: %v, latency: %v", service.Name, isUp, latency)

		s.Services[i] = Service{
			Name:    service.Name,
			URL:     service.URL,
			Up:      isUp,
			Latency: latency,
		}
	}

	return nil
}
