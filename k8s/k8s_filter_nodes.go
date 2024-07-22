package k8s

import (
	"regexp"
	"strings"
)

// FilterNodesListCmd filters through the 'multipass list' command and returns only the cluster nodes in the same 'multipass list' format (assuming that there are other unrelated multipass nodes as well).
func FilterNodesListCmd(multipassListOutput string) string {
	lines := strings.Split(multipassListOutput, "\n")

	var result []string

	re := regexp.MustCompile(`^(haproxy|controller-node-[123]|worker-node-[123])\s+.*`)

	for _, line := range lines {
		if re.MatchString(line) || strings.HasPrefix(line, "Name") {
			result = append(result, line)
		}
	}

	//This result has already one string corresponding to the multipass default header, therefore a length of 1 corresponds to no deployed nodes.
	if len(result) == 1 {
		return ""
	}

	return strings.Join(result, "\n")
}

// GetCurrentNodes takes the filtered output from FilterNodesListCmd and returns a list of cluster node names.
func GetCurrentNodes(multipassListOutputFiltered string) []string {
	lines := strings.Split(multipassListOutputFiltered, "\n")

	var result []string
	// Excluding the first line which contains the multipass specific headers
	for _, line := range lines[1:] {
		nodeName := strings.Fields(line)[0]
		result = append(result, nodeName)
	}

	return result
}
