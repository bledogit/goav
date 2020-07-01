package wrapper

type Scaler struct {
}

func NewScaler() (*Scaler, error) {
	s := &Scaler{}
	return s, nil
}
