package search

type NotFoundError struct {
	message string
}

func (e NotFoundError) Error() string {
	return e.message
}

type MultipleFoundError struct {
	message string
}

func (e MultipleFoundError) Error() string {
	return e.message
}

type DatacenterNotFoundError struct {
	dc      string
	message string
}

func (d DatacenterNotFoundError) Error() string {
	return d.message
}
