package mapUtils

const (
	// ForEachOperationBreak will just continue the loop without doing anything.
	ForEachOperationContinue = 0

	// ForEachOperationBreak will just break the loop without doing anything.
	ForEachOperationBreak = 1

	// ForEachOperationBreak will just remove the current item from the list
	// and continue the loop.
	ForEachOperationRemove = 2

	// ForEachOperationBreak will remove the current item from the list
	// and break the loop.
	ForEachOperationRemoveBreak = 3
)

const (
	checkActionNormal checkAction = iota
	checkActionContinue
	checkActionBreak
	checkActionReturn
)
