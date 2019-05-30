package eddsa

import "github.com/cloudflare/circl/math/fp448"

var tabSign448 = [2][_2w1]pointR3{
	[_2w1]pointR3{
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
	},
	[_2w1]pointR3{
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
		&point448R3{
			addYX: [fp448.Size]byte{},
			subYX: [fp448.Size]byte{},
			dt2:   [fp448.Size]byte{},
		},
	},
}

var tabVerif448 = [numPointsVerif]pointR3{
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
	&point448R3{
		addYX: [fp448.Size]byte{},
		subYX: [fp448.Size]byte{},
		dt2:   [fp448.Size]byte{},
	},
}
