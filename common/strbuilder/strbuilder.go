package strbuilder

// StrBuilder utility
type StrBuilder struct {
	sb         []string
	Delineator string
}

// SetDelineator sets delineator
func (t *StrBuilder) SetDelineator(delineator string) *StrBuilder {
	t.Delineator = delineator
	return t
}

// Add adds strings
func (t *StrBuilder) Add(s ...string) *StrBuilder {
	t.sb = append(t.sb, s...)
	return t
}

// A adds strings
func (t *StrBuilder) A(s ...string) *StrBuilder {
	t.sb = append(t.sb, s...)
	return t
}

// Build returns string
func (t *StrBuilder) Build() string {

	result := ""
	for _, s := range t.sb {
		if result == "" {
			result = s
		} else {
			result = result + t.Delineator + s
		}
	}

	return result
}

// NewStrBuilder returns new instance of StrBuilder with specified delineator
func NewStrBuilder(delineator string) *StrBuilder {
	return &StrBuilder{
		Delineator: delineator,
	}
}
