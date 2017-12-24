package blockchain

// refuse a match
// param : DrID , DrSig , MatchID
func (setup *FabricSetup) RefuseMatch(param []string) (string, error) {

	// Prepare arguments = function + parameters
	var args []string
	args = append(args, "refuseMatch")
	args = append(args, param...)

	return setup.invoke(args)
}
