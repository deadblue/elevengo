package option

type CooldownOption struct{
	// Minimum cooldown duration in millisecond
	Min uint
	// Maximum cooldown duration in millisecond
	Max uint
}

func (o CooldownOption) isOption() {}