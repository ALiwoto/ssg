package mapUtils

const (
	// ForEachOperationBreak will just continue the loop without doing anything.
	ForEachOperationContinue ForEachOperation = iota

	// ForEachOperationBreak will just break the loop without doing anything.
	ForEachOperationBreak

	// ForEachOperationBreak will just remove the current item from the list
	// and continue the loop.
	ForEachOperationRemove

	// ForEachOperationBreak will remove the current item from the list
	// and break the loop.
	ForEachOperationRemoveBreak
)

const (
	checkActionNormal checkAction = iota
	checkActionContinue
	checkActionBreak
	checkActionReturn
)
