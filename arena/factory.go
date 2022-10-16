package arena

type Factory interface {
	Create(Descriptor) Arena
}
