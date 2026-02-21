package imgproc

// JaroszFilter applies a Jarosz box filter to a float32 matrix.
// The filter is an iterative 1-D box average applied first along rows
// then along columns, each direction repeated nreps times.
// windowSize is the half-width of the box kernel; the full kernel width
// is 2*windowSize+1.
func JaroszFilter(buf [][]float32, windowSize, nreps int) {
	if len(buf) == 0 || len(buf[0]) == 0 {
		return
	}
	rows := len(buf)
	cols := len(buf[0])
	for i := 0; i < nreps; i++ {
		boxAlongRows(buf, rows, cols, windowSize)
	}
	for i := 0; i < nreps; i++ {
		boxAlongCols(buf, rows, cols, windowSize)
	}
}

func boxAlongRows(buf [][]float32, rows, cols, windowSize int) {
	w := float32(2*windowSize + 1)
	for i := 0; i < rows; i++ {
		row := buf[i]
		s := row[0] * float32(windowSize)
		for j := 1; j <= windowSize && j < cols; j++ {
			s += row[j]
		}
		s += row[0]
		row[0] = s / w
		for j := 1; j < cols; j++ {
			jn := j + windowSize
			jp := j - windowSize - 1
			var lead, trail float32
			if jn < cols {
				lead = row[jn]
			}
			if jp >= 0 {
				trail = row[jp]
			}
			s += lead - trail
			row[j] = s / w
		}
	}
}

func boxAlongCols(buf [][]float32, rows, cols, windowSize int) {
	w := float32(2*windowSize + 1)
	for j := 0; j < cols; j++ {
		s := buf[0][j] * float32(windowSize)
		for i := 1; i <= windowSize && i < rows; i++ {
			s += buf[i][j]
		}
		s += buf[0][j]
		buf[0][j] = s / w
		for i := 1; i < rows; i++ {
			in := i + windowSize
			ip := i - windowSize - 1
			var lead, trail float32
			if in < rows {
				lead = buf[in][j]
			}
			if ip >= 0 {
				trail = buf[ip][j]
			}
			s += lead - trail
			buf[i][j] = s / w
		}
	}
}
