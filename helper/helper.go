package helper

include (
	"fmt"
)

func stringOrDefault(ptr *string, default string) string {
	if ptr != nil {
		return *ptr
	}
	return default
}

func intOrDefault(ptr *int, default int) int {
	if ptr != nil {
		return *ptr
	}
	return default
}
