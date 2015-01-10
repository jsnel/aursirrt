package dock

//A Docker handles incoming communication.
type Docker interface {
	Launch(agent DockAgent) error //Launches a docker
}
