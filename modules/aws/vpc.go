package aws

import (
	"fmt"
	"strings"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Get a random CIDR block from the range of acceptable private IP addresses per RFC 1918
// (https://tools.ietf.org/html/rfc1918#section-3)
// The routingPrefix refers to the "/28" in 1.2.3.4/28.
// Note that, as written, this function will return a subset of all valid ranges. Since we will probably use this function
// mostly for generating random CIDR ranges for VPCs and Subnets, having comprehensive set coverage is not essential.
func GetRandomPrivateCidrBlock(routingPrefix int) string {

	var o1, o2, o3, o4 int

	switch routingPrefix {
	case 32:
		o1 = random.RandomInt([]int{10, 172, 192})

		switch o1 {
		case 10:
			o2 = random.Random(0, 255)
			o3 = random.Random(0, 255)
			o4 = random.Random(0, 255)
		case 172:
			o2 = random.Random(16, 31)
			o3 = random.Random(0, 255)
			o4 = random.Random(0, 255)
		case 192:
			o2 = 168
			o3 = random.Random(0, 255)
			o4 = random.Random(0, 255)
		}

	case 31, 30, 29, 28, 27, 26, 25:
		fallthrough
	case 24:
		o1 = random.RandomInt([]int{10, 172, 192})

		switch o1 {
		case 10:
			o2 = random.Random(0, 255)
			o3 = random.Random(0, 255)
			o4 = 0
		case 172:
			o2 = 16
			o3 = 0
			o4 = 0
		case 192:
			o2 = 168
			o3 = 0
			o4 = 0
		}
	case 23, 22, 21, 20, 19:
		fallthrough
	case 18:
		o1 = random.RandomInt([]int{10, 172, 192})

		switch o1 {
		case 10:
			o2 = 0
			o3 = 0
			o4 = 0
		case 172:
			o2 = 16
			o3 = 0
			o4 = 0
		case 192:
			o2 = 168
			o3 = 0
			o4 = 0
		}
	}
	return fmt.Sprintf("%d.%d.%d.%d/%d", o1, o2, o3, o4, routingPrefix)
}

func GetFirstTwoOctets(cidrBlock string) string {
	ipAddr := strings.Split(cidrBlock, "/")[0]
	octets := strings.Split(ipAddr, ".")
	return octets[0] + "." + octets[1]
}
