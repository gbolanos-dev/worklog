package prompts

import _ "embed"

//go:embed standup.md
var Standup string

//go:embed summary.md
var Summary string

//go:embed promo.md
var Promo string
