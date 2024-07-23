package expense

import "math"

type SplitType string

var SplitTypes = struct {
	Equal        SplitType
	Proportional SplitType
	Transfer     SplitType
}{
	Equal:        "equal",
	Proportional: "proportional",
	Transfer:     "transfer",
}

func (s SplitType) String() string {
	return string(s)
}

type SplitRatio struct {
	Payer    int
	Receiver int
}

func (s SplitRatio) Type() SplitType {
	if s.Payer == 0 {
		return SplitTypes.Transfer
	}

	if s.Payer == 50 {
		return SplitTypes.Equal
	}

	return SplitTypes.Proportional
}

func NewEqualSplitRatio() SplitRatio {
	return SplitRatio{
		Payer:    50,
		Receiver: 50,
	}
}

func NewProportionalSplitRatio(payerIncome, receiverIncome int) SplitRatio {
	totalIncome := payerIncome + receiverIncome
	payerRatio := int(math.Round(float64(payerIncome) * 100.0 / float64(totalIncome)))

	receiverRatio := 100 - payerRatio

	return SplitRatio{
		Payer:    payerRatio,
		Receiver: receiverRatio,
	}
}

func NewTransferRatio() SplitRatio {
	return SplitRatio{
		Payer:    0,
		Receiver: 100,
	}
}
