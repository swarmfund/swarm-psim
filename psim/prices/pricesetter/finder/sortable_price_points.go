package finder

type sortablePricePointsByTime []providerPricePoint

func (s sortablePricePointsByTime) Len() int {
	return len(s)
}

func (s sortablePricePointsByTime) Less(i, j int) bool {
	return s[i].Time.After(s[j].Time)
}

func (s sortablePricePointsByTime) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}


type sortablePricePointsByPrice []providerPricePoint

func (s sortablePricePointsByPrice) Len() int {
	return len(s)
}

func (s sortablePricePointsByPrice) Less(i, j int) bool {
	return s[i].Price > s[j].Price
}

func (s sortablePricePointsByPrice) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}
