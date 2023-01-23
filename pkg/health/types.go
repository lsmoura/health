package health

import (
	"database/sql/driver"
	"encoding/xml"
	"time"
)

type HealthTime time.Time

func (t *HealthTime) UnmarshalXMLAttr(attr xml.Attr) error {
	parsed, err := time.Parse("2006-01-02 15:04:05 -0700", attr.Value)
	if err != nil {
		return err
	}
	*t = HealthTime(parsed)
	return nil
}
func (t HealthTime) String() string {
	return time.Time(t).String()
}
func (t HealthTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}
func (t HealthTime) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

type Me struct {
	DateOfBirth                 string `xml:"HKCharacteristicTypeIdentifierDateOfBirth,attr"`
	BiologicalSex               string `xml:"HKCharacteristicTypeIdentifierBiologicalSex,attr"`
	BloodType                   string `xml:"HKCharacteristicTypeIdentifierBloodType,attr"`
	FitzpatrickSkinType         string `xml:"HKCharacteristicTypeIdentifierFitzpatrickSkinType,attr"`
	CardioFitnessMedicationsUse string `xml:"HKCharacteristicTypeIdentifierCardioFitnessMedicationsUse,attr"`
}

type MetadataEntry struct {
	Key   string `xml:"key,attr" json:"key"`
	Value string `xml:"value,attr" json:"value"`
}

type InstantaneousBeatsPerMinute struct {
	BPM  int    `xml:"bpm,attr" json:"bpm"`
	Time string `xml:"time,attr" json:"time"`
}

type HeartRateVariabilityMetadataList struct {
	InstantaneousBeatsPerMinute []InstantaneousBeatsPerMinute `xml:"InstantaneousBeatsPerMinute"`
}

type Record struct {
	ID            int64       `xml:"-" db:"id"`
	Type          string      `xml:"type,attr" db:"type"` // required
	Unit          *string     `xml:"unit,attr" db:"unit"`
	Value         *string     `xml:"value,attr" db:"value"`
	SourceName    string      `xml:"sourceName,attr" db:"source_name"` // required
	SourceVersion *string     `xml:"sourceVersion,attr" db:"source_version"`
	Device        *string     `xml:"device,attr" db:"device"`
	CreationDate  *HealthTime `xml:"creationDate,attr" db:"creation_date"`
	StartDate     *HealthTime `xml:"startDate,attr" db:"start_date"` // required
	EndDate       *HealthTime `xml:"endDate,attr" db:"end_date"`     // required

	Metadata             []MetadataEntry                    `xml:"MetadataEntry" db:"metadata"`
	HeartRateVariability []HeartRateVariabilityMetadataList `xml:"HeartRateVariabilityMetadataList" db:"hrv"`
}

type Correlation struct {
	Type          string      `xml:"type,attr" db:"type"`              // required
	SourceName    string      `xml:"sourceName,attr" db:"source_name"` // required
	SourceVersion string      `xml:"sourceVersion,attr" db:"source_version"`
	Device        string      `xml:"device,attr" db:"device"`
	CreationDate  *HealthTime `xml:"creationDate,attr" db:"creation_date"`
	StartDate     *HealthTime `xml:"startDate,attr" db:"start_date"` // required
	EndDate       *HealthTime `xml:"endDate,attr" db:"end_date"`     // required

	Metadata []MetadataEntry `xml:"MetadataEntry" db:"metadata"`
	Records  []Record        `xml:"Record" db:"records"`
}

type WorkoutEvent struct {
	Type         string  `xml:"type,attr" json:"type"` // required
	Date         string  `xml:"date,attr" json:"date"` // required
	Duration     *string `xml:"duration,attr" json:"duration,omitempty"`
	DurationUnit *string `xml:"durationUnit,attr" json:"duration_unit,omitempty"`

	Metadata []MetadataEntry `xml:"MetadataEntry"`
}

type FileReference struct {
	Path string `xml:"path,attr" json:"path"`
}

type WorkoutRoute struct {
	SourceName    string `xml:"sourceName,attr" json:"source_name"` // required
	SourceVersion string `xml:"sourceVersion,attr" json:"source_version,omitempty"`
	Device        string `xml:"device,attr" json:"device,omitempty"`
	CreationDate  string `xml:"creationDate,attr" json:"creation_date,omitempty"`
	StartDate     string `xml:"startDate,attr" json:"start_date"` // required
	EndDate       string `xml:"endDate,attr" json:"end_date"`     // required

	Metadata      []MetadataEntry `xml:"MetadataEntry" db:"metadata,json" json:"metadata,omitempty"`
	FileReference []FileReference `xml:"FileReference" db:"file_reference,json" json:"file_reference,omitempty"`
}

type WorkoutStatistics struct {
	Type      string      `xml:"type,attr" json:"type"`            // required
	StartDate *HealthTime `xml:"startDate,attr" json:"start_date"` // required
	EndDate   *HealthTime `xml:"endDate,attr" json:"end_date"`     // required
	Average   *string     `xml:"average,attr" json:"average,omitempty"`
	Minimum   *string     `xml:"minimum,attr" json:"minimum,omitempty"`
	Maximum   *string     `xml:"maximum,attr" json:"maximum,omitempty"`
	Sum       *string     `xml:"sum,attr" json:"sum,omitempty"`
	Unit      *string     `xml:"unit,attr" json:"unit,omitempty"`
}

type Workout struct {
	ID                    int64       `db:"id"`
	WorkoutActivityType   string      `xml:"workoutActivityType,attr" db:"workout_activity_type"`
	Duration              float64     `xml:"duration,attr" db:"duration"`
	DurationUnit          string      `xml:"durationUnit,attr" db:"duration_unit"`
	TotalDistance         string      `xml:"totalDistance,attr" db:"total_distance"`
	TotalDistanceUnit     string      `xml:"totalDistanceUnit,attr" db:"total_distance_unit"`
	TotalEnergyBurned     string      `xml:"totalEnergyBurned,attr" db:"total_energy_burned"`
	TotalEnergyBurnedUnit string      `xml:"totalEnergyBurnedUnit,attr" db:"total_energy_burned_unit"`
	SourceName            string      `xml:"sourceName,attr" db:"source_name"`
	SourceVersion         string      `xml:"sourceVersion,attr" db:"source_version"`
	Device                string      `xml:"device,attr" db:"device"`
	CreationDate          *HealthTime `xml:"creationDate,attr" db:"creation_date"`
	StartDate             *HealthTime `xml:"startDate,attr" db:"start_date"`
	EndDate               *HealthTime `xml:"endDate,attr" db:"end_date"`

	Metadata          []MetadataEntry     `xml:"MetadataEntry" db:"metadata,json"`
	WorkoutEvent      []WorkoutEvent      `xml:"WorkoutEvent" db:"workout_events,json"`
	WorkoutRoute      []WorkoutRoute      `xml:"WorkoutRoute" db:"workout_routes,json"`
	WorkoutStatistics []WorkoutStatistics `xml:"WorkoutStatistics" db:"workout_statistics,json"`
}

type ActivitySummary struct {
	DateComponents         *string `xml:"dateComponents" db:"date_components"`
	ActiveEnergyBurned     *string `xml:"activeEnergyBurned" db:"active_energy_burned"`
	ActiveEnergyBurnedGoal *string `xml:"activeEnergyBurnedGoal" db:"active_energy_burned_goal"`
	ActiveEnergyBurnedUnit *string `xml:"activeEnergyBurnedUnit" db:"active_energy_burned_unit"`
	AppleMoveTime          *string `xml:"appleMoveTime" db:"apple_move_time"`
	AppleMoveTimeGoal      *string `xml:"appleMoveTimeGoal" db:"apple_move_time_goal"`
	AppleExerciseTime      *string `xml:"appleExerciseTime" db:"apple_exercise_time"`
	AppleExerciseTimeGoal  *string `xml:"appleExerciseTimeGoal" db:"apple_exercise_time_goal"`
	AppleStandHours        *string `xml:"appleStandHours" db:"apple_stand_hours"`
	AppleStandHoursGoal    *string `xml:"appleStandHoursGoal" db:"apple_stand_hours_goal"`
}

type ClinicalRecord struct {
	Type             *string `xml:"type" db:"type"`
	Identifier       *string `xml:"identifier" db:"identifier"`
	SourceName       *string `xml:"sourceName" db:"source_name"`
	SourceURL        *string `xml:"sourceURL" db:"source_url"`
	FhirVersion      *string `xml:"fhirVersion" db:"fhir_version"`
	ReceivedDate     *string `xml:"receivedDate" db:"received_date"`
	ResourceFilePath *string `xml:"resourceFilePath" db:"resource_file_path"`
}

type SensitivityPoint struct {
	FrequencyValue string  `xml:"frequencyValue,attr" json:"frequency_value"`
	FrequencyUnit  string  `xml:"frequencyUnit,attr" json:"frequency_unit"`
	LeftEarValue   *string `xml:"leftEarValue,attr" json:"left_ear_value,omitempty"`
	LeftEarUnit    *string `xml:"leftEarUnit,attr" json:"left_ear_unit,omitempty"`
	RightEarValue  *string `xml:"rightEarValue,attr" json:"right_ear_value,omitempty"`
	RightEarUnit   *string `xml:"rightEarUnit,attr" json:"right_ear_unit,omitempty"`
}

type Audiogram struct {
	Type          string      `xml:"type" db:"type"`
	SourceName    string      `xml:"sourceName" db:"source_name"`
	SourceVersion *string     `xml:"sourceVersion" db:"source_version"`
	Device        *string     `xml:"device" db:"device"`
	CreationDate  *HealthTime `xml:"creationDate" db:"creation_date,omitempty"`
	StartDate     *HealthTime `xml:"startDate" db:"start_date"`
	EndDate       *HealthTime `xml:"endDate" db:"end_date"`

	Metadata         []MetadataEntry    `xml:"MetadataEntry" db:"metadata,json"`
	SensitivityPoint []SensitivityPoint `xml:"SensitivityPoint" db:"sensitivity_points,json"`
}

type Eye struct {
	Sphere          *string `xml:"sphere,attr" json:"sphere,omitempty"`
	SphereUnit      *string `xml:"sphereUnit,attr" json:"sphere_unit,omitempty"`
	Cylinder        *string `xml:"cylinder,attr" json:"cylinder,omitempty"`
	CylinderUnit    *string `xml:"cylinderUnit,attr" json:"cylinder_unit,omitempty"`
	Axis            *string `xml:"axis,attr" json:"axis,omitempty"`
	AxisUnit        *string `xml:"axisUnit,attr" json:"axis_unit,omitempty"`
	Add             *string `xml:"add,attr" json:"add,omitempty"`
	AddUnit         *string `xml:"addUnit,attr" json:"add_unit,omitempty"`
	Vertex          *string `xml:"vertex,attr" json:"vertex,omitempty"`
	VertexUnit      *string `xml:"vertexUnit,attr" json:"vertex_unit,omitempty"`
	PrismAmount     *string `xml:"prismAmount,attr" json:"prism_amount,omitempty"`
	PrismAmountUnit *string `xml:"prismAmountUnit,attr" json:"prism_amount_unit,omitempty"`
	PrismAngle      *string `xml:"prismAngle,attr" json:"prism_angle,omitempty"`
	PrismAngleUnit  *string `xml:"prismAngleUnit,attr" json:"prism_angle_unit,omitempty"`
	FarPD           *string `xml:"farPD,attr" json:"far_pd,omitempty"`
	FarPDUnit       *string `xml:"farPDUnit,attr" json:"far_pdunit,omitempty"`
	NearPD          *string `xml:"nearPD,attr" json:"near_pd,omitempty"`
	NearPDUnit      *string `xml:"nearPDUnit,attr" json:"near_pdunit,omitempty"`
	BaseCurve       *string `xml:"baseCurve,attr" json:"base_curve,omitempty"`
	BaseCurveUnit   *string `xml:"baseCurveUnit,attr" json:"base_curve_unit,omitempty"`
	Diameter        *string `xml:"diameter,attr" json:"diameter,omitempty"`
	DiameterUnit    *string `xml:"diameterUnit,attr" json:"diameter_unit,omitempty"`
}

type Attachment struct {
	Identifier *string `xml:"identifier" json:"identifier,omitempty"`
}

type VisionPrescription struct {
	Type           string  `xml:"type" db:"type"`
	DateIssued     string  `xml:"dateIssued" db:"date_issued"`
	ExpirationDate *string `xml:"expirationDate" db:"expiration_date"`
	Brand          *string `xml:"brand" db:"brand"`

	Metadata   []MetadataEntry `xml:"MetadataEntry" db:"metadata,json"`
	RightEye   []Eye           `xml:"RightEye" db:"right_eye,json"`
	LeftEye    []Eye           `xml:"LeftEye" db:"left_eye,json"`
	Attachment []Attachment    `xml:"Attachment" db:"attachment,json"`
}

type HealthData struct {
	ExportDate         string               `xml:"ExportDate"`
	Me                 Me                   `xml:"Me"`
	Records            []Record             `xml:"Record"`
	Correlations       []Correlation        `xml:"Correlation"`
	Workouts           []Workout            `xml:"Workout"`
	ActivitySummary    []ActivitySummary    `xml:"ActivitySummary"`
	ClinicalRecord     []ClinicalRecord     `xml:"ClinicalRecord"`
	Audiogram          []Audiogram          `xml:"Audiogram"`
	VisionPrescription []VisionPrescription `xml:"VisionPrescription"`
}
