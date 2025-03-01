package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestTlou(t *testing.T) {
	tlou := &mp4.TlouBox{
		LoudnessBaseBox: mp4.LoudnessBaseBox{
			Version: 1,
			Flags:   0,
			LoudnessBases: []*mp4.LoudnessBase{
				{
					EQSetID:                0,
					DownmixID:              0,
					DRCSetID:               0,
					BsSamplePeakLevel:      1087,
					BsTruePeakLevel:        1086,
					MeasurementSystemForTP: 2,
					ReliabilityForTP:       3,
					Measurements: []mp4.Measurement{
						{
							MethodDefinition:  1,
							MethodValue:       121,
							MeasurementSystem: 2,
							Reliability:       3,
						},
						{
							MethodDefinition:  3,
							MethodValue:       122,
							MeasurementSystem: 1,
							Reliability:       3,
						},
					},
				},
				{
					EQSetID:                0,
					DownmixID:              0,
					DRCSetID:               0,
					BsSamplePeakLevel:      1087,
					BsTruePeakLevel:        1086,
					MeasurementSystemForTP: 2,
					ReliabilityForTP:       3,
					Measurements: []mp4.Measurement{
						{
							MethodDefinition:  4,
							MethodValue:       124,
							MeasurementSystem: 1,
							Reliability:       3,
						},
						{
							MethodDefinition:  5,
							MethodValue:       122,
							MeasurementSystem: 1,
							Reliability:       3,
						},
					},
				},
			},
		},
	}
	boxDiffAfterEncodeAndDecode(t, tlou)
}

func TestAlou(t *testing.T) {
	alou := &mp4.AlouBox{
		LoudnessBaseBox: mp4.LoudnessBaseBox{
			Version: 1,
			Flags:   0,
			LoudnessBases: []*mp4.LoudnessBase{
				{
					EQSetID:                0,
					DownmixID:              0,
					DRCSetID:               0,
					BsSamplePeakLevel:      1087,
					BsTruePeakLevel:        1086,
					MeasurementSystemForTP: 2,
					ReliabilityForTP:       3,
					Measurements: []mp4.Measurement{
						{
							MethodDefinition:  1,
							MethodValue:       121,
							MeasurementSystem: 2,
							Reliability:       3,
						},
						{
							MethodDefinition:  3,
							MethodValue:       122,
							MeasurementSystem: 1,
							Reliability:       3,
						},
					},
				},
				{
					EQSetID:                0,
					DownmixID:              0,
					DRCSetID:               0,
					BsSamplePeakLevel:      1087,
					BsTruePeakLevel:        1086,
					MeasurementSystemForTP: 2,
					ReliabilityForTP:       3,
					Measurements: []mp4.Measurement{
						{
							MethodDefinition:  4,
							MethodValue:       124,
							MeasurementSystem: 1,
							Reliability:       3,
						},
						{
							MethodDefinition:  5,
							MethodValue:       122,
							MeasurementSystem: 1,
							Reliability:       3,
						},
					},
				},
			},
		},
	}
	boxDiffAfterEncodeAndDecode(t, alou)
}
