package boot

type Plugin func(*ServiceHub)

type InitService interface {
	Name() string
	InitFlags()
	Run() error
	Get() interface{}
	Stop() <-chan error
}
