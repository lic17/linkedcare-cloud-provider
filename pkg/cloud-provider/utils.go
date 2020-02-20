package cloud_provider

import v1 "k8s.io/api/core/v1"

// NodeList return nodes list in string
func NodeList(nodes []*v1.Node) []string {
	ns := []string{}
	for _, node := range nodes {
		ns = append(ns, node.Name)
	}
	return ns
}
