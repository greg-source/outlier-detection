package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
)

type Machine struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func validateAges(ctx *gin.Context) {
	var machines []Machine

	err := ctx.BindJSON(&machines)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Bad request"})

		return
	}

	timeUnits, err := generateTimeUnits(machines)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Bad request"})

		return
	}

	outliers := parseTimeUnits(timeUnits)

	// initializing map with key Machine.Age for quick comparison with outliers
	ageMachineMap := make(map[string][]Machine)
	for _, k := range machines {
		ageMachineMap[k.Age] = append(ageMachineMap[k.Age], k)
	}

	outlierMachines := make([]Machine, 0, len(machines)/100)
	m := new(sync.Mutex)

	for i := range outliers {
		res, exist := ageMachineMap[i]
		if exist {
			m.Lock()
			outlierMachines = append(outlierMachines, res[0])
			ageMachineMap[i] = res[1:]
			m.Unlock()
		}
	}

	ctx.JSON(200, outlierMachines)
}

// generateTimeUnits returns map, where key is unit(year, month, day, etc),
// and value is a slice of values.
func generateTimeUnits(machines []Machine) (map[string][]float64, error) {
	timeUnits := make(map[string][]float64, len(machines))

	for _, machine := range machines {
		unit := strings.Split(machine.Age, " ")
		if len(unit) != 2 {
			return nil, fmt.Errorf("cant parse %s", machine)
		}

		number, err := strconv.ParseFloat(unit[0], 64)
		if err != nil {
			return nil, fmt.Errorf("cant parse %s", machine)
		}

		timeUnits[unit[1]] = append(timeUnits[unit[1]], number)
	}

	return timeUnits, nil
}

// parseTimeUnits returns channel with detected outliers in original string format
// example "16 days", "45 years".
func parseTimeUnits(data map[string][]float64) <-chan string {
	outliers := make(chan string, len(data)/100)

	go func() {
		var wg sync.WaitGroup

		wg.Add(len(data))

		for unit, values := range data {
			unit := unit
			values := values

			go func() {
				identifyOutlierTimeUnits(unit, values, outliers)
				wg.Done()
			}()
		}

		wg.Wait()

		close(outliers)
	}()

	return outliers
}
