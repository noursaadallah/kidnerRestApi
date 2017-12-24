package blockchain

// CreateDoctor and return its ID
func (setup *FabricSetup) CreateDoctor(signature string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "createDoctor")
	args = append(args, signature)

	return setup.invoke(args)
}
