package Sort

func InsertSortInt64array(slice []int64, reverse bool) {
	n := len(slice)
	if n < 2 {
		return
	}
	if !reverse {
		for i := 1; i < n; i++ {
			for j := i; j > 0 && slice[j] < slice[j-1]; j-- {
				slice[j], slice[j-1] = slice[j-1], slice[j]
			}
		}
	} else {
		for i := 1; i < n; i++ {
			for j := i; j > 0 && slice[j] > slice[j-1]; j-- {
				slice[j], slice[j-1] = slice[j-1], slice[j]
			}
		}
	}
}
