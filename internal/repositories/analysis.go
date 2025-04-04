package repositories

type AnalysisRepository struct {
	Base Repository
}

func NewAnalysisRepository(base Repository) AnalysisRepository {
	return AnalysisRepository{Base: base}
}
