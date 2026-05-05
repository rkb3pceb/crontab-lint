package heatmap

import (
	"testing"
)

func TestBuild_EveryMinute(t *testing.T) {
	hm, err := Build("* * * * * echo hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for dow := 0; dow < 7; dow++ {
		for hour := 0; hour < 24; hour++ {
			if hm.Grid[dow][hour] != 60 {
				t.Errorf("Grid[%d][%d] = %d, want 60", dow, hour, hm.Grid[dow][hour])
			}
		}
	}
	if hm.MaxHits != 60 {
		t.Errorf("MaxHits = %d, want 60", hm.MaxHits)
	}
}

func TestBuild_HourlyJob(t *testing.T) {
	hm, err := Build("0 * * * * echo hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for dow := 0; dow < 7; dow++ {
		for hour := 0; hour < 24; hour++ {
			if hm.Grid[dow][hour] != 1 {
				t.Errorf("Grid[%d][%d] = %d, want 1", dow, hour, hm.Grid[dow][hour])
			}
		}
	}
}

func TestBuild_WeekdayOnly(t *testing.T) {
	// Run at noon on weekdays only (Mon-Fri = 1-5).
	hm, err := Build("0 12 * * 1-5 echo hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for dow := 0; dow < 7; dow++ {
		for hour := 0; hour < 24; hour++ {
			want := 0
			if dow >= 1 && dow <= 5 && hour == 12 {
				want = 1
			}
			if hm.Grid[dow][hour] != want {
				t.Errorf("Grid[%d][%d] = %d, want %d", dow, hour, hm.Grid[dow][hour], want)
			}
		}
	}
}

func TestBuild_InvalidExpression(t *testing.T) {
	_, err := Build("bad expression")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestBuild_Cells_NonZero(t *testing.T) {
	hm, err := Build("0 9 * * 1 echo hi") // 09:00 every Monday
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cells := hm.Cells()
	if len(cells) != 1 {
		t.Fatalf("len(Cells) = %d, want 1", len(cells))
	}
	if cells[0].Day != 1 || cells[0].Hour != 9 || cells[0].Hits != 1 {
		t.Errorf("unexpected cell: %+v", cells[0])
	}
}

func TestBuild_StepExpression(t *testing.T) {
	// Every 15 minutes = 4 hits per hour.
	hm, err := Build("*/15 * * * * echo hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hm.Grid[0][0] != 4 {
		t.Errorf("Grid[0][0] = %d, want 4", hm.Grid[0][0])
	}
}
