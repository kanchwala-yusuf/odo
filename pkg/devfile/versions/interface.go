package versions

type DevfileData interface {
	Validate() error
	MergeDevfiles(interface{}) error
}
