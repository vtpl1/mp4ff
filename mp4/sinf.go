package mp4

// SinfBox -  Protection Scheme Information Box according to ISO/IEC 23001-7
type SinfBox struct {
	Frma     *FrmaBox // Mandatory
	Schm     *SchmBox // Optional
	Schi     *SchiBox // Optional
	Children []Box
}
