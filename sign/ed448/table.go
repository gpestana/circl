package ed448

import "github.com/cloudflare/circl/math/fp448"

var tabSign = [fxV][fx2w1]pointR3{
	[fx2w1]pointR3{
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
	},
	[fx2w1]pointR3{
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		pointR3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
	},
}

var tabVerif = [1 << (omegaFix - 2)]pointR3{
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	pointR3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
}
