package biglin

import (
	"fmt"
	"github.com/jackvalmadre/go-vec"
	"math"
)

type GradientDescent struct {
	LineSearch bool
	StepSize   float64
	Backtrack  bool
}

func (opts GradientDescent) Solve(f Objective, x0 vec.ConstTyped,
	crit TerminationCriteria, callback Callback, verbose bool) (vec.MutableTyped, error) {
	k := 0
	t := opts.StepSize
	f_x := math.Inf(1)
	x := vec.Clone(x0)
	var g_x0 vec.MutableTyped
	var x_prev vec.ConstTyped

	for {
		f_x_prev := f_x
		var g_x vec.MutableTyped
		err := f.Evaluate(x, &f_x, &g_x)
		if err != nil {
			return nil, err
		}
		if k == 0 {
			g_x0 = g_x
		}

		summary := Summarize(k, f_x_prev, f_x, g_x0, g_x, x_prev, x)
		fmt.Println(summary)
		if callback != nil {
			callback(summary)
		}
		converged := crit.Evaluate(summary)
		if converged {
			break
		}

		x_prev = x
		if opts.LineSearch {
			// Update x by exact line search.
			t, err = f.LineSearch(x, g_x)
			if err != nil {
				return nil, err
			}
			// x <- x - t g_x
			vec.CopyTo(x, vec.Plus(x, vec.Scale(t, g_x)))
		} else {
			if !opts.Backtrack {
				// x <- x - t g_x
				vec.CopyTo(x, vec.Plus(x, vec.Scale(-t, g_x)))
			} else {
				var z vec.MutableTyped
				for satisfied := false; !satisfied; {
					// x <- x - t g_x
					vec.CopyTo(z, vec.Plus(x, vec.Scale(-t, g_x)))
					var f_z float64
					err = f.Evaluate(z, &f_z, nil)
					if err != nil {
						return nil, err
					}
					// Check backtrack criterion.
					lhs := f_x - f_z
					rhs := 0.5 * t * vec.SqrNorm(g_x)
					if lhs >= rhs {
						satisfied = true
					} else {
						t /= 2
					}
				}
				x = z
			}
		}

		k += 1
	}

	return x, nil
}
