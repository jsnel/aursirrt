package dock

//A Docker handles incoming communication.
type Docker interface {
	Launch(agent DockAgent) //Launches a docker
}
