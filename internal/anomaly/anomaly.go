package anomaly

import (
	"fmt"
	"time"

	"github.com/christian-gama/produgit/internal/data"
)

// Config represents the configuration for the anomaly command.
type Config struct {
	startDate time.Time
	endDate   time.Time
	quantity  int32
	input     string
	authors   []string
}

// NewConfig creates a new Config.
func NewConfig(
	startDate, endDate time.Time,
	quantity int32,
	input string,
	authors []string,
) (*Config, error) {
	if endDate.IsZero() {
		endDate = time.Now()
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("Quantity must be greater than 0")
	}

	if len(authors) == 0 {
		return nil, fmt.Errorf("No authors provided")
	}

	cfg := &Config{
		startDate: startDate,
		endDate:   endDate,
		quantity:  quantity,
		input:     input,
		authors:   authors,
	}

	return cfg, nil
}

// Anomaly prints the anomalies found in the logs.
func Anomaly(l *data.Logs, config *Config) error {
	logs, err := data.Filter(
		l,
		data.WithDate(config.startDate, config.endDate),
		data.WithAuthors(config.authors),
	)
	if err != nil {
		return err
	}

	found := false
	for _, log := range logs.Logs {
		if log.Plus > config.quantity {
			if !found {
				fmt.Printf("%-6s - %-15s - %s\n", "Plus", "Author", "Path")
			}
			fmt.Printf(
				"%-6d - %-15s - %s\n",
				log.Plus,
				fmt.Sprintf("%.15s", log.Author),
				log.Path,
			)
			found = true
		}
	}

	if !found {
		fmt.Println("No anomalies found")
	}

	return nil
}
