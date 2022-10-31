package conversion

type ConversionRepository interface {
}
type Service struct {
	conversionRepository ConversionRepository
}

func NewService(repository ConversionRepository) *Service {
	if repository == nil {
		return nil
	}
	return &Service{conversionRepository: repository}
}
