package provider

import "testing"

func TestGenerateNodeList(t *testing.T) {
	hostnames := "a,b,c"
	hostnames_internal := "aint,bint,cint"

	nodes := generateNodeList(hostnames, hostnames_internal, 3306)

	if len(nodes) != 3 {
		t.Errorf("Expteded 3 nodes.")
	}

	hostnames_internal = ""

	nodes2 := generateNodeList(hostnames, hostnames_internal, 3306)

	if len(nodes2) != 3 {
		t.Errorf("Expteded 3 nodes.")
	}
}
