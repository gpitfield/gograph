package gograph

import (
	"testing"
)

func checkEqual(t *testing.T, map1, map2 map[string]float64) {
	margin := 0.01
	for key, val := range map1 {
		if val > map2[key]+margin || val < map2[key]-margin {
			t.Fatalf("incorrect pageranks for ex0:\n%v\n\nshould be:\n%v", map1, map2)
		}
		t.Logf("key %v is correct at %v\n", key, val)
	}
}

func TestGraph(t *testing.T) {
	// test inputs taken from http://www.sirgroane.net/google-page-rank/
	ex0 := PageRank(&[]Edge{
		Edge{"A", "B", 0.5},
		Edge{"B", "A", 1.0},
	}, 0.85, 1.0)
	ex0Result := map[string]float64{"A": 0.5, "B": 0.5}
	checkEqual(t, *ex0, ex0Result)

	ex1 := PageRank(&[]Edge{
		Edge{"A", "B", 1.0},
		Edge{"A", "C", 1.0},
		Edge{"B", "C", 1.0},
		Edge{"C", "A", 1.0},
		Edge{"D", "C", 1.0},
	}, 0.85, 4.0)
	ex1Result := map[string]float64{"A": 1.49, "B": 0.78, "C": 1.58, "D": 0.15}
	checkEqual(t, *ex1, ex1Result)
}
